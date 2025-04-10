package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatMessages(t *testing.T) {
	t.Run("List messages success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// 验证请求方法和路径
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v3/chat/message/list", req.URL.Path)

				// 验证查询参数
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))
				assert.Equal(t, "test_chat_id", req.URL.Query().Get("chat_id"))

				// 返回模拟响应
				return mockResponse(http.StatusOK, &listChatsMessagesResp{
					ListChatsMessagesResp: &ListChatsMessagesResp{
						Messages: []*Message{
							{
								ID:             "msg1",
								ConversationID: "test_conversation_id",
								Role:           "user",
								Content:        "Hello",
							},
							{
								ID:             "msg2",
								ConversationID: "test_conversation_id",
								Role:           "assistant",
								Content:        "Hi there!",
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newChatMessages(core)

		resp, err := messages.List(context.Background(), &ListChatsMessagesReq{
			ConversationID: "test_conversation_id",
			ChatID:         "test_chat_id",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		require.Len(t, resp.Messages, 2)

		// 验证第一条消息
		assert.Equal(t, "msg1", resp.Messages[0].ID)
		assert.Equal(t, "test_conversation_id", resp.Messages[0].ConversationID)
		assert.Equal(t, "user", resp.Messages[0].Role.String())
		assert.Equal(t, "Hello", resp.Messages[0].Content)

		// 验证第二条消息
		assert.Equal(t, "msg2", resp.Messages[1].ID)
		assert.Equal(t, "test_conversation_id", resp.Messages[1].ConversationID)
		assert.Equal(t, "assistant", resp.Messages[1].Role.String())
		assert.Equal(t, "Hi there!", resp.Messages[1].Content)
	})

	t.Run("List messages with error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// 返回错误响应
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newChatMessages(core)

		resp, err := messages.List(context.Background(), &ListChatsMessagesReq{
			ConversationID: "invalid_conversation_id",
			ChatID:         "invalid_chat_id",
		})

		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("List messages with empty response", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusOK, &listChatsMessagesResp{
					ListChatsMessagesResp: &ListChatsMessagesResp{
						Messages: []*Message{},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newChatMessages(core)

		resp, err := messages.List(context.Background(), &ListChatsMessagesReq{
			ConversationID: "test_conversation_id",
			ChatID:         "test_chat_id",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Empty(t, resp.Messages)
	})

	t.Run("List messages with missing parameters", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// 验证缺失的参数
				assert.Empty(t, req.URL.Query().Get("conversation_id"))
				assert.Empty(t, req.URL.Query().Get("chat_id"))

				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		messages := newChatMessages(core)

		resp, err := messages.List(context.Background(), &ListChatsMessagesReq{})

		require.Error(t, err)
		assert.Nil(t, resp)
	})
}
