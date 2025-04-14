package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConversationsMessages(t *testing.T) {
	// Test Create method
	t.Run("Create message success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/conversation/message/create", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))

				// Return mock response
				return mockResponse(http.StatusOK, &createMessageResp{
					Message: &CreateMessageResp{
						Message: Message{
							ID:             "msg1",
							ConversationID: "test_conversation_id",
							Role:           "user",
							Content:        "Hello",
							ContentType:    MessageContentTypeText,
							MetaData: map[string]string{
								"key1": "value1",
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newConversationMessage(core)

		resp, err := messages.Create(context.Background(), &CreateMessageReq{
			ConversationID: "test_conversation_id",
			Role:           "user",
			Content:        "Hello",
			ContentType:    MessageContentTypeText,
			MetaData: map[string]string{
				"key1": "value1",
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "msg1", resp.ID)
		assert.Equal(t, "test_conversation_id", resp.ConversationID)
		assert.Equal(t, "user", string(resp.Role))
		assert.Equal(t, "Hello", resp.Content)
		assert.Equal(t, MessageContentTypeText, resp.ContentType)
		assert.Equal(t, "value1", resp.MetaData["key1"])
	})

	// Test List method
	t.Run("List messages success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/conversation/message/list", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))

				// Return mock response
				return mockResponse(http.StatusOK, &listConversationsMessagesResp{
					ListConversationsMessagesResp: &ListConversationsMessagesResp{
						HasMore: true,
						FirstID: "msg1",
						LastID:  "msg2",
						Messages: []*Message{
							{
								ID:             "msg1",
								ConversationID: "test_conversation_id",
								Role:           "user",
								Content:        "Hello",
								ContentType:    MessageContentTypeText,
							},
							{
								ID:             "msg2",
								ConversationID: "test_conversation_id",
								Role:           "assistant",
								Content:        "Hi there!",
								ContentType:    MessageContentTypeText,
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newConversationMessage(core)

		paged, err := messages.List(context.Background(), &ListConversationsMessagesReq{
			ConversationID: "test_conversation_id",
			Limit:          20,
		})

		require.NoError(t, err)
		assert.True(t, paged.HasMore())
		items := paged.Items()
		require.Len(t, items, 2)

		// Verify first message
		assert.Equal(t, "msg1", items[0].ID)
		assert.Equal(t, "test_conversation_id", items[0].ConversationID)
		assert.Equal(t, "user", string(items[0].Role))
		assert.Equal(t, "Hello", items[0].Content)

		// Verify second message
		assert.Equal(t, "msg2", items[1].ID)
		assert.Equal(t, "test_conversation_id", items[1].ConversationID)
		assert.Equal(t, "assistant", string(items[1].Role))
		assert.Equal(t, "Hi there!", items[1].Content)
	})

	// Test Retrieve method
	t.Run("Retrieve message success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/conversation/message/retrieve", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))
				assert.Equal(t, "msg1", req.URL.Query().Get("message_id"))

				// Return mock response
				return mockResponse(http.StatusOK, &retrieveConversationsMessagesResp{
					Message: &RetrieveConversationsMessagesResp{
						Message: Message{
							ID:             "msg1",
							ConversationID: "test_conversation_id",
							Role:           "user",
							Content:        "Hello",
							ContentType:    MessageContentTypeText,
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newConversationMessage(core)

		resp, err := messages.Retrieve(context.Background(), &RetrieveConversationsMessagesReq{
			ConversationID: "test_conversation_id",
			MessageID:      "msg1",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "msg1", resp.ID)
		assert.Equal(t, "test_conversation_id", resp.ConversationID)
		assert.Equal(t, "user", string(resp.Role))
		assert.Equal(t, "Hello", resp.Content)
		assert.Equal(t, MessageContentTypeText, resp.ContentType)
	})

	// Test Update method
	t.Run("Update message success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/conversation/message/modify", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))
				assert.Equal(t, "msg1", req.URL.Query().Get("message_id"))

				// Return mock response
				return mockResponse(http.StatusOK, &updateConversationMessagesResp{
					Message: &UpdateConversationMessagesResp{
						Message: Message{
							ID:             "msg1",
							ConversationID: "test_conversation_id",
							Role:           "user",
							Content:        "Updated content",
							ContentType:    MessageContentTypeText,
							MetaData: map[string]string{
								"key2": "value2",
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newConversationMessage(core)

		resp, err := messages.Update(context.Background(), &UpdateConversationMessagesReq{
			ConversationID: "test_conversation_id",
			MessageID:      "msg1",
			Content:        "Updated content",
			ContentType:    MessageContentTypeText,
			MetaData: map[string]string{
				"key2": "value2",
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "msg1", resp.ID)
		assert.Equal(t, "test_conversation_id", resp.ConversationID)
		assert.Equal(t, "user", string(resp.Role))
		assert.Equal(t, "Updated content", resp.Content)
		assert.Equal(t, MessageContentTypeText, resp.ContentType)
		assert.Equal(t, "value2", resp.MetaData["key2"])
	})

	// Test Delete method
	t.Run("Delete message success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/conversation/message/delete", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))
				assert.Equal(t, "msg1", req.URL.Query().Get("message_id"))

				// Return mock response
				return mockResponse(http.StatusOK, &deleteConversationsMessagesResp{
					Message: &DeleteConversationsMessagesResp{
						Message: Message{
							ID:             "msg1",
							ConversationID: "test_conversation_id",
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newConversationMessage(core)

		resp, err := messages.Delete(context.Background(), &DeleteConversationsMessagesReq{
			ConversationID: "test_conversation_id",
			MessageID:      "msg1",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "msg1", resp.ID)
		assert.Equal(t, "test_conversation_id", resp.ConversationID)
	})

	// Test List method with default limit
	t.Run("List messages with default limit", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return mock response
				return mockResponse(http.StatusOK, &listConversationsMessagesResp{
					ListConversationsMessagesResp: &ListConversationsMessagesResp{
						HasMore:  false,
						Messages: []*Message{},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newConversationMessage(core)

		paged, err := messages.List(context.Background(), &ListConversationsMessagesReq{
			ConversationID: "test_conversation_id",
		})

		require.NoError(t, err)
		assert.False(t, paged.HasMore())
		assert.Empty(t, paged.Items())
	})

	// Test Create message with object context
	t.Run("Create message with object context", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return mock response
				return mockResponse(http.StatusOK, &createMessageResp{
					Message: &CreateMessageResp{
						Message: Message{
							ID:             "msg1",
							ConversationID: "test_conversation_id",
							Role:           "user",
							ContentType:    MessageContentTypeObjectString,
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newConversationMessage(core)

		createReq := &CreateMessageReq{
			ConversationID: "test_conversation_id",
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

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, MessageContentTypeObjectString, resp.ContentType)
	})
}
