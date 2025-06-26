package coze

import (
	"context"
	"net/http"
)

func (r *conversations) List(ctx context.Context, req *ListConversationsReq) (NumberPaged[Conversation], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[Conversation], error) {
			resp := new(listConversationsResp)
			err := r.client.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/conversations",
				Body: &ListConversationsReq{
					BotID:    req.BotID,
					PageNum:  request.PageNum,
					PageSize: request.PageSize,
				},
			}, resp)
			if err != nil {
				return nil, err
			}
			return &pageResponse[Conversation]{
				response: resp.HTTPResponse,
				HasMore:  resp.Data.HasMore,
				Data:     resp.Data.Conversations,
				LogID:    resp.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

// Create 创建会话
//
// docs: https://www.coze.cn/open/docs/developer_guides/create_conversation
func (r *conversations) Create(ctx context.Context, req *CreateConversationsReq) (*CreateConversationsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/conversation/create",
		Body:   req,
	}
	response := new(createConversationsResp)
	err := r.client.rawRequest(ctx, request, response)
	return response.Conversation, err
}

// Retrieve 查看会话信息
//
// docs: https://www.coze.cn/open/docs/developer_guides/retrieve_conversation
func (r *conversations) Retrieve(ctx context.Context, req *RetrieveConversationsReq) (*RetrieveConversationsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    "/v1/conversation/retrieve",
		Body:   req,
	}
	response := new(retrieveConversationsResp)
	err := r.client.rawRequest(ctx, request, response)
	return response.Conversation, err
}

// Clear 清除上下文
//
// docs: https://www.coze.cn/open/docs/developer_guides/clear_conversation_context
func (r *conversations) Clear(ctx context.Context, req *ClearConversationsReq) (*ClearConversationsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/conversations/:conversation_id/clear",
		Body:   req,
	}
	response := new(clearConversationsResp)
	err := r.client.rawRequest(ctx, request, response)
	return response.Data, err
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

// ListConversationsReq represents request for listing conversations
type ListConversationsReq struct {
	// The ID of the bot.
	BotID string `query:"bot_id" json:"-"`

	// The page number.
	PageNum int `query:"page_num" json:"-"`

	// The page size.
	PageSize int `query:"page_size" json:"-"`
}

// ListConversationsResp represents response for listing conversations
type ListConversationsResp struct {
	baseModel
	HasMore       bool            `json:"has_more"`
	Conversations []*Conversation `json:"conversations"`
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

type CreateConversationsResp struct {
	baseModel
	Conversation
}

// RetrieveConversationsReq represents request for retrieving conversation
type RetrieveConversationsReq struct {
	// The ID of the conversation.
	ConversationID string `query:"conversation_id" json:"-"`
}

type RetrieveConversationsResp struct {
	baseModel
	Conversation
}

// ClearConversationsReq represents request for clearing conversation
type ClearConversationsReq struct {
	// The ID of the conversation.
	ConversationID string `path:"conversation_id" json:"-"`
}

type ClearConversationsResp struct {
	baseModel
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
}

// CreateConversationsResp represents response for creating conversation
type createConversationsResp struct {
	baseResponse
	Conversation *CreateConversationsResp `json:"data"`
}

// listConversationsResp represents response for listing conversations
type listConversationsResp struct {
	baseResponse
	Data *ListConversationsResp `json:"data"`
}

// RetrieveConversationsResp represents response for retrieving conversation
type retrieveConversationsResp struct {
	baseResponse
	Conversation *RetrieveConversationsResp `json:"data"`
}

// ClearConversationsResp represents response for clearing conversation
type clearConversationsResp struct {
	baseResponse
	Data *ClearConversationsResp `json:"data"`
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
