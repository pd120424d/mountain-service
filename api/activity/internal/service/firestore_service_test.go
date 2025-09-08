package service

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/firestoretest"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestFirestoreService_ListByUrgency(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()
	fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
		{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": "2025-01-02T10:00:00Z"},
		{"id": int64(2), "urgency_id": int64(3), "employee_id": int64(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
		{"id": int64(3), "urgency_id": int64(2), "employee_id": int64(7), "description": "C", "created_at": "2025-01-04T10:00:00Z"},
	})
	svc := NewFirebaseReadService(fake, log)

	items, err := svc.ListByUrgency(context.Background(), 2, 10)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestFirestoreService_ListAll(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()
	fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
		{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": "2025-01-02T10:00:00Z"},
		{"id": int64(2), "urgency_id": int64(3), "employee_id": int64(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
		{"id": int64(3), "urgency_id": int64(2), "employee_id": int64(7), "description": "C", "created_at": "2025-01-04T10:00:00Z"},
	})
	svc := NewFirebaseReadService(fake, log)

	items, err := svc.ListAll(context.Background(), 2)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestFirestoreService_ListByUrgencyCursor(t *testing.T) {
	log := utils.NewTestLogger()
	t1, _ := time.Parse(time.RFC3339, "2025-01-04T10:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2025-01-03T10:00:00Z")
	t3, _ := time.Parse(time.RFC3339, "2025-01-02T10:00:00Z")
	fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
		{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": t1},
		{"id": int64(2), "urgency_id": int64(2), "employee_id": int64(6), "description": "B", "created_at": t2},
		{"id": int64(3), "urgency_id": int64(2), "employee_id": int64(7), "description": "C", "created_at": t3},
	})
	svc := NewFirebaseReadService(fake, log)

	manualToken := base64.StdEncoding.EncodeToString([]byte("{\"createdAt\":\"" + t1.UTC().Format(time.RFC3339) + "\"}"))
	items, next, err := svc.ListByUrgencyCursor(context.Background(), 2, 2, manualToken)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Empty(t, next)
}

func TestFirestoreService_ListAllCursor(t *testing.T) {
	log := utils.NewTestLogger()
	t1, _ := time.Parse(time.RFC3339, "2025-01-04T10:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2025-01-03T10:00:00Z")
	t3, _ := time.Parse(time.RFC3339, "2025-01-02T10:00:00Z")
	fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
		{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": t1},
		{"id": int64(2), "urgency_id": int64(3), "employee_id": int64(6), "description": "B", "created_at": t2},
		{"id": int64(3), "urgency_id": int64(2), "employee_id": int64(7), "description": "C", "created_at": t3},
	})
	svc := NewFirebaseReadService(fake, log)

	manualToken := base64.StdEncoding.EncodeToString([]byte("{\"createdAt\":\"" + t1.UTC().Format(time.RFC3339) + "\"}"))
	items, next, err := svc.ListAllCursor(context.Background(), 2, manualToken)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Empty(t, next)
}

func TestEncodeDecodeToken(t *testing.T) {
	zero := time.Time{}
	assert.Equal(t, "", encodeToken(zero, 0))

	nonZero := time.Date(2025, 1, 5, 12, 0, 0, 0, time.UTC)
	tok := encodeToken(nonZero, 0)
	assert.NotEmpty(t, tok)

	tDecoded, id, err := decodeToken(tok)
	assert.NoError(t, err)
	assert.True(t, tDecoded.Equal(nonZero))
	assert.Equal(t, uint(0), id)

	_, _, err = decodeToken("???")
	assert.Error(t, err)

	tDecoded2, id2, err := decodeToken("")
	assert.NoError(t, err)
	assert.True(t, tDecoded2.IsZero())
	assert.Equal(t, uint(0), id2)
}

func TestCoerceTime(t *testing.T) {
	now := time.Date(2025, 2, 1, 8, 30, 0, 0, time.UTC)
	assert.True(t, coerceTime(now).Equal(now))

	str := "2025-02-01T08:30:00Z"
	ct := coerceTime(str)
	assert.Equal(t, str, ct.UTC().Format(time.RFC3339))

	invalid := coerceTime("not-a-time")
	assert.True(t, invalid.IsZero())
}

func TestFirestoreService_ListByUrgencyCursor_NextToken(t *testing.T) {
	log := utils.NewTestLogger()
	fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
		{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": "2025-01-04T10:00:00Z"},
		{"id": int64(2), "urgency_id": int64(2), "employee_id": int64(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
		{"id": int64(3), "urgency_id": int64(2), "employee_id": int64(7), "description": "C", "created_at": "2025-01-02T10:00:00Z"},
	})
	svc := NewFirebaseReadService(fake, log)

	items, next, err := svc.ListByUrgencyCursor(context.Background(), 2, 2, "")
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.NotEmpty(t, next)
}

func TestFirestoreService_ListAllCursor_InvalidToken(t *testing.T) {
	log := utils.NewTestLogger()
	fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
		{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": "2025-01-04T10:00:00Z"},
		{"id": int64(2), "urgency_id": int64(3), "employee_id": int64(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
	})
	svc := NewFirebaseReadService(fake, log)

	items, next, err := svc.ListAllCursor(context.Background(), 1, "!!!")
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	// token is ignored when invalid, so we still get a next token when more results exist
	assert.NotEmpty(t, next)
}

func TestFirestoreService_ListByUrgency_ClientNil(t *testing.T) {
	svc := NewFirebaseReadService(nil, utils.NewTestLogger())
	_, err := svc.ListByUrgency(context.Background(), 1, 10)
	assert.Error(t, err)
}

func TestFirestoreService_ListAll_ClientNil(t *testing.T) {
	svc := NewFirebaseReadService(nil, utils.NewTestLogger())
	_, err := svc.ListAll(context.Background(), 10)
	assert.Error(t, err)
}
