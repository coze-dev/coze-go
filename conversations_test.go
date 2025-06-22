package coze

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCoreWithTransport(transport http.RoundTripper) *core {
	return newCore(&clientOption{
		baseURL:  ComBaseURL,
		client:   &http.Client{Transport: transport},
		logLevel: LogLevelInfo,
		auth:     NewTokenAuth("token"),
	})
}

func randomString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func TestConversations(t *testing.T) {
	as := assert.New(t)

	t.Run("list success", func(t *testing.T) {
		botID := randomString(10)
		conversationIDs := []string{randomString(10), randomString(10)}
		sectionIDs := []string{randomString(10), randomString(10)}
		mockTransport := newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/conversations", req.URL.Path)

			as.Equal(botID, req.URL.Query().Get("bot_id"))
			as.Equal("1", req.URL.Query().Get("page_num"))
			as.Equal("20", req.URL.Query().Get("page_size"))

			return mockResponse(http.StatusOK, &listConversationsResp{
				Data: &ListConversationsResp{
					HasMore: true,
					Conversations: []*Conversation{
						{
							ID:            conversationIDs[0],
							CreatedAt:     1234567890,
							LastSectionID: sectionIDs[0],
							MetaData: map[string]string{
								"key1": "value1",
							},
						},
						{
							ID:            conversationIDs[1],
							CreatedAt:     1234567891,
							LastSectionID: sectionIDs[1],
							MetaData: map[string]string{
								"key2": "value2",
							},
						},
					},
				},
			})
		})

		core := newCoreWithTransport(mockTransport)
		conversations := newConversations(core)

		paged, err := conversations.List(context.Background(), &ListConversationsReq{
			BotID:    botID,
			PageNum:  1,
			PageSize: 20,
		})

		as.Nil(err)
		as.True(paged.HasMore())
		items := paged.Items()
		require.Len(t, items, 2)

		// Verify first conversation
		as.Equal(conversationIDs[0], items[0].ID)
		as.Equal(1234567890, items[0].CreatedAt)
		as.Equal(sectionIDs[0], items[0].LastSectionID)
		as.Equal("value1", items[0].MetaData["key1"])

		// Verify second conversation
		as.Equal(conversationIDs[1], items[1].ID)
		as.Equal(1234567891, items[1].CreatedAt)
		as.Equal(sectionIDs[1], items[1].LastSectionID)
		as.Equal("value2", items[1].MetaData["key2"])
	})

	t.Run("list failed: http error", func(t *testing.T) {
		botID := randomString(10)
		mockTransport := newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/conversations", req.URL.Path)

			return nil, fmt.Errorf("http error")
		})

		core := newCoreWithTransport(mockTransport)
		conversations := newConversations(core)

		_, err := conversations.List(context.Background(), &ListConversationsReq{
			BotID:    botID,
			PageNum:  1,
			PageSize: 20,
		})
		as.NotNil(t, err)
		as.Contains(err.Error(), "http error")
	})

	t.Run("list with default pagination", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				as.Equal("1", req.URL.Query().Get("page_num"))
				as.Equal("20", req.URL.Query().Get("page_size"))

				return mockResponse(http.StatusOK, &listConversationsResp{
					Data: &ListConversationsResp{
						HasMore:       false,
						Conversations: []*Conversation{},
					},
				})
			},
		}

		core := newCoreWithTransport(mockTransport)
		conversations := newConversations(core)

		paged, err := conversations.List(context.Background(), &ListConversationsReq{
			BotID: "test_bot_id",
		})
		as.Nil(err)
		as.False(paged.HasMore())
		as.Empty(paged.Items())
	})

	t.Run("create success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodPost, req.Method)
				as.Equal("/v1/conversation/create", req.URL.Path)

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

		core := newCoreWithTransport(mockTransport)
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

		as.Nil(err)
		as.Equal("test_log_id", resp.LogID())
		as.Equal("conv1", resp.ID)
		as.Equal(1234567890, resp.CreatedAt)
		as.Equal("section1", resp.LastSectionID)
		as.Equal("value1", resp.MetaData["key1"])
	})

	t.Run("retrieve success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/conversation/retrieve", req.URL.Path)

				as.Equal("conv1", req.URL.Query().Get("conversation_id"))

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

		core := newCoreWithTransport(mockTransport)
		conversations := newConversations(core)

		resp, err := conversations.Retrieve(context.Background(), &RetrieveConversationsReq{
			ConversationID: "conv1",
		})

		as.Nil(err)
		as.Equal("test_log_id", resp.LogID())
		as.Equal("conv1", resp.ID)
		as.Equal(1234567890, resp.CreatedAt)
		as.Equal("section1", resp.LastSectionID)
		as.Equal("value1", resp.MetaData["key1"])
	})

	t.Run("clear success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				as.Equal(http.MethodPost, req.Method)
				as.Equal("/v1/conversations/conv1/clear", req.URL.Path)

				// Return mock response
				return mockResponse(http.StatusOK, &clearConversationsResp{
					Data: &ClearConversationsResp{
						ConversationID: "conv1",
						ID:             "new_section",
					},
				})
			},
		}

		core := newCoreWithTransport(mockTransport)
		conversations := newConversations(core)

		resp, err := conversations.Clear(context.Background(), &ClearConversationsReq{
			ConversationID: "conv1",
		})

		as.Nil(err)
		as.Equal("test_log_id", resp.LogID())
		as.Equal("conv1", resp.ConversationID)
		as.Equal("new_section", resp.ID)
	})
}
