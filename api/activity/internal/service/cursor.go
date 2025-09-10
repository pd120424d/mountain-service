package service

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type cursorToken struct {
	CreatedAt string `json:"createdAt"`
	ID        uint   `json:"id,omitempty"`
}

func encodeToken(t time.Time, id uint) string {
	if t.IsZero() {
		return ""
	}
	b, _ := json.Marshal(cursorToken{CreatedAt: t.UTC().Format(time.RFC3339Nano), ID: id})
	return base64.RawURLEncoding.EncodeToString(b)
}

func decodeToken(token string) (time.Time, uint, error) {
	if token == "" {
		return time.Time{}, 0, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return time.Time{}, 0, err
	}
	var ct cursorToken
	if err := json.Unmarshal(raw, &ct); err != nil {
		return time.Time{}, 0, err
	}
	if ct.CreatedAt == "" {
		return time.Time{}, ct.ID, nil
	}
	t, err := time.Parse(time.RFC3339Nano, ct.CreatedAt)
	if err != nil {
		if t2, err2 := time.Parse(time.RFC3339, ct.CreatedAt); err2 == nil {
			return t2, ct.ID, nil
		}
		return time.Time{}, ct.ID, err
	}
	return t, ct.ID, nil
}
