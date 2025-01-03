package coze

import (
	"bufio"
	"context"
	"net/http"
)

type workflowsChat struct {
	client *core
}

func (r *workflowsChat) Stream(ctx context.Context, req *WorkflowsChatStreamReq) (Stream[ChatEvent], error) {
	method := http.MethodPost
	uri := "/v1/workflows/chat"
	resp, err := r.client.RawRequest(ctx, method, uri, req)
	if err != nil {
		return nil, err
	}

	return &streamReader[ChatEvent]{
		ctx:          ctx,
		response:     resp,
		reader:       bufio.NewReader(resp.Body),
		processor:    parseChatEvent,
		httpResponse: newHTTPResponse(resp),
	}, nil
}

func newWorkflowsChat(core *core) *workflowsChat {
	return &workflowsChat{
		client: core,
	}
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
