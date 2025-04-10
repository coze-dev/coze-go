package coze

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a mock streaming response
func mockStreamResponse(data string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(data)),
		Header:     http.Header{"X-Log-Id": []string{"test_log_id"}},
	}, nil
}

func TestChats(t *testing.T) {
	// Test Create method
	t.Run("Create chat success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v3/chat", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))

				// Return mock response
				return mockResponse(http.StatusOK, &createChatsResp{
					Chat: &CreateChatsResp{Chat: Chat{
						ID:             "chat1",
						ConversationID: "test_conversation_id",
						BotID:          "bot1",
						Status:         ChatStatusCreated,
					}},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		chats := newChats(core)

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

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "chat1", resp.Chat.ID)
		assert.Equal(t, ChatStatusCreated, resp.Chat.Status)
	})

	// Test CreateAndPoll method
	t.Run("CreateAndPoll success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
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
			},
		}

		t.Run("CreateAndPoll success", func(t *testing.T) {
			core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
			chats := newChats(core)

			timeout := 5
			resp, err := chats.CreateAndPoll(context.Background(), &CreateChatsReq{
				ConversationID: "test_conversation_id",
				BotID:          "bot1",
				UserID:         "user1",
			}, &timeout)

			require.NoError(t, err)
			assert.Equal(t, "chat1", resp.Chat.ID)
			assert.Equal(t, ChatStatusCompleted, resp.Chat.Status)
			require.Len(t, resp.Messages, 1)
			assert.Equal(t, "Hello!", resp.Messages[0].Content)
		})
		t.Run("CreateAndPoll success with cancel chat", func(t *testing.T) {
			core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
			chats := newChats(core)

			timeout := 0
			resp, err := chats.CreateAndPoll(context.Background(), &CreateChatsReq{
				ConversationID: "test_conversation_id",
				BotID:          "bot1",
				UserID:         "user1",
			}, &timeout)

			require.NoError(t, err)
			assert.Equal(t, "chat1", resp.Chat.ID)
			assert.Equal(t, ChatStatusCancelled, resp.Chat.Status)
			require.Len(t, resp.Messages, 1)
			assert.Equal(t, "Hello!", resp.Messages[0].Content)
		})
	})

	// Test Stream method
	t.Run("Stream chat success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v3/chat", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))

				// Return mock response with streaming data
				return mockStreamResponse(`event: conversation.chat.created
data: {"id":"chat1","conversation_id":"test_conversation_id","bot_id":"bot1","status":"created"}

event: conversation.message.delta
data: {"id":"msg1","conversation_id":"test_conversation_id","role":"assistant","content":"Hello"}

event: done
data: 
`)
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		chats := newChats(core)

		reader, err := chats.Stream(context.Background(), &CreateChatsReq{
			ConversationID: "test_conversation_id",
			BotID:          "bot1",
			UserID:         "user1",
		})

		require.NoError(t, err)
		defer reader.Close()

		// Read and verify events
		event, err := reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventConversationChatCreated, event.Event)
		assert.Equal(t, "chat1", event.Chat.ID)

		event, err = reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventConversationMessageDelta, event.Event)
		assert.Equal(t, "Hello", event.Message.Content)

		event, err = reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventDone, event.Event)
	})

	// Test Cancel method
	t.Run("Cancel chat success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v3/chat/cancel", req.URL.Path)

				// Return mock response
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
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		chats := newChats(core)

		resp, err := chats.Cancel(context.Background(), &CancelChatsReq{
			ConversationID: "test_conversation_id",
			ChatID:         "chat1",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, ChatStatusCancelled, resp.Chat.Status)
	})

	// Test Retrieve method
	t.Run("Retrieve chat success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v3/chat/retrieve", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))
				assert.Equal(t, "chat1", req.URL.Query().Get("chat_id"))

				// Return mock response
				return mockResponse(http.StatusOK, &retrieveChatsResp{
					Chat: &RetrieveChatsResp{
						Chat: Chat{
							ID:             "chat1",
							ConversationID: "test_conversation_id",
							Status:         ChatStatusCompleted,
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		chats := newChats(core)

		resp, err := chats.Retrieve(context.Background(), &RetrieveChatsReq{
			ConversationID: "test_conversation_id",
			ChatID:         "chat1",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, ChatStatusCompleted, resp.Chat.Status)
	})

	// Test SubmitToolOutputs method
	t.Run("SubmitToolOutputs success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v3/chat/submit_tool_outputs", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))
				assert.Equal(t, "chat1", req.URL.Query().Get("chat_id"))

				// Return mock response
				return mockResponse(http.StatusOK, &submitToolOutputsChatResp{
					Chat: &SubmitToolOutputsChatResp{Chat: Chat{
						ID:             "chat1",
						ConversationID: "test_conversation_id",
						Status:         ChatStatusInProgress,
					}},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		chats := newChats(core)

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

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, ChatStatusInProgress, resp.Chat.Status)
	})

	// Test StreamSubmitToolOutputs method
	t.Run("StreamSubmitToolOutputs success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v3/chat/submit_tool_outputs", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test_conversation_id", req.URL.Query().Get("conversation_id"))
				assert.Equal(t, "chat1", req.URL.Query().Get("chat_id"))

				// Return mock streaming response
				return mockStreamResponse(`event: conversation.chat.in_progress
data: {"id":"chat1","conversation_id":"test_conversation_id","status":"in_progress"}

event: conversation.message.delta
data: {"id":"msg1","conversation_id":"test_conversation_id","role":"assistant","content":"Processing tool output"}

event: done
data: 
`)
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		chats := newChats(core)

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

		require.NoError(t, err)
		defer reader.Close()

		// Read and verify events
		event, err := reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventConversationChatInProgress, event.Event)
		assert.Equal(t, "chat1", event.Chat.ID)

		event, err = reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventConversationMessageDelta, event.Event)
		assert.Equal(t, "Processing tool output", event.Message.Content)

		event, err = reader.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventDone, event.Event)
	})
}
