package data

import (
	"context"
	"neelbhat88/nest-monitor/m/v2/internal/models"
)

type Controller interface {
	WriteEventMessage(ctx context.Context, msg []byte) error
	WriteHvacEvent(ctx context.Context, event models.HVACEvent) error
}
