package coze

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestStreamStopsAfterTerminalEvent checks that Chat.Stream and
// Workflows.Runs.Stream stop reading the HTTP response once the terminal event
// has been delivered, both when the server appends a trailing event and when it
// keeps the connection open.
func TestStreamStopsAfterTerminalEvent(t *testing.T) {
	for _, protocol := range []string{"chat", "workflow"} {
		for _, keepOpen := range []bool{false, true} {
			name := protocol + "/trailing_event"
			if keepOpen {
				name = protocol + "/open_connection"
			}
			t.Run(name, func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/event-stream")
					if protocol == "chat" {
						_, _ = io.WriteString(w, "event: done\ndata: [DONE]\n\n")
						if !keepOpen {
							_, _ = io.WriteString(w, "event: conversation.message.delta\ndata: {\"content\":\"after done\"}\n\n")
						}
					} else {
						_, _ = io.WriteString(w, "id: 1\nevent: Done\ndata: {\"debug_url\":\"https://example.com/debug\"}\n\n")
						if !keepOpen {
							_, _ = io.WriteString(w, "id: 2\nevent: Message\ndata: {\"content\":\"after done\"}\n\n")
						}
					}
					w.(http.Flusher).Flush()
					if keepOpen {
						<-r.Context().Done()
					}
				}))
				defer server.Close()
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				client := NewCozeAPI(NewTokenAuth("test-token"), WithBaseURL(server.URL), WithHttpClient(server.Client()))
				if protocol == "chat" {
					stream, err := client.Chat.Stream(ctx, &CreateChatsReq{BotID: "bot", UserID: "user"})
					require.NoError(t, err)
					defer stream.Close()
					event, err := stream.Recv()
					require.NoError(t, err)
					require.Equal(t, ChatEventDone, event.Event)
					assertFinishedStream(t, stream, cancel)
				} else {
					stream, err := client.Workflows.Runs.Stream(ctx, &RunWorkflowsReq{WorkflowID: "workflow"})
					require.NoError(t, err)
					defer stream.Close()
					event, err := stream.Recv()
					require.NoError(t, err)
					require.Equal(t, WorkflowEventTypeDone, event.Event)
					require.Equal(t, "https://example.com/debug", event.DebugURL.URL)
					assertFinishedStream(t, stream, cancel)
				}
			})
		}
	}
}

// assertFinishedStream asserts that repeated Recv calls on a finished stream
// return io.EOF promptly instead of waiting for more HTTP data. It cancels the
// request context before failing so the test server can shut down.
func assertFinishedStream[T streamable](t *testing.T, stream Stream[T], cancel context.CancelFunc) {
	t.Helper()
	type result struct {
		event *T
		err   error
	}
	for i := 0; i < 2; i++ {
		received := make(chan result, 1)
		go func() {
			event, err := stream.Recv()
			received <- result{event, err}
		}()
		select {
		case got := <-received:
			require.Nil(t, got.event)
			require.ErrorIs(t, got.err, io.EOF)
		case <-time.After(time.Second):
			cancel()
			<-received
			t.Fatal("Recv waited for HTTP data after the terminal event")
		}
	}
}
