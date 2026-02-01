package observer

import (
	"context"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
	"go.uber.org/zap"
)

type AuditPublisher struct {
	log   *zap.SugaredLogger
	sinks []Observer
}

func NewAuditPublisher(log *zap.SugaredLogger) *AuditPublisher {
	return &AuditPublisher{
		log: log,
	}
}
func (a AuditPublisher) Register(o Observer) {
	a.sinks = append(a.sinks, o)
}

func (a AuditPublisher) NotifyObservers(ctx context.Context, ev models.AuditEvent) {
	for _, sink := range a.sinks {
		err := sink.Consume(ctx, ev)
		if err != nil {
			a.log.Errorw("failed to consume event", "event", ev, "err", err)
		}
	}
}
