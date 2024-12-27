package coze

import (
	"context"
	"net/http"
)

func (r *conversationsMessages) Create(ctx context.Context, req *CreateMessageReq) (*CreateMessageResp, error) {
	method := http.MethodPost
	uri := "/v1/conversation/message/create"
	resp := &createMessageResp{}

	err := r.core.Request(ctx, method, uri, req, resp,
		withHTTPQuery("conversation_id", req.ConversationID))
	if err != nil {
		return nil, err
	}
	resp.Message.setHTTPResponse(resp.HTTPResponse)
	return resp.Message, nil
}

func (r *conversationsMessages) List(ctx context.Context, req *ListConversationsMessagesReq) (LastIDPaged[Message], error) {
	if req.Limit == 0 {
		req.Limit = 20
	}
	return NewLastIDPaged[Message](
		func(request *pageRequest) (*pageResponse[Message], error) {
			uri := "/v1/conversation/message/list"
			resp := &listConversationsMessagesResp{}
			doReq := &ListConversationsMessagesReq{
				Order:    req.Order,
				ChatID:   req.ChatID,
				BotID:    req.BotID,
				BeforeID: req.BeforeID,
				Limit:    request.PageSize,
			}
			if request.PageToken != "" {
				doReq.AfterID = ptr(request.PageToken)
			}
			err := r.core.Request(ctx, http.MethodPost, uri, doReq, resp,
				withHTTPQuery("conversation_id", req.ConversationID))
			if err != nil {
				return nil, err
			}
			return &pageResponse[Message]{
				HasMore: resp.HasMore,
				Data:    resp.Messages,
				LastID:  resp.FirstID,
				NextID:  resp.LastID,
				LogID:   resp.HTTPResponse.LogID(),
			}, nil
		}, req.Limit, req.AfterID)
}

func (r *conversationsMessages) Retrieve(ctx context.Context, req *RetrieveConversationsMessagesReq) (*RetrieveConversationsMessagesResp, error) {
	method := http.MethodGet
	uri := "/v1/conversation/message/retrieve"
	resp := &retrieveConversationsMessagesResp{}
	err := r.core.Request(ctx, method, uri, nil, resp,
		withHTTPQuery("conversation_id", req.ConversationID),
		withHTTPQuery("message_id", req.MessageID),
	)
	if err != nil {
		return nil, err
	}
	resp.Message.setHTTPResponse(resp.HTTPResponse)
	return resp.Message, nil
}

func (r *conversationsMessages) Update(ctx context.Context, req *UpdateConversationMessagesReq) (*UpdateConversationMessagesResp, error) {
	method := http.MethodPost
	uri := "/v1/conversation/message/modify"
	resp := &updateConversationMessagesResp{}
	conversationID := req.ConversationID
	messageID := req.MessageID
	req.ConversationID = ""
	req.MessageID = ""
	err := r.core.Request(ctx, method, uri, req, resp,
		withHTTPQuery("conversation_id", conversationID),
		withHTTPQuery("message_id", messageID),
	)
	if err != nil {
		return nil, err
	}
	resp.Message.setHTTPResponse(resp.HTTPResponse)
	return resp.Message, nil
}

func (r *conversationsMessages) Delete(ctx context.Context, req *DeleteConversationsMessagesReq) (*DeleteConversationsMessagesResp, error) {
	method := http.MethodPost
	uri := "/v1/conversation/message/delete"
	resp := &deleteConversationsMessagesResp{}
	err := r.core.Request(ctx, method, uri, nil, resp,
		withHTTPQuery("conversation_id", req.ConversationID),
		withHTTPQuery("message_id", req.MessageID),
	)
	if err != nil {
		return nil, err
	}
	resp.Message.setHTTPResponse(resp.HTTPResponse)
	return resp.Message, nil
}

type conversationsMessages struct {
	core *core
}

func newConversationMessage(core *core) *conversationsMessages {
	return &conversationsMessages{core: core}
}

// CreateMessageReq represents request for creating message
type CreateMessageReq struct {
	// The ID of the conversation.
	ConversationID string `json:"-"`

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
	ConversationID string `json:"-"`

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

// RetrieveConversationsMessagesReq represents request for retrieving message
type RetrieveConversationsMessagesReq struct {
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
}

// UpdateConversationMessagesReq represents request for updating message
type UpdateConversationMessagesReq struct {
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`

	// The ID of the message.
	MessageID string `json:"message_id"`

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
	ConversationID string `json:"conversation_id"`

	// message id
	MessageID string `json:"message_id"`
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
