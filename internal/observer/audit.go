package observer

import (
	"context"
	"fmt"

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
func (a *AuditPublisher) Register(o Observer) {
	fmt.Println("added")
	fmt.Println("obs: ", o)
	a.sinks = append(a.sinks, o)
}

func (a *AuditPublisher) NotifyObservers(ctx context.Context, ev models.AuditEvent) {
	fmt.Println("notifying observers")
	for _, sink := range a.sinks {
		fmt.Println("notifying sink", sink)
		fmt.Println("sending event", ev)
		err := sink.Consume(ctx, ev)
		if err != nil {
			a.log.Errorw("failed to consume event", "event", ev, "err", err)
		}
	}
}
