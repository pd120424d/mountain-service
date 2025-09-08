package service

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/firestoretest"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestFirestoreService_ListByUrgency(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when Firebase client is valid", func(t *testing.T) {
		log := utils.NewTestLogger()
		fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
			{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": "2025-01-02T10:00:00Z"},
			{"id": int64(2), "urgency_id": int64(3), "employee_id": int64(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
			{"id": int64(3), "urgency_id": int64(2), "employee_id": int64(7), "description": "C", "created_at": "2025-01-04T10:00:00Z"},
		})
		svc := NewFirebaseReadService(fake, log)
		items, err := svc.ListByUrgency(t.Context(), 2, 10)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
	})

	t.Run("it returns error when Firebase client is nil", func(t *testing.T) {
		svc := NewFirebaseReadService(nil, utils.NewTestLogger())
		_, err := svc.ListByUrgency(t.Context(), 1, 10)
		assert.Error(t, err)
	})
}

func TestFirestoreService_ListAll(t *testing.T) {
	t.Parallel()

	t.Run("it suceeds when Firebase client is valid", func(t *testing.T) {
		log := utils.NewTestLogger()
		fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
			{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": "2025-01-02T10:00:00Z"},
			{"id": int64(2), "urgency_id": int64(3), "employee_id": int64(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
			{"id": int64(3), "urgency_id": int64(2), "employee_id": int64(7), "description": "C", "created_at": "2025-01-04T10:00:00Z"},
		})
		svc := NewFirebaseReadService(fake, log)
		items, err := svc.ListAll(t.Context(), 2)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
	})

	t.Run("it returns error when Firebase client is nil", func(t *testing.T) {
		svc := NewFirebaseReadService(nil, utils.NewTestLogger())
		_, err := svc.ListAll(t.Context(), 10)
		assert.Error(t, err)
	})
}

func TestFirestoreService_ListByUrgencyCursor(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds with manual token", func(t *testing.T) {
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
		items, next, err := svc.ListByUrgencyCursor(t.Context(), 2, 2, manualToken)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Empty(t, next)
	})

	t.Run("it succeeds with two pages with distinct timestamps", func(t *testing.T) {
		log := utils.NewTestLogger()
		t1, _ := time.Parse(time.RFC3339, "2025-03-03T10:00:00Z")
		t2, _ := time.Parse(time.RFC3339, "2025-03-02T10:00:00Z")
		t3, _ := time.Parse(time.RFC3339, "2025-03-01T10:00:00Z")
		fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
			{"id": int64(3), "urgency_id": int64(9), "employee_id": int64(1), "description": "C", "created_at": t1},
			{"id": int64(2), "urgency_id": int64(9), "employee_id": int64(1), "description": "B", "created_at": t2},
			{"id": int64(1), "urgency_id": int64(9), "employee_id": int64(1), "description": "A", "created_at": t3},
		})
		svc := NewFirebaseReadService(fake, log)
		items1, next1, err := svc.ListByUrgencyCursor(t.Context(), 9, 2, "")
		assert.NoError(t, err)
		if assert.Len(t, items1, 2) {
			assert.Equal(t, uint(3), items1[0].ID)
			assert.Equal(t, uint(2), items1[1].ID)
		}
		assert.NotEmpty(t, next1)
		items2, next2, err := svc.ListByUrgencyCursor(t.Context(), 9, 2, next1)
		assert.NoError(t, err)
		assert.Len(t, items2, 1)
		assert.Equal(t, uint(1), items2[0].ID)
		assert.Empty(t, next2)
	})

	t.Run("it succeeds with next token from first page", func(t *testing.T) {
		log := utils.NewTestLogger()
		fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
			{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": "2025-01-04T10:00:00Z"},
			{"id": int64(2), "urgency_id": int64(2), "employee_id": int64(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
			{"id": int64(3), "urgency_id": int64(2), "employee_id": int64(7), "description": "C", "created_at": "2025-01-02T10:00:00Z"},
		})
		svc := NewFirebaseReadService(fake, log)
		items, next, err := svc.ListByUrgencyCursor(t.Context(), 2, 2, "")
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.NotEmpty(t, next)
	})

	t.Run("it falls back on invalid token", func(t *testing.T) {
		log := utils.NewTestLogger()
		fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
			{"id": int64(2), "urgency_id": int64(5), "employee_id": int64(1), "description": "B", "created_at": "2025-01-02T10:00:00Z"},
			{"id": int64(1), "urgency_id": int64(5), "employee_id": int64(1), "description": "A", "created_at": "2025-01-01T10:00:00Z"},
		})
		svc := NewFirebaseReadService(fake, log)
		items, next, err := svc.ListByUrgencyCursor(t.Context(), 5, 1, "!!!")
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.NotEmpty(t, next)
	})

	t.Run("it returns zero items for non-empty token beyond range", func(t *testing.T) {
		log := utils.NewTestLogger()
		fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
			{"id": int64(2), "urgency_id": int64(6), "employee_id": int64(1), "description": "B", "created_at": "2025-01-02T10:00:00Z"},
			{"id": int64(1), "urgency_id": int64(6), "employee_id": int64(1), "description": "A", "created_at": "2025-01-01T10:00:00Z"},
		})
		svc := NewFirebaseReadService(fake, log)
		oldTok := base64.RawURLEncoding.EncodeToString([]byte("{\"createdAt\":\"0001-01-01T00:00:00Z\"}"))
		items, next, err := svc.ListByUrgencyCursor(t.Context(), 6, 5, oldTok)
		assert.NoError(t, err)
		assert.Len(t, items, 0)
		assert.Empty(t, next)
	})
}

func TestFirestoreService_ListAllCursor(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds with manual token", func(t *testing.T) {
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
		items, next, err := svc.ListAllCursor(t.Context(), 2, manualToken)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Empty(t, next)
	})

	t.Run("it succeeds with invalid token fallback", func(t *testing.T) {
		log := utils.NewTestLogger()
		fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
			{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": "2025-01-04T10:00:00Z"},
			{"id": int64(2), "urgency_id": int64(3), "employee_id": int64(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
		})
		svc := NewFirebaseReadService(fake, log)
		items, next, err := svc.ListAllCursor(t.Context(), 1, "!!!")
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.NotEmpty(t, next)
	})

	t.Run("it succeeds with page size capping and URL safe token", func(t *testing.T) {
		log := utils.NewTestLogger()
		base, _ := time.Parse(time.RFC3339, "2025-04-01T00:00:00Z")
		docs := make([]map[string]interface{}, 0, 105)
		for i := 0; i < 105; i++ {
			id := int64(200 - i)
			docs = append(docs, map[string]interface{}{
				"id": id, "urgency_id": int64(1), "employee_id": int64(1), "description": "x",
				"created_at": base.Add(-time.Duration(i) * time.Minute),
			})
		}
		fake := firestoretest.NewFake().WithCollection("activities", docs)
		svc := NewFirebaseReadService(fake, log)
		items, next, err := svc.ListAllCursor(t.Context(), 1000, "")
		assert.NoError(t, err)
		assert.Len(t, items, 100)
		assert.NotEmpty(t, next)
		assert.NotContains(t, next, "+")
		assert.NotContains(t, next, "/")
		assert.NotContains(t, next, "=")
		items2, next2, err := svc.ListAllCursor(t.Context(), 1000, next)
		assert.NoError(t, err)
		assert.Len(t, items2, 5)
		assert.Empty(t, next2)
	})

	t.Run("it succeeds with default page size when zero", func(t *testing.T) {
		log := utils.NewTestLogger()
		fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
			{"id": int64(3), "urgency_id": int64(1), "employee_id": int64(1), "description": "c", "created_at": "2025-01-03T00:00:00Z"},
			{"id": int64(2), "urgency_id": int64(1), "employee_id": int64(1), "description": "b", "created_at": "2025-01-02T00:00:00Z"},
			{"id": int64(1), "urgency_id": int64(1), "employee_id": int64(1), "description": "a", "created_at": "2025-01-01T00:00:00Z"},
		})
		svc := NewFirebaseReadService(fake, log)
		items, _, err := svc.ListAllCursor(t.Context(), 0, "")
		assert.NoError(t, err)
		assert.Len(t, items, 3)
	})

	t.Run("it returns zero items for non-empty token beyond range", func(t *testing.T) {
		log := utils.NewTestLogger()
		fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
			{"id": int64(2), "urgency_id": int64(1), "employee_id": int64(1), "description": "B", "created_at": "2025-01-02T00:00:00Z"},
			{"id": int64(1), "urgency_id": int64(1), "employee_id": int64(1), "description": "A", "created_at": "2025-01-01T00:00:00Z"},
		})
		svc := NewFirebaseReadService(fake, log)
		oldTok := base64.RawURLEncoding.EncodeToString([]byte("{\"createdAt\":\"0001-01-01T00:00:00Z\"}"))
		items, next, err := svc.ListAllCursor(t.Context(), 5, oldTok)
		assert.NoError(t, err)
		assert.Len(t, items, 0)
		assert.Empty(t, next)
	})

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

	// roundtrip with non-zero id
	tok2 := encodeToken(nonZero, 123)
	tDecoded3, id3, err := decodeToken(tok2)
	assert.NoError(t, err)
	assert.True(t, tDecoded3.Equal(nonZero))
	assert.Equal(t, uint(123), id3)

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

	// also supports full-offset format
	str2 := "2025-02-01T08:30:00+00:00"
	ct2 := coerceTime(str2)
	assert.Equal(t, "2025-02-01T08:30:00Z", ct2.UTC().Format(time.RFC3339))

	invalid := coerceTime("not-a-time")
	assert.True(t, invalid.IsZero())
}
