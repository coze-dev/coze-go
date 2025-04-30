package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConversations(t *testing.T) {
	// Test List method
	t.Run("List conversations success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/conversations", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_bot_id", req.URL.Query().Get("bot_id"))
				assert.Equal(t, "1", req.URL.Query().Get("page_num"))
				assert.Equal(t, "20", req.URL.Query().Get("page_size"))

				// Return mock response
				return mockResponse(http.StatusOK, &listConversationsResp{
					Data: &ListConversationsResp{
						HasMore: true,
						Conversations: []*Conversation{
							{
								ID:            "conv1",
								CreatedAt:     1234567890,
								LastSectionID: "section1",
								MetaData: map[string]string{
									"key1": "value1",
								},
							},
							{
								ID:            "conv2",
								CreatedAt:     1234567891,
								LastSectionID: "section2",
								MetaData: map[string]string{
									"key2": "value2",
								},
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		conversations := newConversations(core)

		paged, err := conversations.List(context.Background(), &ListConversationsReq{
			BotID:    "test_bot_id",
			PageNum:  1,
			PageSize: 20,
		})

		require.NoError(t, err)
		assert.True(t, paged.HasMore())
		items := paged.Items()
		require.Len(t, items, 2)

		// Verify first conversation
		assert.Equal(t, "conv1", items[0].ID)
		assert.Equal(t, 1234567890, items[0].CreatedAt)
		assert.Equal(t, "section1", items[0].LastSectionID)
		assert.Equal(t, "value1", items[0].MetaData["key1"])

		// Verify second conversation
		assert.Equal(t, "conv2", items[1].ID)
		assert.Equal(t, 1234567891, items[1].CreatedAt)
		assert.Equal(t, "section2", items[1].LastSectionID)
		assert.Equal(t, "value2", items[1].MetaData["key2"])
	})

	// Test Create method
	t.Run("Create conversation success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/conversation/create", req.URL.Path)

				// Return mock response
				return mockResponse(http.StatusOK, &createConversationsResp{
					Conversation: &CreateConversationsResp{
						Conversation: Conversation{
							ID:            "conv1",
							CreatedAt:     1234567890,
							LastSectionID: "section1",
							MetaData: map[string]string{
								"key1": "value1",
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		conversations := newConversations(core)

		resp, err := conversations.Create(context.Background(), &CreateConversationsReq{
			Messages: []*Message{
				{
					Role:    "user",
					Content: "Hello",
				},
			},
			MetaData: map[string]string{
				"key1": "value1",
			},
			BotID: "test_bot_id",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "conv1", resp.ID)
		assert.Equal(t, 1234567890, resp.CreatedAt)
		assert.Equal(t, "section1", resp.LastSectionID)
		assert.Equal(t, "value1", resp.MetaData["key1"])
	})

	// Test Retrieve method
	t.Run("Retrieve conversation success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/conversation/retrieve", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "conv1", req.URL.Query().Get("conversation_id"))

				// Return mock response
				return mockResponse(http.StatusOK, &retrieveConversationsResp{
					Conversation: &RetrieveConversationsResp{
						Conversation: Conversation{
							ID:            "conv1",
							CreatedAt:     1234567890,
							LastSectionID: "section1",
							MetaData: map[string]string{
								"key1": "value1",
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		conversations := newConversations(core)

		resp, err := conversations.Retrieve(context.Background(), &RetrieveConversationsReq{
			ConversationID: "conv1",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "conv1", resp.ID)
		assert.Equal(t, 1234567890, resp.CreatedAt)
		assert.Equal(t, "section1", resp.LastSectionID)
		assert.Equal(t, "value1", resp.MetaData["key1"])
	})

	// Test Clear method
	t.Run("Clear conversation success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/conversations/conv1/clear", req.URL.Path)

				// Return mock response
				return mockResponse(http.StatusOK, &clearConversationsResp{
					Data: &ClearConversationsResp{
						ConversationID: "conv1",
						ID:             "new_section",
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		conversations := newConversations(core)

		resp, err := conversations.Clear(context.Background(), &ClearConversationsReq{
			ConversationID: "conv1",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "conv1", resp.ConversationID)
		assert.Equal(t, "new_section", resp.ID)
	})

	// Test List method with default pagination
	t.Run("List conversations with default pagination", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify default pagination parameters
				assert.Equal(t, "1", req.URL.Query().Get("page_num"))
				assert.Equal(t, "20", req.URL.Query().Get("page_size"))

				// Return mock response
				return mockResponse(http.StatusOK, &listConversationsResp{
					Data: &ListConversationsResp{
						HasMore:       false,
						Conversations: []*Conversation{},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		conversations := newConversations(core)

		paged, err := conversations.List(context.Background(), &ListConversationsReq{
			BotID: "test_bot_id",
		})

		require.NoError(t, err)
		assert.False(t, paged.HasMore())
		assert.Empty(t, paged.Items())
	})
}
