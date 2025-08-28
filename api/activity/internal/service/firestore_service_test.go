package service

import (
	"context"
	"testing"

	"github.com/pd120424d/mountain-service/api/shared/firestoretest"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestFirestoreService_ListByUrgency(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()
	fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
		{"id": uint(1), "urgency_id": uint(2), "employee_id": uint(5), "description": "A", "created_at": "2025-01-02T10:00:00Z"},
		{"id": uint(2), "urgency_id": uint(3), "employee_id": uint(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
		{"id": uint(3), "urgency_id": uint(2), "employee_id": uint(7), "description": "C", "created_at": "2025-01-04T10:00:00Z"},
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
		{"id": uint(1), "urgency_id": uint(2), "employee_id": uint(5), "description": "A", "created_at": "2025-01-02T10:00:00Z"},
		{"id": uint(2), "urgency_id": uint(3), "employee_id": uint(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
		{"id": uint(3), "urgency_id": uint(2), "employee_id": uint(7), "description": "C", "created_at": "2025-01-04T10:00:00Z"},
	})
	svc := NewFirebaseReadService(fake, log)

	items, err := svc.ListAll(context.Background(), 2)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
}
