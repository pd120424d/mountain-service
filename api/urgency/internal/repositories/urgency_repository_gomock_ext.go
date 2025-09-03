package repositories

import (
	"context"
	"reflect"

	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"go.uber.org/mock/gomock"
)

// Extend the generated mock with the new method to avoid regenerating mocks right now.
func (m *MockUrgencyRepository) GetByIDPrimary(ctx context.Context, id uint, urgency *model.Urgency) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIDPrimary", ctx, id, urgency)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockUrgencyRepositoryMockRecorder) GetByIDPrimary(ctx, id, urgency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDPrimary", reflect.TypeOf((*MockUrgencyRepository)(nil).GetByIDPrimary), ctx, id, urgency)
}
