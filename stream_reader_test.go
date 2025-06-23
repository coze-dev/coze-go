package coze

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockHTTPResponse() *httpResponse {
	header := http.Header{}
	header.Set(httpLogIDKey, "test_log_id")
	return &httpResponse{
		Header: header,
	}
}

// Mock event processor for testing
func mockEventProcessor(ctx context.Context, core *core, line []byte, reader *bufio.Reader) (*WorkflowEvent, bool, error) {
	if len(line) == 0 {
		return nil, false, nil
	}

	// Parse event data
	event := &WorkflowEvent{
		ID:    0,
		Event: WorkflowEventTypeMessage,
		Message: &WorkflowEventMessage{
			Content: string(line),
		},
	}

	// Check if this is the last event
	isDone := string(line) == "done"
	if isDone {
		event.Event = WorkflowEventTypeDone
	}
	return event, isDone, nil
}

func TestStreamReader(t *testing.T) {
	as := assert.New(t)
	ctx := context.Background()
	t.Run("successful event processing", func(t *testing.T) {
		// Create mock response with multiple events
		events := []string{
			"first",
			"second",
			"done",
		}
		resp := createMockResponse(events)

		// Create stream reader
		reader := &streamReader[WorkflowEvent]{
			ctx:          ctx,
			reader:       bufio.NewReader(resp.Body),
			response:     resp,
			processor:    mockEventProcessor,
			httpResponse: mockHTTPResponse(),
		}
		defer reader.Close()

		// Read first event
		event, err := reader.Recv()
		require.NoError(t, err)
		as.Equal(WorkflowEventTypeMessage, event.Event)
		as.Equal("first", event.Message.Content)
		as.False(reader.isFinished)

		// Read second event
		event, err = reader.Recv()
		require.NoError(t, err)
		as.Equal(WorkflowEventTypeMessage, event.Event)
		as.Equal("second", event.Message.Content)
		as.False(reader.isFinished)

		// Read final event
		event, err = reader.Recv()
		require.NoError(t, err)
		as.Equal(WorkflowEventTypeDone, event.Event)
		as.True(reader.isFinished)

		// Try reading after done
		event, err = reader.Recv()
		as.Equal(io.EOF, err)
		as.Nil(event)
	})

	t.Run("empty lines are skipped", func(t *testing.T) {
		events := []string{
			"",
			"test",
			"",
			"done",
		}
		resp := createMockResponse(events)

		reader := &streamReader[WorkflowEvent]{
			ctx:          ctx,
			reader:       bufio.NewReader(resp.Body),
			response:     resp,
			processor:    mockEventProcessor,
			httpResponse: mockHTTPResponse(),
		}
		defer reader.Close()

		// First non-empty event
		event, err := reader.Recv()
		as.Nil(err)
		as.Equal(WorkflowEventTypeMessage, event.Event)
		as.Equal("test", event.Message.Content)

		// Second non-empty event
		event, err = reader.Recv()
		as.Nil(err)
		as.Equal(WorkflowEventTypeDone, event.Event)
	})

	t.Run("error response handling", func(t *testing.T) {
		// Create mock error response
		errorResp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(strings.NewReader(`{
				"log_id": "error_log_id",
				"error": {
					"code": 400,
					"message": "Bad Request"
				}
			}`)),
		}

		reader := &streamReader[WorkflowEvent]{
			ctx:          ctx,
			reader:       bufio.NewReader(errorResp.Body),
			response:     errorResp,
			processor:    mockEventProcessor,
			httpResponse: mockHTTPResponse(),
		}
		defer reader.Close()

		// Attempt to read should return error
		event, err := reader.Recv()
		as.Error(err)
		as.Nil(event)
	})

	t.Run("LogID method", func(t *testing.T) {
		reader := &streamReader[WorkflowEvent]{
			ctx:          ctx,
			httpResponse: mockHTTPResponse(),
		}
		as.Equal("test_log_id", reader.httpResponse.LogID())
	})
}

// Helper function to create mock response with events
func createMockResponse(events []string) *http.Response {
	// Join events with newlines
	body := strings.Join(events, "\n")

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
