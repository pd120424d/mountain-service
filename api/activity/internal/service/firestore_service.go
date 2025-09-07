package service

//go:generate mockgen -source=firestore_service.go -destination=firestore_service_gomock.go -package=service mountain_service/activity/internal/service -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/firestorex"
	sharedModels "github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type FirestoreService interface {
	ListByUrgency(ctx context.Context, urgencyID uint, limit int) ([]sharedModels.Activity, error)
	ListAll(ctx context.Context, limit int) ([]sharedModels.Activity, error)
	ListByUrgencyCursor(ctx context.Context, urgencyID uint, pageSize int, pageToken string) ([]sharedModels.Activity, string, error)
	ListAllCursor(ctx context.Context, pageSize int, pageToken string) ([]sharedModels.Activity, string, error)
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
	log := s.logger.WithContext(ctx)
	log.Infof("Listing activities by urgency: %d", urgencyID)
	defer utils.TimeOperation(log, "FirestoreService.ListByUrgency")()

	if s.client == nil {
		return nil, fmt.Errorf("firestore client is nil")
	}

	q := s.client.Collection(s.collection).
		Where("urgency_id", "==", int64(urgencyID)).
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
			log.Errorf("failed to iterate firestore docs: %v", err)
			return nil, err
		}

		var a struct {
			ID          int64       `firestore:"id"`
			Description string      `firestore:"description"`
			EmployeeID  int64       `firestore:"employee_id"`
			UrgencyID   int64       `firestore:"urgency_id"`
			CreatedAt   interface{} `firestore:"created_at"`
			UpdatedAt   interface{} `firestore:"updated_at"`
		}
		if err := doc.DataTo(&a); err != nil {
			log.Errorf("failed to unmarshal firestore doc: %v", err)
			continue
		}
		item := sharedModels.Activity{
			ID:          uint(a.ID),
			Description: a.Description,
			EmployeeID:  uint(a.EmployeeID),
			UrgencyID:   uint(a.UrgencyID),
			CreatedAt:   coerceTime(a.CreatedAt),
			UpdatedAt:   coerceTime(a.UpdatedAt),
		}
		items = append(items, item)
	}

	log.Infof("Successfully listed %d activities by urgency %d", len(items), urgencyID)
	return items, nil
}

func (s *firestoreService) ListAll(ctx context.Context, limit int) ([]sharedModels.Activity, error) {
	log := s.logger.WithContext(ctx)
	log.Infof("Listing all activities with limit: %d", limit)
	defer utils.TimeOperation(log, "FirestoreService.ListAll")()

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
			log.Errorf("failed to iterate firestore docs: %v", err)
			return nil, err
		}
		var a struct {
			ID          int64     `firestore:"id"`
			Description string    `firestore:"description"`
			EmployeeID  int64     `firestore:"employee_id"`
			UrgencyID   int64     `firestore:"urgency_id"`
			CreatedAt   time.Time `firestore:"created_at"`
			UpdatedAt   time.Time `firestore:"updated_at"`
		}
		if err := doc.DataTo(&a); err != nil {
			log.Errorf("failed to unmarshal firestore doc: %v", err)
			continue
		}
		item := sharedModels.Activity{
			ID:          uint(a.ID),
			Description: a.Description,
			EmployeeID:  uint(a.EmployeeID),
			UrgencyID:   uint(a.UrgencyID),
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		}
		items = append(items, item)
	}

	log.Infof("Successfully listed %d activities", len(items))
	return items, nil
}

// Cursor token helpers
type cursorToken struct {
	CreatedAt string `json:"createdAt"`
}

func encodeToken(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	b, _ := json.Marshal(cursorToken{CreatedAt: t.UTC().Format(time.RFC3339)})
	return base64.StdEncoding.EncodeToString(b)
}

func decodeToken(token string) (time.Time, error) {
	if token == "" {
		return time.Time{}, nil
	}
	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return time.Time{}, err
	}
	var ct cursorToken
	if err := json.Unmarshal(raw, &ct); err != nil {
		return time.Time{}, err
	}
	if ct.CreatedAt == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, ct.CreatedAt)
}

func (s *firestoreService) ListByUrgencyCursor(ctx context.Context, urgencyID uint, pageSize int, pageToken string) ([]sharedModels.Activity, string, error) {
	log := s.logger.WithContext(ctx)
	defer utils.TimeOperation(log, "FirestoreService.ListByUrgencyCursor")()
	if s.client == nil {
		return nil, "", fmt.Errorf("firestore client is nil")
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	q := s.client.Collection(s.collection).
		Where("urgency_id", "==", int64(urgencyID)).
		OrderBy("created_at", firestorex.Desc)
	if pageToken != "" {
		if t, err := decodeToken(pageToken); err == nil && !t.IsZero() {
			q = q.StartAfter(t)
		}
	}
	q = q.Limit(pageSize + 1)
	it := q.Documents(ctx)
	defer it.Stop()
	var items []sharedModels.Activity
	for {
		doc, err := it.Next()
		if isDone(err) {
			break
		}
		if err != nil {
			return nil, "", err
		}
		var a struct {
			ID          int64       `firestore:"id"`
			Description string      `firestore:"description"`
			EmployeeID  int64       `firestore:"employee_id"`
			UrgencyID   int64       `firestore:"urgency_id"`
			CreatedAt   interface{} `firestore:"created_at"`
			UpdatedAt   interface{} `firestore:"updated_at"`
		}
		if err := doc.DataTo(&a); err != nil {
			continue
		}
		items = append(items, sharedModels.Activity{
			ID: uint(a.ID), Description: a.Description,
			EmployeeID: uint(a.EmployeeID), UrgencyID: uint(a.UrgencyID),
			CreatedAt: coerceTime(a.CreatedAt), UpdatedAt: coerceTime(a.UpdatedAt),
		})
	}
	var next string
	if len(items) > pageSize {
		last := items[pageSize-1]
		next = encodeToken(last.CreatedAt)
		items = items[:pageSize]
	}
	return items, next, nil
}

func (s *firestoreService) ListAllCursor(ctx context.Context, pageSize int, pageToken string) ([]sharedModels.Activity, string, error) {
	log := s.logger.WithContext(ctx)
	defer utils.TimeOperation(log, "FirestoreService.ListAllCursor")()
	if s.client == nil {
		return nil, "", fmt.Errorf("firestore client is nil")
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	q := s.client.Collection(s.collection).OrderBy("created_at", firestorex.Desc)
	if pageToken != "" {
		if t, err := decodeToken(pageToken); err == nil && !t.IsZero() {
			q = q.StartAfter(t)
		}
	}
	q = q.Limit(pageSize + 1)
	it := q.Documents(ctx)
	defer it.Stop()
	var items []sharedModels.Activity
	for {
		doc, err := it.Next()
		if isDone(err) {
			break
		}
		if err != nil {
			return nil, "", err
		}
		var a struct {
			ID          int64       `firestore:"id"`
			Description string      `firestore:"description"`
			EmployeeID  int64       `firestore:"employee_id"`
			UrgencyID   int64       `firestore:"urgency_id"`
			CreatedAt   interface{} `firestore:"created_at"`
			UpdatedAt   interface{} `firestore:"updated_at"`
		}
		if err := doc.DataTo(&a); err != nil {
			continue
		}
		items = append(items, sharedModels.Activity{
			ID: uint(a.ID), Description: a.Description,
			EmployeeID: uint(a.EmployeeID), UrgencyID: uint(a.UrgencyID),
			CreatedAt: coerceTime(a.CreatedAt), UpdatedAt: coerceTime(a.UpdatedAt),
		})
	}
	var next string
	if len(items) > pageSize {
		last := items[pageSize-1]
		next = encodeToken(last.CreatedAt)
		items = items[:pageSize]
	}
	return items, next, nil
}

func coerceTime(v interface{}) time.Time {
	switch t := v.(type) {
	case time.Time:
		return t
	case string:
		if parsed, err := time.Parse(time.RFC3339, t); err == nil {
			return parsed
		}
		if parsed, err := time.Parse("2006-01-02T15:04:05Z07:00", t); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

func isDone(err error) bool { return firestorex.IsDone(err) }
