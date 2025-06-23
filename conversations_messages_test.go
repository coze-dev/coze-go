package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversationMessage(t *testing.T) {
	as := assert.New(t)

	t.Run("create success", func(t *testing.T) {
		conversationID := randomString(10)
		messages := newConversationMessage(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/conversation/message/create", req.URL.Path)
			as.Equal(conversationID, req.URL.Query().Get("conversation_id"))

			return mockResponse(http.StatusOK, &createMessageResp{
				Message: &CreateMessageResp{
					Message: Message{
						ID:             "msg1",
						ConversationID: conversationID,
						Role:           "user",
						Content:        "Hello",
						ContentType:    MessageContentTypeText,
						MetaData: map[string]string{
							"key1": "value1",
						},
					},
				},
			})
		})))
		resp, err := messages.Create(context.Background(), &CreateMessageReq{
			ConversationID: conversationID,
			Role:           "user",
			Content:        "Hello",
			ContentType:    MessageContentTypeText,
			MetaData: map[string]string{
				"key1": "value1",
			},
		})
		as.Nil(err)
		as.Equal("test_log_id", resp.LogID())
		as.Equal(conversationID, resp.ConversationID)
		as.Equal("user", string(resp.Role))
		as.Equal("Hello", resp.Content)
		as.Equal(MessageContentTypeText, resp.ContentType)
		as.Equal("value1", resp.MetaData["key1"])
	})

	t.Run("list success", func(t *testing.T) {
		conversationID := randomString(10)
		messages := newConversationMessage(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/conversation/message/list", req.URL.Path)
			as.Equal(conversationID, req.URL.Query().Get("conversation_id"))

			return mockResponse(http.StatusOK, &listConversationsMessagesResp{
				ListConversationsMessagesResp: &ListConversationsMessagesResp{
					HasMore: true,
					FirstID: "msg1",
					LastID:  "msg2",
					Messages: []*Message{
						{
							ID:             "msg1",
							ConversationID: conversationID,
							Role:           "user",
							Content:        "Hello",
							ContentType:    MessageContentTypeText,
						},
						{
							ID:             "msg2",
							ConversationID: conversationID,
							Role:           "assistant",
							Content:        "Hi there!",
							ContentType:    MessageContentTypeText,
						},
					},
				},
			})
		})))
		paged, err := messages.List(context.Background(), &ListConversationsMessagesReq{
			ConversationID: conversationID,
			Limit:          20,
		})
		as.Nil(err)
		as.True(paged.HasMore())
		items := paged.Items()
		as.Len(items, 2)

		as.Equal("msg1", items[0].ID)
		as.Equal(conversationID, items[0].ConversationID)
		as.Equal("user", string(items[0].Role))
		as.Equal("Hello", items[0].Content)

		as.Equal("msg2", items[1].ID)
		as.Equal(conversationID, items[1].ConversationID)
		as.Equal("assistant", string(items[1].Role))
		as.Equal("Hi there!", items[1].Content)
	})

	t.Run("retrieve success", func(t *testing.T) {
		conversationID := randomString(10)
		messages := newConversationMessage(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/conversation/message/retrieve", req.URL.Path)
			as.Equal(conversationID, req.URL.Query().Get("conversation_id"))
			as.Equal("msg1", req.URL.Query().Get("message_id"))

			return mockResponse(http.StatusOK, &retrieveConversationsMessagesResp{
				Message: &RetrieveConversationsMessagesResp{
					Message: Message{
						ID:             "msg1",
						ConversationID: conversationID,
						Role:           "user",
						Content:        "Hello",
						ContentType:    MessageContentTypeText,
					},
				},
			})
		})))
		resp, err := messages.Retrieve(context.Background(), &RetrieveConversationsMessagesReq{
			ConversationID: conversationID,
			MessageID:      "msg1",
		})
		as.Nil(err)
		as.Equal("test_log_id", resp.LogID())
		as.Equal("msg1", resp.ID)
		as.Equal(conversationID, resp.ConversationID)
		as.Equal("user", string(resp.Role))
		as.Equal("Hello", resp.Content)
		as.Equal(MessageContentTypeText, resp.ContentType)
	})

	t.Run("update success", func(t *testing.T) {
		conversationID := randomString(10)
		messages := newConversationMessage(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/conversation/message/modify", req.URL.Path)
			as.Equal(conversationID, req.URL.Query().Get("conversation_id"))
			as.Equal("msg1", req.URL.Query().Get("message_id"))

			return mockResponse(http.StatusOK, &updateConversationMessagesResp{
				Message: &UpdateConversationMessagesResp{
					Message: Message{
						ID:             "msg1",
						ConversationID: conversationID,
						Role:           "user",
						Content:        "Updated content",
						ContentType:    MessageContentTypeText,
						MetaData: map[string]string{
							"key2": "value2",
						},
					},
				},
			})
		})))
		resp, err := messages.Update(context.Background(), &UpdateConversationMessagesReq{
			ConversationID: conversationID,
			MessageID:      "msg1",
			Content:        "Updated content",
			ContentType:    MessageContentTypeText,
			MetaData: map[string]string{
				"key2": "value2",
			},
		})
		as.Nil(err)
		as.Equal("test_log_id", resp.LogID())
		as.Equal("msg1", resp.ID)
		as.Equal(conversationID, resp.ConversationID)
		as.Equal("user", string(resp.Role))
		as.Equal("Updated content", resp.Content)
		as.Equal(MessageContentTypeText, resp.ContentType)
		as.Equal("value2", resp.MetaData["key2"])
	})

	t.Run("delete success", func(t *testing.T) {
		conversationID := randomString(10)
		messages := newConversationMessage(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/conversation/message/delete", req.URL.Path)
			as.Equal(conversationID, req.URL.Query().Get("conversation_id"))
			as.Equal("msg1", req.URL.Query().Get("message_id"))

			// Return mock response
			return mockResponse(http.StatusOK, &deleteConversationsMessagesResp{
				Message: &DeleteConversationsMessagesResp{
					Message: Message{
						ID:             "msg1",
						ConversationID: conversationID,
					},
				},
			})
		})))
		resp, err := messages.Delete(context.Background(), &DeleteConversationsMessagesReq{
			ConversationID: conversationID,
			MessageID:      "msg1",
		})
		as.Nil(err)
		as.Equal("test_log_id", resp.LogID())
		as.Equal("msg1", resp.ID)
		as.Equal(conversationID, resp.ConversationID)
	})

	t.Run("list with default limit", func(t *testing.T) {
		conversationID := randomString(10)
		messages := newConversationMessage(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return mockResponse(http.StatusOK, &listConversationsMessagesResp{
				ListConversationsMessagesResp: &ListConversationsMessagesResp{
					HasMore:  false,
					Messages: []*Message{},
				},
			})
		})))
		paged, err := messages.List(context.Background(), &ListConversationsMessagesReq{
			ConversationID: conversationID,
		})
		as.Nil(err)
		as.False(paged.HasMore())
		as.Empty(paged.Items())
	})

	t.Run("create message with object context", func(t *testing.T) {
		conversationID := randomString(10)
		messages := newConversationMessage(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return mockResponse(http.StatusOK, &createMessageResp{
				Message: &CreateMessageResp{
					Message: Message{
						ID:             "msg1",
						ConversationID: conversationID,
						Role:           "user",
						ContentType:    MessageContentTypeObjectString,
					},
				},
			})
		})))
		createReq := &CreateMessageReq{
			ConversationID: conversationID,
			Role:           "user",
		}
		createReq.SetObjectContext([]*MessageObjectString{
			NewFileMessageObjectByID("file_id"),
			NewAudioMessageObjectByURL("audio_url"),
			NewAudioMessageObjectByID("audio_id"),
			NewFileMessageObjectByURL("file_url"),
			NewImageMessageObjectByID("image_id"),
			NewImageMessageObjectByURL("image_url"),
			NewTextMessageObject("text"),
		})

		resp, err := messages.Create(context.Background(), createReq)
		as.Nil(err)
		as.Equal("test_log_id", resp.LogID())
		as.Equal(MessageContentTypeObjectString, resp.ContentType)
	})
}
