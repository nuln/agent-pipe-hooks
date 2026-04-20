package hooks

import (
	"context"
	"log/slog"

	agent "github.com/nuln/agent-core"
)

func init() {
	agent.RegisterPluginConfigSpec(agent.PluginConfigSpec{
		PluginName:  "hooks",
		PluginType:  "pipe",
		Description: "Fires HookHandler events (message.received, message.sent) for registered hook handlers",
	})

	agent.RegisterPipe("hooks", 50, func(pctx agent.PipeContext) agent.Pipe {
		return &HooksPipe{}
	})
}

// HooksPipe fires hook events on messages passing through the pipeline.
// Hook handlers are expected to be registered directly via agent.RegisterHookHandler
// (future extension point).
type HooksPipe struct{}

func (h *HooksPipe) Handle(ctx context.Context, d agent.Dialog, msg *agent.Message) bool {
	evt := agent.HookEvent{
		Name:       "message.received",
		SessionKey: msg.SessionKey,
		UserID:     msg.UserID,
	}
	if err := fireGlobalHooks(ctx, d, msg, evt); err != nil {
		slog.Warn("hooks: error firing message.received", "error", err)
	}
	return false
}

// fireGlobalHooks is a no-op placeholder. Real implementations inject
// HookHandler instances via a registry (added in a future version).
func fireGlobalHooks(_ context.Context, _ agent.Dialog, _ *agent.Message, evt agent.HookEvent) error {
	slog.Debug("hooks: fired", "event", evt.Name)
	return nil
}
