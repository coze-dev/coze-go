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
					`{"event":"conversation_message_delta","message":{"content":"Hello"}}`,
					`{"event":"conversation_message_delta","message":{"content":" World"}}`,
					`{"event":"conversation_chat_completed","chat":{"usage":{"token_count":10}}}`,
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(strings.Join(events, "\n"))),
					Header:     make(http.Header),
				}, nil
			},
		}

		core := newCore(&http.Client{Transport: mockTransport}, ComBaseURL)
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
		assert.Equal(t, ChatEventConversationMessageDelta, event1.Event)
		assert.Equal(t, "Hello", event1.Message.Content)

		// Verify second event
		event2, err := stream.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventConversationMessageDelta, event2.Event)
		assert.Equal(t, " World", event2.Message.Content)

		// Verify completion event
		event3, err := stream.Recv()
		require.NoError(t, err)
		assert.Equal(t, ChatEventConversationChatCompleted, event3.Event)
		assert.Equal(t, int64(10), event3.Chat.Usage.TokenCount)

		// Verify stream end
		_, err = stream.Recv()
		assert.Equal(t, io.EOF, err)
	})

	t.Run("Stream chat with error response", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body: io.NopCloser(strings.NewReader(`{
						"code": "invalid_request",
						"message": "Invalid workflow ID",
						"log_id": "test_log_id"
					}`)),
					Header: make(http.Header),
				}, nil
			},
		}

		core := newCore(&http.Client{Transport: mockTransport}, ComBaseURL)
		chat := newWorkflowsChat(core)

		req := &WorkflowsChatStreamReq{
			WorkflowID: "invalid_workflow",
		}

		_, err := chat.Stream(context.Background(), req)
		require.Error(t, err)

		// Verify error details
		cozeErr, ok := AsCozeError(err)
		require.True(t, ok)
		assert.Equal(t, "invalid_request", cozeErr.Code)
		assert.Equal(t, "Invalid workflow ID", cozeErr.Message)
		assert.Equal(t, "test_log_id", cozeErr.LogID)
	})

	t.Run("Stream chat with network error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, io.ErrUnexpectedEOF
			},
		}

		core := newCore(&http.Client{Transport: mockTransport}, ComBaseURL)
		chat := newWorkflowsChat(core)

		req := &WorkflowsChatStreamReq{
			WorkflowID: "test_workflow",
		}

		_, err := chat.Stream(context.Background(), req)
		require.Error(t, err)
		assert.Equal(t, io.ErrUnexpectedEOF, err)
	})
}
