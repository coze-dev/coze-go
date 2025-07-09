package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversationMessageFeedback(t *testing.T) {
	as := assert.New(t)

	t.Run("create", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			conversationID := randomString(10)
			messageID := randomString(10)
			feedback := newConversationsMessagesFeedback(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodPost, req.Method)
				as.Equal("/v1/conversations/"+conversationID+"/messages/"+messageID+"/feedback", req.URL.Path)

				return mockResponse(http.StatusOK, &createConversationMessageFeedbackResp{
					Data: &CreateConversationMessageFeedbackResp{},
				})
			})))

			resp, err := feedback.Create(context.Background(), &CreateConversationMessageFeedbackReq{
				ConversationID: conversationID,
				MessageID:      messageID,
				FeedbackType:   FeedbackTypeLike,
				ReasonTypes:    []string{"helpful", "accurate"},
				Comment:        ptr("Great response!"),
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
		})

		t.Run("error", func(t *testing.T) {
			conversationID := randomString(10)
			messageID := randomString(10)
			feedback := newConversationsMessagesFeedback(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodPost, req.Method)
				as.Equal("/v1/conversations/"+conversationID+"/messages/"+messageID+"/feedback", req.URL.Path)
				return nil, errors.New("test error")
			})))

			_, err := feedback.Create(context.Background(), &CreateConversationMessageFeedbackReq{
				ConversationID: conversationID,
				MessageID:      messageID,
				FeedbackType:   FeedbackTypeLike,
				ReasonTypes:    []string{"helpful", "accurate"},
				Comment:        ptr("Great response!"),
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			conversationID := randomString(10)
			messageID := randomString(10)
			feedback := newConversationsMessagesFeedback(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodDelete, req.Method)
				as.Equal("/v1/conversations/"+conversationID+"/messages/"+messageID+"/feedback", req.URL.Path)

				return mockResponse(http.StatusOK, &deleteConversationMessageFeedbackResp{
					Data: &DeleteConversationMessageFeedbackResp{},
				})
			})))

			resp, err := feedback.Delete(context.Background(), &DeleteConversationMessageFeedbackReq{
				ConversationID: conversationID,
				MessageID:      messageID,
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
		})

		t.Run("error", func(t *testing.T) {
			conversationID := randomString(10)
			messageID := randomString(10)
			feedback := newConversationsMessagesFeedback(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodDelete, req.Method)
				as.Equal("/v1/conversations/"+conversationID+"/messages/"+messageID+"/feedback", req.URL.Path)
				return nil, errors.New("test error")
			})))

			_, err := feedback.Delete(context.Background(), &DeleteConversationMessageFeedbackReq{
				ConversationID: conversationID,
				MessageID:      messageID,
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("feedback type constants", func(t *testing.T) {
		as.Equal("like", string(FeedbackTypeLike))
		as.Equal("unlike", string(FeedbackTypeUnlike))
	})
}
