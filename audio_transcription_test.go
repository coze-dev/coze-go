package coze

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAudioTranscription(t *testing.T) {
	// Test Transcription method
	t.Run("Transcription with different text", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/audio/transcriptions", req.URL.Path)
				result := map[string]map[string]string{
					"data": {
						"text": "this_test",
					},
				}
				v, _ := json.Marshal(result)
				// Return mock response with audio data
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{},
					Body:       io.NopCloser(strings.NewReader(string(v))),
				}
				resp.Header.Set(logIDHeader, "test_log_id")
				return resp, nil
			},
		}

		core := newCore(&http.Client{Transport: mockTransport}, ComBaseURL)
		transcription := newTranscription(core)
		reader := strings.NewReader("testmp3")
		resp, err := transcription.Create(context.Background(), &AudioSpeechTranscriptionsReq{
			Filename: "testmp3",
			Audio:    reader,
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.HTTPResponse.LogID())

		// Read and verify response body
		require.NoError(t, err)
		assert.NotEqual(t, resp.Data.Text, "this_test_2")
	})

	t.Run("Transcription error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}
		core := newCore(&http.Client{Transport: mockTransport}, ComBaseURL)
		transcription := newTranscription(core)
		reader := strings.NewReader("testmp3")
		resp, err := transcription.Create(context.Background(), &AudioSpeechTranscriptionsReq{
			Filename: "testmp3",
			Audio:    reader,
		})

		require.Error(t, err)
		assert.Nil(t, resp)
	})
}
