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

func TestAudioSpeech(t *testing.T) {
	// Test Create method
	t.Run("Create speech success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/audio/speech", req.URL.Path)

				// Return mock response with audio data
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{},
					Body:       io.NopCloser(strings.NewReader("mock audio data")),
				}
				resp.Header.Set(httpLogIDKey, "test_log_id")
				return resp, nil
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		speech := newSpeech(core)

		resp, err := speech.Create(context.Background(), &CreateAudioSpeechReq{
			Input:          "Hello, world!",
			VoiceID:        "voice1",
			ResponseFormat: AudioFormatMP3.Ptr(),
			Speed:          ptr[float32](1.0),
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.HTTPResponse.LogID())

		// Read and verify response body
		data, err := io.ReadAll(resp.Data)
		require.NoError(t, err)
		assert.Equal(t, "mock audio data", string(data))
		resp.Data.Close()
	})

	// Test Create method with error
	t.Run("Create speech with error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		speech := newSpeech(core)

		resp, err := speech.Create(context.Background(), &CreateAudioSpeechReq{
			Input:          "Hello, world!",
			VoiceID:        "invalid_voice",
			ResponseFormat: AudioFormatMP3.Ptr(),
			Speed:          ptr[float32](1.0),
		})

		require.Error(t, err)
		assert.Nil(t, resp)
	})

	// Test Create method with invalid speed
	t.Run("Create speech with invalid speed", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response for invalid speed
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		speech := newSpeech(core)

		resp, err := speech.Create(context.Background(), &CreateAudioSpeechReq{
			Input:          "Hello, world!",
			VoiceID:        "voice1",
			ResponseFormat: AudioFormatMP3.Ptr(),
			Speed:          ptr[float32](-1.0), // Invalid speed
		})

		require.Error(t, err)
		assert.Nil(t, resp)
	})
}
