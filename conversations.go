package coze

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

func (r *conversations) List(ctx context.Context, req *ListConversationsReq) (NumberPaged[Conversation], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged[Conversation](
		func(request *pageRequest) (*pageResponse[Conversation], error) {
			uri := "/v1/conversations"
			resp := &listConversationsResp{}
			err := r.client.Request(ctx, http.MethodGet, uri, nil, resp,
				withHTTPQuery("bot_id", req.BotID),
				withHTTPQuery("page_num", strconv.Itoa(request.PageNum)),
				withHTTPQuery("page_size", strconv.Itoa(request.PageSize)))
			if err != nil {
				return nil, err
			}
			return &pageResponse[Conversation]{
				HasMore: resp.Data.HasMore,
				Data:    resp.Data.Conversations,
				LogID:   resp.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

func (r *conversations) Create(ctx context.Context, req *CreateConversationsReq) (*CreateConversationsResp, error) {
	uri := "/v1/conversation/create"
	resp := &createConversationsResp{}
	err := r.client.Request(ctx, http.MethodPost, uri, req, resp)
	if err != nil {
		return nil, err
	}
	resp.Conversation.setHTTPResponse(resp.HTTPResponse)
	return resp.Conversation, nil
}

func (r *conversations) Retrieve(ctx context.Context, req *RetrieveConversationsReq) (*RetrieveConversationsResp, error) {
	uri := "/v1/conversation/retrieve"
	resp := &retrieveConversationsResp{}
	err := r.client.Request(ctx, http.MethodGet, uri, nil, resp, withHTTPQuery("conversation_id", req.ConversationID))
	if err != nil {
		return nil, err
	}
	resp.Conversation.setHTTPResponse(resp.HTTPResponse)
	return resp.Conversation, nil
}

func (r *conversations) Clear(ctx context.Context, req *ClearConversationsReq) (*ClearConversationsResp, error) {
	uri := fmt.Sprintf("/v1/conversations/%s/clear", req.ConversationID)
	resp := &clearConversationsResp{}
	err := r.client.Request(ctx, http.MethodPost, uri, nil, resp)
	if err != nil {
		return nil, err
	}
	resp.Data.setHTTPResponse(resp.HTTPResponse)
	return resp.Data, nil
}

type conversations struct {
	client   *core
	Messages *conversationsMessages
}

func newConversations(core *core) *conversations {
	return &conversations{
		client:   core,
		Messages: newConversationMessage(core),
	}
}

// Conversation represents conversation information
type Conversation struct {
	// The ID of the conversation
	ID string `json:"id"`

	// Indicates the create time of the conversation. The value format is Unix timestamp in seconds.
	CreatedAt int `json:"created_at"`

	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages.
	MetaData map[string]string `json:"meta_data,omitempty"`

	// section_id is used to distinguish the context sections of the session history.
	// The same section is one context.
	LastSectionID string `json:"last_section_id"`
}

// CreateConversationsReq represents request for creating conversation
type CreateConversationsReq struct {
	// Messages in the conversation. For more information, see EnterMessage object.
	Messages []*Message `json:"messages,omitempty"`

	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages.
	MetaData map[string]string `json:"meta_data,omitempty"`

	// Bind and isolate conversation on different bots.
	BotID string `json:"bot_id,omitempty"`

	// Optional: Specify a connector ID. Supports passing in 999 (Chat SDK) and 1024 (API). If not provided, the default is 1024 (API).
	ConnectorID string `json:"connector_id"`
}

// ListConversationsReq represents request for listing conversations
type ListConversationsReq struct {
	// The ID of the bot.
	BotID string `json:"bot_id"`

	// The page number.
	PageNum int `json:"page_num,omitempty"`

	// The page size.
	PageSize int `json:"page_size,omitempty"`
}

// RetrieveConversationsReq represents request for retrieving conversation
type RetrieveConversationsReq struct {
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`
}

// ClearConversationsReq represents request for clearing conversation
type ClearConversationsReq struct {
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`
}

// CreateConversationsResp represents response for creating conversation
type createConversationsResp struct {
	baseResponse
	Conversation *CreateConversationsResp `json:"data"`
}

type CreateConversationsResp struct {
	baseModel
	Conversation
}

// listConversationsResp represents response for listing conversations
type listConversationsResp struct {
	baseResponse
	Data *ListConversationsResp `json:"data"`
}

// ListConversationsResp represents response for listing conversations
type ListConversationsResp struct {
	baseModel
	HasMore       bool            `json:"has_more"`
	Conversations []*Conversation `json:"conversations"`
}

// RetrieveConversationsResp represents response for retrieving conversation
type retrieveConversationsResp struct {
	baseResponse
	Conversation *RetrieveConversationsResp `json:"data"`
}

type RetrieveConversationsResp struct {
	baseModel
	Conversation
}

// ClearConversationsResp represents response for clearing conversation
type clearConversationsResp struct {
	baseResponse
	Data *ClearConversationsResp `json:"data"`
}

type ClearConversationsResp struct {
	baseModel
	ConversationID string `json:"conversation_id"`
}
