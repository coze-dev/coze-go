package coze

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatMessages(t *testing.T) {
	as := assert.New(t)
	t.Run("List messages success", func(t *testing.T) {
		messages := newChatMessages(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v3/chat/message/list", req.URL.Path)
			as.Equal("test_conversation_id", req.URL.Query().Get("conversation_id"))
			as.Equal("test_chat_id", req.URL.Query().Get("chat_id"))
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
		})))
		resp, err := messages.List(context.Background(), &ListChatsMessagesReq{
			ConversationID: "test_conversation_id",
			ChatID:         "test_chat_id",
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal(2, len(resp.Messages))
		as.Equal("msg1", resp.Messages[0].ID)
		as.Equal("test_conversation_id", resp.Messages[0].ConversationID)
		as.Equal("user", resp.Messages[0].Role.String())
		as.Equal("Hello", resp.Messages[0].Content)
		as.Equal("msg2", resp.Messages[1].ID)
		as.Equal("test_conversation_id", resp.Messages[1].ConversationID)
		as.Equal("assistant", resp.Messages[1].Role.String())
		as.Equal("Hi there!", resp.Messages[1].Content)
	})

	t.Run("List messages with error", func(t *testing.T) {
		messages := newChatMessages(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("test error")
		})))
		_, err := messages.List(context.Background(), &ListChatsMessagesReq{
			ConversationID: "invalid_conversation_id",
			ChatID:         "invalid_chat_id",
		})
		as.NotNil(err)
	})

	t.Run("List messages with empty response", func(t *testing.T) {
		messages := newChatMessages(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return mockResponse(http.StatusOK, &listChatsMessagesResp{
				ListChatsMessagesResp: &ListChatsMessagesResp{
					Messages: []*Message{},
				},
			})
		})))
		resp, err := messages.List(context.Background(), &ListChatsMessagesReq{
			ConversationID: "test_conversation_id",
			ChatID:         "test_chat_id",
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Empty(resp.Messages)
	})

	t.Run("List messages with missing parameters", func(t *testing.T) {
		messages := newChatMessages(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Empty(req.URL.Query().Get("conversation_id"))
			as.Empty(req.URL.Query().Get("chat_id"))
			return nil, fmt.Errorf("test error")
		})))
		_, err := messages.List(context.Background(), &ListChatsMessagesReq{})
		as.NotNil(err)
	})
}
