package events

import (
	"encoding/json"
	"testing"
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
)

func evt(t string) activityV1.ActivityEvent {
	return activityV1.ActivityEvent{Type: t, ActivityID: 101, UrgencyID: 2, EmployeeID: 3, Description: "x", CreatedAt: time.Now().UTC()}
}

func Test_Parse_succeeds_when_parsing_envelope_with_json_eventData(t *testing.T) {
	e := evt("CREATE")
	o := activityV1.CreateOutboxEvent(activityV1.ActivityEventCreated, 101, e)
	b, _ := json.Marshal(o)
	got, strat, err := Parse(b, nil)
	if err != nil || got.ActivityID != 101 || got.Type != "CREATE" || strat != "envelope" {
		t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
	}
}

func Test_Parse_succeeds_when_parsing_legacy_event(t *testing.T) {
	e := evt("CREATE")
	b, _ := json.Marshal(e)
	got, strat, err := Parse(b, nil)
	if err != nil || got.ActivityID != 101 || got.Type != "CREATE" || strat != "legacy" {
		t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
	}
}

func Test_Parse_succeeds_when_parsing_envelope_with_quoted_eventData(t *testing.T) {
	e := evt("CREATE")
	j, _ := json.Marshal(e)
	o := activityV1.OutboxEvent{EventType: string(activityV1.ActivityEventCreated), AggregateID: "activity-101", EventData: string(json.RawMessage(j))}
	b, _ := json.Marshal(o)
	got, strat, err := Parse(b, nil)
	if err != nil || got.ActivityID != 101 || got.Type != "CREATE" || strat != "envelope" {
		t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
	}
}

func Test_Parse_succeeds_when_parsing_quoted_message_payload(t *testing.T) {
	e := evt("CREATE")
	oj := activityV1.CreateOutboxEvent(activityV1.ActivityEventCreated, 101, e)
	b, _ := json.Marshal(oj)
	qb, _ := json.Marshal(string(b)) // quoted whole message
	got, strat, err := Parse(qb, nil)
	if err != nil || got.ActivityID != 101 || got.Type != "CREATE" || strat != "quoted-envelope" {
		t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
	}
}

func Test_Parse_succeeds_when_normalizing_activity_dot_created(t *testing.T) {
	e := evt("activity.created")
	b, _ := json.Marshal(e)
	got, strat, err := Parse(b, nil)
	if err != nil || got.Type != "CREATE" || strat != "legacy" {
		t.Fatalf("unexpected: err=%v strat=%s got=%+v", err, strat, got)
	}
}

