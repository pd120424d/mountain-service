package events

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
)

// Parse tries multiple strategies to extract an ActivityEvent from raw Pub/Sub bytes
// Returns the parsed and normalized event, the strategy label, or an error
func Parse(data []byte, attrs map[string]string) (activityV1.ActivityEvent, string, error) {
	agg := attrs["aggregateId"]
	// 1) Envelope OutboxEvent with eventData JSON (possibly double-encoded)
	if ev, ok := tryEnvelope(data, agg); ok {
		return *ev, "envelope", nil
	}

	// 2) Legacy direct ActivityEvent
	if ev, ok := tryLegacy(data, agg, "legacy"); ok {
		return ev, "legacy", nil
	}

	// 3) Whole message is quoted JSON string -> unquote then retry
	var msgAsString string
	if err := json.Unmarshal(data, &msgAsString); err == nil && msgAsString != "" {
		raw := []byte(msgAsString)
		if ev, ok := tryEnvelope(raw, agg); ok {
			return *ev, "quoted-envelope", nil
		}
		if ev, ok := tryLegacy(raw, agg, "quoted-legacy"); ok {
			return ev, "quoted-legacy", nil
		}
	}

	// 4) Base64 decoding then retry strategies
	if decoded, decErr := base64.StdEncoding.DecodeString(string(data)); decErr == nil && len(decoded) > 0 {
		if ev, ok := tryEnvelope(decoded, agg); ok {
			return *ev, "b64-envelope", nil
		}
		if ev, ok := tryLegacy(decoded, agg, "b64-legacy"); ok {
			return ev, "b64-legacy", nil
		}
		var s2 string
		if err := json.Unmarshal(decoded, &s2); err == nil && s2 != "" {
			raw2 := []byte(s2)
			if ev, ok := tryEnvelope(raw2, agg); ok {
				return *ev, "b64-quoted-envelope", nil
			}
			if ev, ok := tryLegacy(raw2, agg, "b64-quoted-legacy"); ok {
				return ev, "b64-quoted-legacy", nil
			}
		}
	}
	// 5) Base64 Raw (no padding) decoding then retry strategies
	if decoded, decErr := base64.RawStdEncoding.DecodeString(string(data)); decErr == nil && len(decoded) > 0 {
		if ev, ok := tryEnvelope(decoded, agg); ok {
			return *ev, "b64raw-envelope", nil
		}
		if ev, ok := tryLegacy(decoded, agg, "b64raw-legacy"); ok {
			return ev, "b64raw-legacy", nil
		}
		var s2 string
		if err := json.Unmarshal(decoded, &s2); err == nil && s2 != "" {
			raw2 := []byte(s2)
			if ev, ok := tryEnvelope(raw2, agg); ok {
				return *ev, "b64raw-quoted-envelope", nil
			}
			if ev, ok := tryLegacy(raw2, agg, "b64raw-quoted-legacy"); ok {
				return ev, "b64raw-quoted-legacy", nil
			}
		}
	}

	return activityV1.ActivityEvent{}, "", fmt.Errorf("unrecognized payload format")
}

func tryEnvelope(data []byte, aggAttr string) (*activityV1.ActivityEvent, bool) {
	var env activityV1.OutboxEvent
	if err := json.Unmarshal(data, &env); err != nil || (env.EventData == "" && env.EventType == "") {
		return nil, false
	}
	agg := env.AggregateID
	if agg == "" {
		agg = aggAttr
	}
	// First attempt: direct eventData
	if ev, err := env.GetEventData(); err == nil {
		postProcess(ev, agg)
		return ev, true
	}
	// Second: eventData may be quoted JSON string
	if len(env.EventData) > 0 && env.EventData[0] == '"' {
		if unq, uErr := strconv.Unquote(env.EventData); uErr == nil {
			var ev activityV1.ActivityEvent
			if jErr := json.Unmarshal([]byte(unq), &ev); jErr == nil {
				postProcess(&ev, env.AggregateID)
				return &ev, true
			}
		}
	}
	return nil, false
}

func tryLegacy(data []byte, aggAttr string, _ string) (activityV1.ActivityEvent, bool) {
	var legacy activityV1.ActivityEvent
	if err := json.Unmarshal(data, &legacy); err == nil {
		postProcess(&legacy, aggAttr)
		if legacy.ActivityID != 0 { // only accept if valid
			return legacy, true
		}
	}
	return activityV1.ActivityEvent{}, false
}

func postProcess(ev *activityV1.ActivityEvent, _ string) {
	if ev == nil {
		return
	}
	// ID recovery from aggregateId is DISABLED by request. Only normalize type here.
	t := strings.ToUpper(ev.Type)
	switch t {
	case "ACTIVITY.CREATED", "CREATED":
		t = "CREATE"
	case "ACTIVITY.UPDATED", "UPDATED":
		t = "UPDATE"
	case "ACTIVITY.DELETED", "DELETED":
		t = "DELETE"
	}
	ev.Type = t
}
