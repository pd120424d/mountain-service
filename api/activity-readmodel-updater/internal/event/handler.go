package events

//go:generate mockgen -source=handler.go -destination=handler_gomock.go -package=events github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/event -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/service"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type EventHandler interface {
	Handle(ctx context.Context, msg *pubsub.Message) error
}

type Handler struct {
	fb     service.FirebaseService
	logger utils.Logger
}

func NewHandler(fb service.FirebaseService, logger utils.Logger) *Handler {
	return &Handler{fb: fb, logger: logger.WithName("eventHandler")}
}

func (h *Handler) Handle(ctx context.Context, msg *pubsub.Message) error {
	h.logger.Infof("Received activity event: message_id=%s, publish_time=%v", msg.ID, msg.PublishTime)

	ev, strat, err := Parse(msg.Data, msg.Attributes)
	if err != nil {
		h.logger.Errorf("Unrecognized event payload format, cannot parse message_id=%s", msg.ID)
		return fmt.Errorf("unrecognized event payload format")
	}
	h.logger.Infof("Parsed activity event via strategy=%s", strat)

	// Normalize type and process
	normalizeType(&ev)
	return h.fb.SyncActivity(ctx, ev)
}

func normalizeType(ev *activityV1.ActivityEvent) {
	if ev == nil {
		return
	}
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
