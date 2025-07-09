package coze

import (
	"context"
	"net/http"
)

// Create adds feedback to a conversation message
func (r *conversationsMessagesFeedback) Create(ctx context.Context, req *CreateConversationMessageFeedbackReq, options ...CozeAPIOption) (*CreateConversationMessageFeedbackResp, error) {
	request := &RawRequestReq{
		Method:  http.MethodPost,
		URL:     "/v1/conversations/:conversation_id/messages/:message_id/feedback",
		Body:    req,
		options: options,
	}
	response := new(createConversationMessageFeedbackResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// Delete removes feedback from a conversation message
func (r *conversationsMessagesFeedback) Delete(ctx context.Context, req *DeleteConversationMessageFeedbackReq, options ...CozeAPIOption) (*DeleteConversationMessageFeedbackResp, error) {
	request := &RawRequestReq{
		Method:  http.MethodDelete,
		URL:     "/v1/conversations/:conversation_id/messages/:message_id/feedback",
		Body:    req,
		options: options,
	}
	response := new(deleteConversationMessageFeedbackResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// FeedbackType represents the type of feedback
type FeedbackType string

const (
	// FeedbackTypeLike represents a like feedback
	FeedbackTypeLike FeedbackType = "like"
	// FeedbackTypeUnlike represents an unlike feedback
	FeedbackTypeUnlike FeedbackType = "unlike"
)

// CreateConversationMessageFeedbackReq represents request for creating message feedback
type CreateConversationMessageFeedbackReq struct {
	// The ID of the conversation.
	ConversationID string `path:"conversation_id" json:"-"`
	// The ID of the message.
	MessageID string `path:"message_id" json:"-"`
	// The type of feedback.
	FeedbackType FeedbackType `json:"feedback_type"`
	// Optional reasons for the feedback.
	ReasonTypes []string `json:"reason_types,omitempty"`
	// Optional comment for the feedback.
	Comment *string `json:"comment,omitempty"`
}

// CreateConversationMessageFeedbackResp represents response for creating message feedback
type CreateConversationMessageFeedbackResp struct {
	baseModel
}

// DeleteConversationMessageFeedbackReq represents request for deleting message feedback
type DeleteConversationMessageFeedbackReq struct {
	// The ID of the conversation.
	ConversationID string `path:"conversation_id" json:"-"`
	// The ID of the message.
	MessageID string `path:"message_id" json:"-"`
}

// DeleteConversationMessageFeedbackResp represents response for deleting message feedback
type DeleteConversationMessageFeedbackResp struct {
	baseModel
}

// createConversationMessageFeedbackResp represents response for creating message feedback
type createConversationMessageFeedbackResp struct {
	baseResponse
	Data *CreateConversationMessageFeedbackResp `json:"data"`
}

// deleteConversationMessageFeedbackResp represents response for deleting message feedback
type deleteConversationMessageFeedbackResp struct {
	baseResponse
	Data *DeleteConversationMessageFeedbackResp `json:"data"`
}

// conversationsMessagesFeedback handles feedback operations for conversation messages
type conversationsMessagesFeedback struct {
	core *core
}

func newConversationsMessagesFeedback(core *core) *conversationsMessagesFeedback {
	return &conversationsMessagesFeedback{core: core}
}
