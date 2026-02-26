package coze

import (
	"context"
	"net/http"
)

func (r *chatMessages) List(ctx context.Context, req *ListChatsMessagesReq) (*ListChatsMessagesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    "/v3/chat/message/list",
		Body:   req,
	}
	response := new(listChatsMessagesResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.ListChatsMessagesResp, err
}

// ListChatsMessagesReq represents the request to list messages
type ListChatsMessagesReq struct {
	// The Conversation ID can be viewed in the 'conversation_id' field of the Response when
	// initiating a conversation through the Chat API.
	ConversationID string `query:"conversation_id" json:"-"`

	// The Chat ID can be viewed in the 'id' field of the Response when initiating a chat through the
	// Chat API. If it is a streaming response, check the 'id' field in the chat event of the Response.
	ChatID string `query:"chat_id" json:"-"`
}

type ListChatsMessagesResp struct {
	baseModel
	Messages []*Message `json:"data"`
}

type listChatsMessagesResp struct {
	baseResponse
	*ListChatsMessagesResp
}

type chatMessages struct {
	core *core
}

func newChatMessages(core *core) *chatMessages {
	return &chatMessages{core: core}
}
