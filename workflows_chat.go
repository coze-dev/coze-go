package coze

import (
	"context"
	"net/http"
)

func (r *workflowsChat) Stream(ctx context.Context, req *WorkflowsChatStreamReq) (Stream[ChatEvent], error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/workflows/chat",
		Body:   req,
	}
	response := new(createChatsResp)
	err := r.client.rawRequest(ctx, request, response)
	return newStream(ctx, r.client, response.HTTPResponse, parseChatEvent), err
}

// WorkflowsChatStreamReq 表示工作流聊天流式请求
type WorkflowsChatStreamReq struct {
	WorkflowID         string            `json:"workflow_id"`               // 工作流ID
	AdditionalMessages []*Message        `json:"additional_messages"`       // 额外的消息信息
	Parameters         map[string]any    `json:"parameters,omitempty"`      // 工作流参数
	AppID              *string           `json:"app_id,omitempty"`          // 应用ID
	BotID              *string           `json:"bot_id,omitempty"`          // 机器人ID
	ConversationID     *string           `json:"conversation_id,omitempty"` // 会话ID
	Ext                map[string]string `json:"ext,omitempty"`             // 扩展信息
}

type workflowsChat struct {
	client *core
}

func newWorkflowsChat(core *core) *workflowsChat {
	return &workflowsChat{
		client: core,
	}
}
