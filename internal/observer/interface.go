package observer

import (
	"context"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

type Observer interface {
	Consume(ctx context.Context, ev models.AuditEvent) error
	Close() error
}

type Publisher interface {
	Register(o Observer)
	NotifyObservers(ctx context.Context, ev models.AuditEvent)
}
