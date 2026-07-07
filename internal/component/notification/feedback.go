package notification

import (
	"context"

	auditcore "github.com/neteast-software/go-module/audit/operate"
	eventcore "github.com/neteast-software/go-module/fault/event"
)

type feedback struct {
	auditAction   string
	auditResource string
	eventKind     string
	eventMessage  string
	context       map[string]any
}

func (p *Component) recordFeedback(ctx context.Context, item feedback) error {
	_ = p.audit.Record(ctx, auditcore.New(item.auditAction,
		auditcore.WithTransport(auditcore.Bin),
		auditcore.WithApplication("app2"),
		auditcore.WithResource(item.auditResource),
		auditcore.WithContext(item.context),
	))
	return p.event.Record(ctx, eventcore.New(eventcore.Info, item.eventMessage,
		eventcore.WithSource(ID.String()),
		eventcore.WithKind(item.eventKind),
		eventcore.WithContext(item.context),
	))
}
