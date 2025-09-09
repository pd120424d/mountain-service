package events

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/stretchr/testify/assert"
)

func evt(t string) activityV1.ActivityEvent {
	return activityV1.ActivityEvent{Type: t, ActivityID: 101, UrgencyID: 2, EmployeeID: 3, Description: "x", CreatedAt: time.Now().UTC()}
}

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when parsing envelope with json eventData", func(t *testing.T) {
		e := evt("CREATE")
		o := activityV1.CreateOutboxEvent(101, e)
		b, _ := json.Marshal(o)
		got, strat, err := Parse(b, nil)
		if err != nil || got.ActivityID != uint(101) || got.Type != "CREATE" || strat != "envelope" {
			t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
		}
	})

	t.Run("it succeeds when parsing legacy event", func(t *testing.T) {
		e := evt("CREATE")
		b, _ := json.Marshal(e)
		got, strat, err := Parse(b, nil)
		if err != nil || got.ActivityID != uint(101) || got.Type != "CREATE" || strat != "legacy" {
			t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
		}
	})

	t.Run("it succeeds when parsing envelope with quoted eventData", func(t *testing.T) {
		e := evt("CREATE")
		o := activityV1.CreateOutboxEvent(101, e)
		q := strconv.Quote(o.EventData)
		o.EventData = q
		b, _ := json.Marshal(o)
		got, strat, err := Parse(b, nil)
		if err != nil || got.ActivityID != uint(101) || got.Type != "CREATE" || strat != "envelope" {
			t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
		}
	})

	t.Run("it succeeds when parsing legacy event with quoted whole message", func(t *testing.T) {
		e := evt("CREATE")
		b, _ := json.Marshal(e)
		qb, _ := json.Marshal(string(b))
		got, strat, err := Parse(qb, nil)
		if err != nil || got.ActivityID != uint(101) || got.Type != "CREATE" || strat != "quoted-legacy" {
			t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
		}
	})

	t.Run("it succeeds when parsing base64 encoded envelope", func(t *testing.T) {
		e := evt("CREATE")
		o := activityV1.CreateOutboxEvent(101, e)
		b, _ := json.Marshal(o)
		enc := base64.StdEncoding.EncodeToString(b)
		got, strat, err := Parse([]byte(enc), nil)
		if err != nil || got.ActivityID != uint(101) || got.Type != "CREATE" || strat != "b64-envelope" {
			t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
		}
	})

	t.Run("it succeeds when parsing base64 encoded legacy event", func(t *testing.T) {
		e := evt("CREATE")
		b, _ := json.Marshal(e)
		enc := base64.StdEncoding.EncodeToString(b)
		got, strat, err := Parse([]byte(enc), nil)
		if err != nil || got.ActivityID != uint(101) || got.Type != "CREATE" || strat != "b64-legacy" {
			t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
		}
	})

	// base64-quoted envelope (accept any b64* strategy)
	t.Run("it succeeds when parsing base64 quoted envelope", func(t *testing.T) {
		e := evt("CREATE")
		o := activityV1.CreateOutboxEvent(101, e)
		b, _ := json.Marshal(o)
		qb, _ := json.Marshal(string(b))
		enc := base64.StdEncoding.EncodeToString(qb)
		got, strat, err := Parse([]byte(enc), nil)
		if err != nil || got.ActivityID != uint(101) || !strings.HasPrefix(strat, "b64") {
			t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
		}
	})

	// base64raw-envelope (accept any b64* strategy)
	t.Run("it succeeds when parsing base64raw envelope", func(t *testing.T) {
		e := evt("CREATE")
		o := activityV1.CreateOutboxEvent(101, e)
		b, _ := json.Marshal(o)
		rawEnc := base64.RawStdEncoding.EncodeToString(b)
		got2, strat2, err2 := Parse([]byte(rawEnc), nil)
		if err2 != nil || got2.ActivityID != uint(101) || !strings.HasPrefix(strat2, "b64") {
			t.Fatalf("unexpected: err=%v strat=%s got=%+v", err2, strat2, got2)
		}
	})

	// base64raw-quoted-legacy (accept any b64* strategy)
	t.Run("it succeeds when parsing base64raw quoted legacy", func(t *testing.T) {
		e := evt("CREATE")
		lb, _ := json.Marshal(e)
		lqb, _ := json.Marshal(string(lb))
		rawEnc2 := base64.RawStdEncoding.EncodeToString(lqb)
		got3, strat3, err3 := Parse([]byte(rawEnc2), nil)
		if err3 != nil || got3.ActivityID != uint(101) || !strings.HasPrefix(strat3, "b64") {
			t.Fatalf("unexpected: err=%v strat=%s got=%+v", err3, strat3, got3)
		}
	})

	// failure on invalid payload
	t.Run("it fails on invalid payload", func(t *testing.T) {
		_, _, err := Parse([]byte("notjson"), nil)
		if err == nil {
			t.Fatalf("expected error for invalid payload")
		}
	})

	// type normalization
	t.Run("it normalizes event type", func(t *testing.T) {
		e := evt("activity.updated")
		b, _ := json.Marshal(e)
		got, _, err := Parse(b, nil)
		if err != nil || got.Type != "UPDATE" {
			t.Fatalf("expected UPDATE, got: %s (err=%v)", got.Type, err)
		}
	})
}

func TestParser_tryEnvelope(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when parsing envelope with json eventData", func(t *testing.T) {
		e := evt("CREATE")
		o := activityV1.CreateOutboxEvent(101, e)
		b, _ := json.Marshal(o)
		got, ok := tryEnvelope(b, "")
		assert.True(t, ok)
		assert.Equal(t, uint(101), got.ActivityID)
		assert.Equal(t, "CREATE", got.Type)
	})

	t.Run("it succeeds when parsing envelope with quoted eventData", func(t *testing.T) {
		e := evt("CREATE")
		o := activityV1.CreateOutboxEvent(101, e)
		o.EventData = strconv.Quote(o.EventData)
		b, _ := json.Marshal(o)
		got, ok := tryEnvelope(b, "")

		assert.True(t, ok)
		assert.Equal(t, uint(101), got.ActivityID)
		assert.Equal(t, "CREATE", got.Type)
	})
}

func Test_tryLegacy_rejects_zero_id(t *testing.T) {
	le := activityV1.ActivityEvent{Type: "CREATE", ActivityID: 0}
	b, _ := json.Marshal(le)
	_, ok := tryLegacy(b, "", "legacy")
	assert.False(t, ok)
}
