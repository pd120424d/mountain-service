package clients

import (
	"context"

	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
)

// UrgencyClient defines the methods activity service needs from urgency service.
type UrgencyClient interface {
	GetUrgencyByID(ctx context.Context, id uint) (*urgencyV1.UrgencyResponse, error)
}

