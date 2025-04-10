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

func TestWorkflowsChat(t *testing.T) {
	t.Run("Stream chat success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/workflows/chat", req.URL.Path)

				// Return mock response with chat events
				events := []string{
					`event: conversation.chat.created
data: {"id":"chat1","conversation_id":"test_conversation_id","bot_id":"bot1","status":"created"}

event: conversation.message.delta
data: {"id":"msg1","conversation_id":"test_conversation_id","role":"assistant","content":"Hello"}

event: done
data: {}

`,
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(strings.Join(events, "\n"))),
					Header:     make(http.Header),
				}, nil
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		chat := newWorkflowsChat(core)

		// Create test request
		req := &WorkflowsChatStreamReq{
			WorkflowID: "test_workflow",
			AdditionalMessages: []*Message{
				{
					Role:    MessageRoleUser,
					Content: "Hello",
				},
			},
			Parameters: map[string]any{
				"test": "value",
			},
		}

		// Test streaming
		stream, err := chat.Stream(context.Background(), req)
		require.NoError(t, err)
		defer stream.Close()

		// Verify first event
		event1, err := stream.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventConversationChatCreated, event1.Event)

		// Verify second event
		event2, err := stream.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventConversationMessageDelta, event2.Event)

		// Verify completion event
		event3, err := stream.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventDone, event3.Event)

		// Verify stream end
		_, err = stream.Recv()
		assert.Equal(t, io.EOF, err)
	})

	t.Run("Stream chat with error response", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				mockResp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"code": 100,
						"msg": "Invalid workflow ID"
					}`)),
					Header: make(http.Header),
				}
				mockResp.Header.Set("Content-Type", "application/json")
				mockResp.Header.Set(httpLogIDKey, "test_log_id")
				return mockResp, nil
			},
		}

		core := newCore(&clientOption{
			baseURL: ComBaseURL,
			client:  &http.Client{Transport: mockTransport},
		})
		chat := newWorkflowsChat(core)

		req := &WorkflowsChatStreamReq{
			WorkflowID: "invalid_workflow",
		}

		_, err := chat.Stream(context.Background(), req)
		require.Error(t, err)

		// Verify error details
		cozeErr, ok := AsCozeError(err)
		require.True(t, ok)
		assert.Equal(t, 100, cozeErr.Code)
		assert.Equal(t, "Invalid workflow ID", cozeErr.Message)
		assert.Equal(t, "test_log_id", cozeErr.LogID)
	})
}
