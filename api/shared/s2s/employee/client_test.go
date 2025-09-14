package employee

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type stubHTTP struct{ responses []*http.Response; errs []error; idx int }

func (s *stubHTTP) Get(ctx context.Context, endpoint string) (*http.Response, error) {
	if s.idx < len(s.errs) && s.errs[s.idx] != nil {
		err := s.errs[s.idx]
		s.idx++
		return nil, err
	}
	if s.idx < len(s.responses) {
		r := s.responses[s.idx]
		s.idx++
		return r, nil
	}
	return &http.Response{StatusCode: 500, Body: ioutil.NopCloser(bytes.NewReader(nil))}, nil
}

func jsonBody(v interface{}) *http.Response {
	b, _ := json.Marshal(v)
	return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader(b))}
}

func status(code int) *http.Response { return &http.Response{StatusCode: code, Body: ioutil.NopCloser(bytes.NewReader(nil))} }

func TestEmployeeClient_GetEmployeeByID(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()

	t.Run("ok", func(t *testing.T) {
		stub := &stubHTTP{responses: []*http.Response{jsonBody(employeeV1.EmployeeResponse{ID: 7, FirstName: "A", LastName: "B"})}}
		c := &clientImpl{http: stub, logger: log, maxRetries: 1}
		got, err := c.GetEmployeeByID(t.Context(), 7)
		if err != nil { t.Fatalf("err: %v", err) }
		if got == nil || got.ID != 7 || got.FirstName != "A" { t.Fatalf("bad resp: %+v", got) }
	})

	t.Run("not_found", func(t *testing.T) {
		stub := &stubHTTP{responses: []*http.Response{status(404)}}
		c := &clientImpl{http: stub, logger: log, maxRetries: 0}
		_, err := c.GetEmployeeByID(t.Context(), 9)
		if err == nil { t.Fatalf("expected error") }
	})

	t.Run("retry_then_ok", func(t *testing.T) {
		stub := &stubHTTP{responses: []*http.Response{status(500), jsonBody(employeeV1.EmployeeResponse{ID: 3})}}
		c := &clientImpl{http: stub, logger: log, maxRetries: 2}
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		got, err := c.GetEmployeeByID(ctx, 3)
		if err != nil { t.Fatalf("err: %v", err) }
		if got == nil || got.ID != 3 { t.Fatalf("bad resp: %+v", got) }
	})
}

func TestEmployeeClient_Collections(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()

	t.Run("get_all", func(t *testing.T) {
		stub := &stubHTTP{responses: []*http.Response{jsonBody(employeeV1.AllEmployeesResponse{Employees: []employeeV1.EmployeeResponse{{ID: 1}, {ID: 2}}})}}
		c := &clientImpl{http: stub, logger: log, maxRetries: 1}
		list, err := c.GetAllEmployees(t.Context())
		if err != nil { t.Fatalf("err: %v", err) }
		if len(list) != 2 { t.Fatalf("want 2, got %d", len(list)) }
	})

	t.Run("on_call", func(t *testing.T) {
		stub := &stubHTTP{responses: []*http.Response{jsonBody(employeeV1.OnCallEmployeesResponse{Employees: []employeeV1.EmployeeResponse{{ID: 5}}})}}
		c := &clientImpl{http: stub, logger: log, maxRetries: 1}
		list, err := c.GetOnCallEmployees(t.Context(), 30*time.Minute)
		if err != nil { t.Fatalf("err: %v", err) }
		if len(list) != 1 || list[0].ID != 5 { t.Fatalf("unexpected: %+v", list) }
	})

	t.Run("active_emergencies", func(t *testing.T) {
		stub := &stubHTTP{responses: []*http.Response{jsonBody(employeeV1.ActiveEmergenciesResponse{HasActiveEmergencies: true})}}
		c := &clientImpl{http: stub, logger: log, maxRetries: 1}
		has, err := c.CheckActiveEmergencies(t.Context(), 11)
		if err != nil { t.Fatalf("err: %v", err) }
		if !has { t.Fatalf("expected true") }
	})
}

