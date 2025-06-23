package coze

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockStreamResponse(data string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(data)),
		Header: http.Header{
			httpLogIDKey:   []string{"test_log_id"},
			"Content-Type": []string{"text/event-stream"},
		},
	}, nil
}

func TestChat(t *testing.T) {
	as := assert.New(t)
	t.Run("create chat success", func(t *testing.T) {
		chats := newChats(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v3/chat", req.URL.Path)
			as.Equal("test_conversation_id", req.URL.Query().Get("conversation_id"))
			return mockResponse(http.StatusOK, &createChatsResp{
				Chat: &CreateChatsResp{Chat: Chat{
					ID:             "chat1",
					ConversationID: "test_conversation_id",
					BotID:          "bot1",
					Status:         ChatStatusCreated,
				}},
			})
		})))
		resp, err := chats.Create(context.Background(), &CreateChatsReq{
			ConversationID: "test_conversation_id",
			BotID:          "bot1",
			UserID:         "user1",
			Messages: []*Message{
				BuildUserQuestionText("hello", nil),
				BuildUserQuestionObjects([]*MessageObjectString{
					NewFileMessageObjectByURL("url"),
				}, nil),
				BuildAssistantAnswer("hello", nil),
			},
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal("chat1", resp.Chat.ID)
		as.Equal(ChatStatusCreated, resp.Chat.Status)
	})

	t.Run("CreateAndPoll success", func(t *testing.T) {
		chats := newChats(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			switch req.URL.Path {
			case "/v3/chat":
				// Return create response
				return mockResponse(http.StatusOK, &createChatsResp{
					Chat: &CreateChatsResp{Chat: Chat{
						ID:             "chat1",
						ConversationID: "test_conversation_id",
						BotID:          "bot1",
						Status:         ChatStatusInProgress,
					}},
				})
			case "/v3/chat/retrieve":
				// Return retrieve response with completed status
				return mockResponse(http.StatusOK, &retrieveChatsResp{
					Chat: &RetrieveChatsResp{
						Chat: Chat{
							ID:             "chat1",
							ConversationID: "test_conversation_id",
							Status:         ChatStatusCompleted,
						},
					},
				})
			case "/v3/chat/message/list":
				// Return message list response
				return mockResponse(http.StatusOK, &listChatsMessagesResp{
					ListChatsMessagesResp: &ListChatsMessagesResp{
						Messages: []*Message{
							{
								ID:             "msg1",
								ConversationID: "test_conversation_id",
								Role:           "assistant",
								Content:        "Hello!",
							},
						},
					},
				})
			case "/v3/chat/cancel":
				return mockResponse(http.StatusOK, &cancelChatsResp{
					Chat: &CancelChatsResp{
						Chat: Chat{
							ID:             "chat1",
							ConversationID: "test_conversation_id",
							BotID:          "bot1",
							Status:         ChatStatusCancelled,
						},
					},
				})
			default:
				t.Fatalf("Unexpected request path: %s", req.URL.Path)
				return nil, nil
			}
		})))
		t.Run("CreateAndPoll success", func(t *testing.T) {
			timeout := 5
			resp, err := chats.CreateAndPoll(context.Background(), &CreateChatsReq{
				ConversationID: "test_conversation_id",
				BotID:          "bot1",
				UserID:         "user1",
			}, &timeout)
			as.Nil(err)
			as.NotNil(resp)
			as.Equal("chat1", resp.Chat.ID)
			as.Equal(ChatStatusCompleted, resp.Chat.Status)
			as.Equal(1, len(resp.Messages))
			as.Equal("Hello!", resp.Messages[0].Content)
		})
		t.Run("CreateAndPoll success with cancel chat", func(t *testing.T) {
			timeout := 0
			resp, err := chats.CreateAndPoll(context.Background(), &CreateChatsReq{
				ConversationID: "test_conversation_id",
				BotID:          "bot1",
				UserID:         "user1",
			}, &timeout)
			as.Nil(err)
			as.NotNil(resp)
			as.Equal("chat1", resp.Chat.ID)
			as.Equal(ChatStatusCancelled, resp.Chat.Status)
			as.Equal(1, len(resp.Messages))
			as.Equal("Hello!", resp.Messages[0].Content)
		})
	})

	t.Run("Stream chat success", func(t *testing.T) {
		chats := newChats(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v3/chat", req.URL.Path)
			as.Equal("test_conversation_id", req.URL.Query().Get("conversation_id"))
			return mockStreamResponse(`event: conversation.chat.created
data: {"id":"chat1","conversation_id":"test_conversation_id1","bot_id":"bot1","status":"created"}

event: conversation.message.delta
data: {"id":"msg1","conversation_id":"test_conversation_id2","role":"assistant","content":"Hello"}

event: done
data: 
`)
		})))
		stream, err := chats.Stream(context.Background(), &CreateChatsReq{
			ConversationID: "test_conversation_id",
			BotID:          "bot1",
			UserID:         "user1",
		})
		as.Nil(err)
		defer stream.Close()

		event, err := stream.Recv()
		as.Nil(err)
		as.Equal(ChatEventConversationChatCreated, event.Event)
		as.Equal("chat1", event.Chat.ID)
		as.Equal("test_conversation_id1", event.Chat.ConversationID)
		as.Equal("bot1", event.Chat.BotID)

		event, err = stream.Recv()
		as.Nil(err)
		as.Equal(ChatEventConversationMessageDelta, event.Event)
		as.Equal("Hello", event.Message.Content)
		as.Equal("test_conversation_id2", event.Message.ConversationID)
		as.Equal(MessageRole("assistant"), event.Message.Role)

		event, err = stream.Recv()
		as.Nil(err)
		as.Equal(ChatEventDone, event.Event)
	})

	t.Run("cancel chat success", func(t *testing.T) {
		chats := newChats(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v3/chat/cancel", req.URL.Path)
			return mockResponse(http.StatusOK, &cancelChatsResp{
				Chat: &CancelChatsResp{
					Chat: Chat{
						ID:             "chat1",
						ConversationID: "test_conversation_id",
						BotID:          "bot1",
						Status:         ChatStatusCancelled,
					},
				},
			})
		})))
		resp, err := chats.Cancel(context.Background(), &CancelChatsReq{
			ConversationID: "test_conversation_id",
			ChatID:         "chat1",
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal(ChatStatusCancelled, resp.Chat.Status)
	})

	// Test Retrieve method
	t.Run("Retrieve chat success", func(t *testing.T) {
		chats := newChats(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v3/chat/retrieve", req.URL.Path)
			as.Equal("test_conversation_id", req.URL.Query().Get("conversation_id"))
			as.Equal("chat1", req.URL.Query().Get("chat_id"))
			return mockResponse(http.StatusOK, &retrieveChatsResp{
				Chat: &RetrieveChatsResp{
					Chat: Chat{
						ID:             "chat1",
						ConversationID: "test_conversation_id",
						Status:         ChatStatusCompleted,
					},
				},
			})
		})))
		resp, err := chats.Retrieve(context.Background(), &RetrieveChatsReq{
			ConversationID: "test_conversation_id",
			ChatID:         "chat1",
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal(ChatStatusCompleted, resp.Chat.Status)
	})

	t.Run("SubmitToolOutputs success", func(t *testing.T) {
		chats := newChats(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v3/chat/submit_tool_outputs", req.URL.Path)
			as.Equal("test_conversation_id", req.URL.Query().Get("conversation_id"))
			as.Equal("chat1", req.URL.Query().Get("chat_id"))
			return mockResponse(http.StatusOK, &submitToolOutputsChatResp{
				Chat: &SubmitToolOutputsChatResp{Chat: Chat{
					ID:             "chat1",
					ConversationID: "test_conversation_id",
					Status:         ChatStatusInProgress,
				}},
			})
		})))
		resp, err := chats.SubmitToolOutputs(context.Background(), &SubmitToolOutputsChatReq{
			ConversationID: "test_conversation_id",
			ChatID:         "chat1",
			ToolOutputs: []*ToolOutput{
				{
					ToolCallID: "tool1",
					Output:     "result1",
				},
			},
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal(ChatStatusInProgress, resp.Chat.Status)
	})

	t.Run("StreamSubmitToolOutputs success", func(t *testing.T) {
		chats := newChats(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v3/chat/submit_tool_outputs", req.URL.Path)
			as.Equal("test_conversation_id", req.URL.Query().Get("conversation_id"))
			as.Equal("chat1", req.URL.Query().Get("chat_id"))
			return mockStreamResponse(`event: conversation.chat.in_progress
data: {"id":"chat1","conversation_id":"test_conversation_id","status":"in_progress"}

event: conversation.message.delta
data: {"id":"msg1","conversation_id":"test_conversation_id","role":"assistant","content":"Processing tool output"}

event: done
data: 
`)
		})))
		reader, err := chats.StreamSubmitToolOutputs(context.Background(), &SubmitToolOutputsChatReq{
			ConversationID: "test_conversation_id",
			ChatID:         "chat1",
			ToolOutputs: []*ToolOutput{
				{
					ToolCallID: "tool1",
					Output:     "result1",
				},
			},
		})
		as.Nil(err)
		defer reader.Close()
		event, err := reader.Recv()
		as.Nil(err)
		as.Equal(ChatEventConversationChatInProgress, event.Event)
		as.Equal("chat1", event.Chat.ID)

		event, err = reader.Recv()
		as.Nil(err)
		as.Equal(ChatEventConversationMessageDelta, event.Event)
		as.Equal("Processing tool output", event.Message.Content)

		event, err = reader.Recv()
		as.Nil(err)
		as.Equal(ChatEventDone, event.Event)
	})
}
