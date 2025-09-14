package urgency

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
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

func TestUrgencyClient_GetUrgencyByID(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()

	t.Run("ok", func(t *testing.T) {
		stub := &stubHTTP{responses: []*http.Response{jsonBody(urgencyV1.UrgencyResponse{ID: 42, Status: urgencyV1.Open})}}
		c := &clientImpl{http: stub, logger: log, maxRetries: 1}
		got, err := c.GetUrgencyByID(t.Context(), 42)
		if err != nil { t.Fatalf("err: %v", err) }
		if got == nil || got.ID != 42 { t.Fatalf("bad resp: %+v", got) }
	})

	t.Run("not_found", func(t *testing.T) {
		stub := &stubHTTP{responses: []*http.Response{status(404)}}
		c := &clientImpl{http: stub, logger: log, maxRetries: 0}
		_, err := c.GetUrgencyByID(t.Context(), 9)
		if err == nil { t.Fatalf("expected error") }
	})

	t.Run("retry_then_ok", func(t *testing.T) {
		stub := &stubHTTP{responses: []*http.Response{status(500), jsonBody(urgencyV1.UrgencyResponse{ID: 7, Status: urgencyV1.InProgress})}}
		c := &clientImpl{http: stub, logger: log, maxRetries: 2}
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		got, err := c.GetUrgencyByID(ctx, 7)
		if err != nil { t.Fatalf("err: %v", err) }
		if got == nil || got.ID != 7 { t.Fatalf("bad resp: %+v", got) }
	})
}

