package service

//go:generate mockgen -source=firebase_service.go -destination=firebase_service_gomock.go -package=service mountain_service/activity/internal/service -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"fmt"
	"strconv"
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/firestorex"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

// FirebaseService handles Firebase Firestore operations for read model
type FirebaseService interface {
	GetActivitiesByUrgency(ctx context.Context, urgencyID uint) ([]*models.Activity, error)
	GetAllActivities(ctx context.Context, limit int) ([]*models.Activity, error)
	SyncActivity(ctx context.Context, eventData activityV1.ActivityEvent) error
	HealthCheck(ctx context.Context) error
}

type firebaseService struct {
	client     firestorex.Client
	logger     utils.Logger
	collection string
}

// FirebaseActivityDoc represents the document structure in Firestore
type FirebaseActivityDoc struct {
	ID           int64     `firestore:"id"`
	UrgencyID    int64     `firestore:"urgency_id"`
	EmployeeID   int64     `firestore:"employee_id"`
	Description  string    `firestore:"description"`
	CreatedAt    time.Time `firestore:"created_at"`
	EmployeeName string    `firestore:"employee_name"`
	UrgencyTitle string    `firestore:"urgency_title"`
	UrgencyLevel string    `firestore:"urgency_level"`
	SyncedAt     time.Time `firestore:"synced_at"`
	Version      int       `firestore:"version"`
}

func NewFirebaseService(client firestorex.Client, logger utils.Logger) FirebaseService {
	return &firebaseService{
		client:     client,
		logger:     logger.WithName("firebaseService"),
		collection: "activities",
	}
}
func isDone(err error) bool { return firestorex.IsDone(err) }

func (s *firebaseService) GetActivitiesByUrgency(ctx context.Context, urgencyID uint) ([]*models.Activity, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Firestore client is nil")
	}

	log := s.logger.WithContext(ctx)
	log.Infof("Getting activities from Firebase for urgency: %d", urgencyID)

	iter := s.client.Collection(s.collection).
		Where("urgency_id", "==", int64(urgencyID)).
		OrderBy("created_at", firestorex.Desc).
		Documents(ctx)

	var activities []*models.Activity

	for {
		doc, err := iter.Next()
		if isDone(err) {
			break
		}
		if err != nil {
			s.logger.Errorf("Failed to iterate Firebase documents: %v", err)
			return nil, fmt.Errorf("failed to get activities from Firebase: %w", err)
		}

		var fbDoc FirebaseActivityDoc
		if err := doc.DataTo(&fbDoc); err != nil {
			log.Errorf("Failed to unmarshal Firebase document: %v", err)
			continue
		}

		activity := &models.Activity{
			ID:          uint(fbDoc.ID),
			Description: fbDoc.Description,
		}

		activities = append(activities, activity)
	}

	log.Infof("Retrieved %d activities from Firebase for urgency %d", len(activities), urgencyID)
	return activities, nil
}

func (s *firebaseService) GetAllActivities(ctx context.Context, limit int) ([]*models.Activity, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Firestore client is nil")
	}

	log := s.logger.WithContext(ctx)
	log.Infof("Getting all activities from Firebase with limit: %d", limit)

	query := s.client.Collection(s.collection).OrderBy("created_at", firestorex.Desc)
	if limit > 0 {
		query = query.Limit(limit)
	}

	iter := query.Documents(ctx)
	var activities []*models.Activity

	for {
		doc, err := iter.Next()
		if isDone(err) {
			break
		}
		if err != nil {
			s.logger.Errorf("Failed to iterate Firebase documents: %v", err)
			return nil, fmt.Errorf("failed to get activities from Firebase: %w", err)
		}

		var fbDoc FirebaseActivityDoc
		if err := doc.DataTo(&fbDoc); err != nil {
			log.Errorf("Failed to unmarshal Firebase document: %v", err)
			continue
		}

		// Convert Firebase doc to Activity model
		activity := &models.Activity{
			ID:          uint(fbDoc.ID),
			Description: fbDoc.Description,
			// Map other fields as needed
		}

		activities = append(activities, activity)
	}

	log.Infof("Retrieved %d activities from Firebase", len(activities))
	return activities, nil
}

func (s *firebaseService) SyncActivity(ctx context.Context, eventData activityV1.ActivityEvent) error {
	if s.client == nil {
		return fmt.Errorf("Firestore client is nil")
	}
	if eventData.ActivityID == 0 {
		return fmt.Errorf("invalid activity id: 0")
	}

	log := s.logger.WithContext(ctx)
	log.Infof("Syncing activity to Firebase: activity_id=%d, type=%s", eventData.ActivityID, eventData.Type)

	docID := strconv.Itoa(int(eventData.ActivityID))
	docRef := s.client.Collection(s.collection).Doc(docID)

	switch eventData.Type {
	case "CREATE":
		fbDoc := FirebaseActivityDoc{
			ID:           int64(eventData.ActivityID),
			UrgencyID:    int64(eventData.UrgencyID),
			EmployeeID:   int64(eventData.EmployeeID),
			Description:  eventData.Description,
			CreatedAt:    eventData.CreatedAt.UTC(),
			EmployeeName: eventData.EmployeeName,
			UrgencyTitle: eventData.UrgencyTitle,
			UrgencyLevel: eventData.UrgencyLevel,
			SyncedAt:     time.Now().UTC(),
			Version:      1,
		}

		_, err := docRef.Set(ctx, fbDoc)
		if err != nil {
			s.logger.Errorf("Failed to create activity in Firebase: activity_id=%d, error=%v", eventData.ActivityID, err)
			return fmt.Errorf("failed to create activity in Firebase: %w", err)
		}

	case "UPDATE":
		updates := []firestorex.Update{
			{Path: "description", Value: eventData.Description},
			{Path: "employee_name", Value: eventData.EmployeeName},
			{Path: "urgency_title", Value: eventData.UrgencyTitle},
			{Path: "urgency_level", Value: eventData.UrgencyLevel},
			{Path: "synced_at", Value: firestorex.ServerTimestamp()},
			{Path: "version", Value: firestorex.Increment(1)},
		}

		_, err := docRef.Update(ctx, updates)
		if err != nil {
			s.logger.Errorf("Failed to update activity in Firebase: activity_id=%d, error=%v", eventData.ActivityID, err)
			return fmt.Errorf("failed to update activity in Firebase: %w", err)
		}

	case "DELETE":
		_, err := docRef.Delete(ctx)
		if err != nil {
			s.logger.Errorf("Failed to delete activity in Firebase: activity_id=%d, error=%v", eventData.ActivityID, err)
			return fmt.Errorf("failed to delete activity in Firebase: %w", err)
		}

	default:
		return fmt.Errorf("unknown event type: %s", eventData.Type)
	}

	log.Infof("Activity synced to Firebase successfully: activity_id=%d, type=%s", eventData.ActivityID, eventData.Type)
	return nil
}

func (s *firebaseService) HealthCheck(ctx context.Context) error {
	if s.client == nil {
		return fmt.Errorf("failed to check health: Firestore client is nil")
	}

	// Try to read a single document to test connectivity
	iter := s.client.Collection(s.collection).Limit(1).Documents(ctx)
	_, err := iter.Next()
	log := s.logger.WithContext(ctx)
	if err != nil && !isDone(err) {
		log.Errorf("Firebase health check failed: %v", err)
		return fmt.Errorf("failed to check health: %w", err)
	}

	log.Debug("Firebase health check passed")
	return nil
}
