package coze

import (
	"context"
	"net/http"
)

// List 查看消息列表
//
// docs: https://www.coze.cn/open/docs/developer_guides/list_message
func (r *conversationsMessages) List(ctx context.Context, req *ListConversationsMessagesReq) (LastIDPaged[Message], error) {
	if req.Limit == 0 {
		req.Limit = 20
	}
	return NewLastIDPaged(
		func(request *pageRequest) (*pageResponse[Message], error) {
			response := new(listConversationsMessagesResp)
			if err := r.core.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodPost,
				URL:    "/v1/conversation/message/list",
				Body:   req.toReq(request),
			}, response); err != nil {
				return nil, err
			}
			return &pageResponse[Message]{
				response: response.HTTPResponse,
				HasMore:  response.HasMore,
				Data:     response.Messages,
				LastID:   response.FirstID,
				NextID:   response.LastID,
				LogID:    response.HTTPResponse.LogID(),
			}, nil
		}, req.Limit, req.AfterID)
}

// Create 创建消息
//
// https://www.coze.cn/open/docs/developer_guides/create_message
func (r *conversationsMessages) Create(ctx context.Context, req *CreateMessageReq) (*CreateMessageResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/conversation/message/create",
		Body:   req,
	}
	response := new(createMessageResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Message, err
}

// Retrieve 查看消息详情
//
// docs: https://www.coze.cn/open/docs/developer_guides/retrieve_message
func (r *conversationsMessages) Retrieve(ctx context.Context, req *RetrieveConversationsMessagesReq) (*RetrieveConversationsMessagesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    "/v1/conversation/message/retrieve",
		Body:   req,
	}
	response := new(retrieveConversationsMessagesResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Message, err
}

// Update 修改消息
//
// docs: https://www.coze.cn/open/docs/developer_guides/modify_message
func (r *conversationsMessages) Update(ctx context.Context, req *UpdateConversationMessagesReq) (*UpdateConversationMessagesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/conversation/message/modify",
		Body:   req,
	}
	response := new(updateConversationMessagesResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Message, err
}

// Delete 删除消息
//
// docs: https://www.coze.cn/open/docs/developer_guides/delete_message
func (r *conversationsMessages) Delete(ctx context.Context, req *DeleteConversationsMessagesReq) (*DeleteConversationsMessagesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/conversation/message/delete",
		Body:   req,
	}
	response := new(deleteConversationsMessagesResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Message, err
}

type conversationsMessages struct {
	core     *core
	Feedback *conversationsMessagesFeedback
}

func newConversationMessage(core *core) *conversationsMessages {
	return &conversationsMessages{
		core:     core,
		Feedback: newConversationsMessagesFeedback(core),
	}
}

// CreateMessageReq represents request for creating message
type CreateMessageReq struct {
	// The ID of the conversation.
	ConversationID string `query:"conversation_id" json:"-"`

	// The entity that sent this message.
	Role MessageRole `json:"role"`

	// The content of the message, supporting pure text, multimodal (mixed input of text, images, files),
	// cards, and various types of content.
	Content string `json:"content"`

	// The type of message content.
	ContentType MessageContentType `json:"content_type"`

	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages.
	MetaData map[string]string `json:"meta_data,omitempty"`
}

func (c *CreateMessageReq) SetObjectContext(objs []*MessageObjectString) {
	c.ContentType = MessageContentTypeObjectString
	c.Content = mustToJson(objs)
}

// ListConversationsMessagesReq represents request for listing messages
type ListConversationsMessagesReq struct {
	// The ID of the conversation.
	ConversationID string `query:"conversation_id" json:"-"`

	// The sorting method for the message list.
	Order *string `json:"order,omitempty"`

	// The ID of the Chat.
	ChatID *string `json:"chat_id,omitempty"`

	// Get messages before the specified position.
	BeforeID *string `json:"before_id,omitempty"`

	// Get messages after the specified position.
	AfterID *string `json:"after_id,omitempty"`

	// The amount of data returned per query. Default is 50, with a range of 1 to 50.
	Limit int `json:"limit,omitempty"`

	BotID *string `json:"bot_id,omitempty"`
}

func (r ListConversationsMessagesReq) toReq(page *pageRequest) *ListConversationsMessagesReq {
	return &ListConversationsMessagesReq{
		ConversationID: r.ConversationID,
		Order:          r.Order,
		ChatID:         r.ChatID,
		BotID:          r.BotID,
		BeforeID:       r.BeforeID,
		AfterID:        ptrNotZero(page.PageToken),
		Limit:          page.PageSize,
	}
}

// RetrieveConversationsMessagesReq represents request for retrieving message
type RetrieveConversationsMessagesReq struct {
	ConversationID string `query:"conversation_id" json:"-"`
	MessageID      string `query:"message_id" json:"-"`
}

// UpdateConversationMessagesReq represents request for updating message
type UpdateConversationMessagesReq struct {
	// The ID of the conversation.
	ConversationID string `query:"conversation_id" json:"-"`

	// The ID of the message.
	MessageID string `query:"message_id" json:"-"`

	// The content of the message, supporting pure text, multimodal (mixed input of text, images, files),
	// cards, and various types of content.
	Content string `json:"content,omitempty"`

	MetaData map[string]string `json:"meta_data,omitempty"`

	// The type of message content.
	ContentType MessageContentType `json:"content_type,omitempty"`
}

// DeleteConversationsMessagesReq represents request for deleting message
type DeleteConversationsMessagesReq struct {
	// The ID of the conversation.
	ConversationID string `query:"conversation_id" json:"-"`

	// message id
	MessageID string `query:"message_id" json:"-"`
}

// createMessageResp represents response for creating message
type createMessageResp struct {
	baseResponse
	Message *CreateMessageResp `json:"data"`
}

// CreateMessageResp represents response for creating message
type CreateMessageResp struct {
	baseModel
	Message
}

// ListConversationsMessagesResp represents response for listing messages
type listConversationsMessagesResp struct {
	baseResponse
	*ListConversationsMessagesResp
}

type ListConversationsMessagesResp struct {
	baseModel
	HasMore  bool       `json:"has_more"`
	FirstID  string     `json:"first_id"`
	LastID   string     `json:"last_id"`
	Messages []*Message `json:"data"`
}

// retrieveConversationsMessagesResp represents response for retrieving message
type retrieveConversationsMessagesResp struct {
	baseResponse
	Message *RetrieveConversationsMessagesResp `json:"data"`
}

// RetrieveConversationsMessagesResp represents response for creating message
type RetrieveConversationsMessagesResp struct {
	baseModel
	Message
}

// updateConversationMessagesResp represents response for updating message
type updateConversationMessagesResp struct {
	baseResponse
	Message *UpdateConversationMessagesResp `json:"message"`
}

// UpdateConversationMessagesResp represents response for creating message
type UpdateConversationMessagesResp struct {
	baseModel
	Message
}

// deleteConversationsMessagesResp represents response for deleting message
type deleteConversationsMessagesResp struct {
	baseResponse
	Message *DeleteConversationsMessagesResp `json:"data"`
}

// DeleteConversationsMessagesResp represents response for creating message
type DeleteConversationsMessagesResp struct {
	baseModel
	Message
}
