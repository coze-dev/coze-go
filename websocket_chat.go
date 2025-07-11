package coze

import (
	"context"
	"strconv"
)

func (c *websocketChatBuilder) Create(ctx context.Context, req *CreateWebsocketChatReq) *WebSocketChat {
	return newWebsocketChatClient(ctx, c.core, req)
}

type CreateWebsocketChatReq struct {
	WebSocketClientOption *WebSocketClientOption

	// BotID is the ID of the bot.
	BotID *string `json:"bot_id"`
	// WorkflowID is the ID of the workflow.
	WorkflowID *string `json:"workflow_id"`
	// DeviceID is the ID of the device.
	DeviceID *int64 `json:"device_id"`
}

func (c *CreateWebsocketChatReq) toQuery() map[string]string {
	q := map[string]string{}
	if c.BotID != nil {
		q["bot_id"] = *c.BotID
	}
	if c.WorkflowID != nil {
		q["workflow_id"] = *c.WorkflowID
	}
	if c.DeviceID != nil {
		q["device_id"] = strconv.FormatInt(*c.DeviceID, 10)
	}
	return q
}

type websocketChatBuilder struct {
	core *core
}

func newWebsocketChat(core *core) *websocketChatBuilder {
	return &websocketChatBuilder{
		core: core,
	}
}
