package service

//go:generate mockgen -source=firestore_service.go -destination=firestore_service_gomock.go -package=service mountain_service/activity/internal/service -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/firestorex"
	sharedModels "github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type FirestoreService interface {
	ListByUrgency(ctx context.Context, urgencyID uint, limit int) ([]sharedModels.Activity, error)
	ListAll(ctx context.Context, limit int) ([]sharedModels.Activity, error)
}

type firestoreService struct {
	client     firestorex.Client
	logger     utils.Logger
	collection string
}

func NewFirebaseReadService(client firestorex.Client, logger utils.Logger) FirestoreService {
	return &firestoreService{
		client:     client,
		logger:     logger.WithName("firebaseReadService"),
		collection: "activities",
	}
}

func (s *firestoreService) ListByUrgency(ctx context.Context, urgencyID uint, limit int) ([]sharedModels.Activity, error) {
	if s.client == nil {
		return nil, fmt.Errorf("firestore client is nil")
	}

	q := s.client.Collection(s.collection).
		Where("urgency_id", "==", urgencyID).
		OrderBy("created_at", firestorex.Desc)
	if limit > 0 {
		q = q.Limit(limit)
	}

	iter := q.Documents(ctx)
	defer iter.Stop()

	var items []sharedModels.Activity
	for {
		doc, err := iter.Next()
		if isDone(err) {
			break
		}
		if err != nil {
			log := s.logger.WithContext(ctx)
			log.Errorf("failed to iterate firestore docs: %v", err)
			return nil, err
		}

		var a struct {
			ID          uint   `firestore:"id"`
			Description string `firestore:"description"`
			EmployeeID  uint   `firestore:"employee_id"`
			UrgencyID   uint   `firestore:"urgency_id"`
			CreatedAt   string `firestore:"created_at"`
			UpdatedAt   string `firestore:"updated_at"`
		}
		if err := doc.DataTo(&a); err != nil {
			log := s.logger.WithContext(ctx)
			log.Errorf("failed to unmarshal firestore doc: %v", err)
			continue
		}
		item := sharedModels.Activity{
			ID:          a.ID,
			Description: a.Description,
			EmployeeID:  a.EmployeeID,
			UrgencyID:   a.UrgencyID,
			CreatedAt:   s.parseTime(a.CreatedAt),
			UpdatedAt:   s.parseTime(a.UpdatedAt),
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *firestoreService) ListAll(ctx context.Context, limit int) ([]sharedModels.Activity, error) {
	if s.client == nil {
		return nil, fmt.Errorf("firestore client is nil")
	}

	q := s.client.Collection(s.collection).OrderBy("created_at", firestorex.Desc)
	if limit > 0 {
		q = q.Limit(limit)
	}

	iter := q.Documents(ctx)
	defer iter.Stop()

	var items []sharedModels.Activity
	for {
		doc, err := iter.Next()
		if isDone(err) {
			break
		}
		if err != nil {
			log := s.logger.WithContext(ctx)
			log.Errorf("failed to iterate firestore docs: %v", err)
			return nil, err
		}
		var a struct {
			ID          uint   `firestore:"id"`
			Description string `firestore:"description"`
			EmployeeID  uint   `firestore:"employee_id"`
			UrgencyID   uint   `firestore:"urgency_id"`
			CreatedAt   string `firestore:"created_at"`
			UpdatedAt   string `firestore:"updated_at"`
		}
		if err := doc.DataTo(&a); err != nil {
			log := s.logger.WithContext(ctx)
			log.Errorf("failed to unmarshal firestore doc: %v", err)
			continue
		}
		item := sharedModels.Activity{
			ID:          a.ID,
			Description: a.Description,
			EmployeeID:  a.EmployeeID,
			UrgencyID:   a.UrgencyID,
			CreatedAt:   s.parseTime(a.CreatedAt),
			UpdatedAt:   s.parseTime(a.UpdatedAt),
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *firestoreService) parseTime(ts string) (t time.Time) {
	if ts == "" {
		return
	}
	if parsed, err := time.Parse(time.RFC3339, ts); err == nil {
		return parsed
	}
	return
}

func isDone(err error) bool { return firestorex.IsDone(err) }
