package coze

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAudioVoices(t *testing.T) {
	// Test Clone method
	t.Run("Clone voice success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/audio/voices/clone", req.URL.Path)

				// Return mock response
				return mockResponse(http.StatusOK, &cloneAudioVoicesResp{
					Data: &CloneAudioVoicesResp{
						VoiceID: "voice1",
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		voices := newVoice(core)

		// Create mock audio file
		audioData := strings.NewReader("mock audio data")
		audioFormat := AudioFormatMP3
		language := LanguageCodeEN
		voiceID := "base_voice"
		previewText := "Hello"
		text := "Sample text"
		description := "Test voice"
		spaceID := "test_space"

		resp, err := voices.Clone(context.Background(), &CloneAudioVoicesReq{
			VoiceName:   "test_voice",
			File:        audioData,
			AudioFormat: audioFormat,
			Language:    &language,
			VoiceID:     &voiceID,
			PreviewText: &previewText,
			Text:        &text,
			Description: &description,
			SpaceID:     &spaceID,
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "voice1", resp.VoiceID)
	})

	// Test Clone method with error
	t.Run("Clone voice with error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		voices := newVoice(core)

		resp, err := voices.Clone(context.Background(), &CloneAudioVoicesReq{
			VoiceName: "test_voice",
			File:      strings.NewReader("invalid audio data"),
		})

		require.Error(t, err)
		assert.Nil(t, resp)
	})

	// Test Clone method with nil file
	t.Run("Clone voice with nil file", func(t *testing.T) {
		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{}})
		voices := newVoice(core)

		resp, err := voices.Clone(context.Background(), &CloneAudioVoicesReq{
			VoiceName: "test_voice",
			File:      nil, // Invalid: nil file
		})

		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "file is required")
	})

	// Test List method
	t.Run("List voices success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/audio/voices", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "1", req.URL.Query().Get("page_num"))
				assert.Equal(t, "20", req.URL.Query().Get("page_size"))
				assert.Equal(t, "true", req.URL.Query().Get("filter_system_voice"))

				// Return mock response
				return mockResponse(http.StatusOK, &ListAudioVoicesResp{
					Data: struct {
						VoiceList []*Voice `json:"voice_list"`
					}{
						VoiceList: []*Voice{
							{
								VoiceID:                "voice1",
								Name:                   "Voice 1",
								IsSystemVoice:          false,
								LanguageCode:           "en-US",
								LanguageName:           "English (US)",
								PreviewText:            "Hello",
								PreviewAudio:           "url1",
								AvailableTrainingTimes: 5,
								CreateTime:             1234567890,
								UpdateTime:             1234567891,
							},
							{
								VoiceID:                "voice2",
								Name:                   "Voice 2",
								IsSystemVoice:          true,
								LanguageCode:           "zh-CN",
								LanguageName:           "Chinese (Simplified)",
								PreviewText:            "你好",
								PreviewAudio:           "url2",
								AvailableTrainingTimes: 3,
								CreateTime:             1234567892,
								UpdateTime:             1234567893,
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		voices := newVoice(core)

		paged, err := voices.List(context.Background(), &ListAudioVoicesReq{
			FilterSystemVoice: true,
			PageNum:           1,
			PageSize:          20,
		})

		require.NoError(t, err)
		assert.False(t, paged.HasMore())
		items := paged.Items()
		require.Len(t, items, 2)

		// Verify first voice
		assert.Equal(t, "voice1", items[0].VoiceID)
		assert.Equal(t, "Voice 1", items[0].Name)
		assert.False(t, items[0].IsSystemVoice)
		assert.Equal(t, "en-US", items[0].LanguageCode)
		assert.Equal(t, "English (US)", items[0].LanguageName)
		assert.Equal(t, "Hello", items[0].PreviewText)
		assert.Equal(t, "url1", items[0].PreviewAudio)
		assert.Equal(t, 5, items[0].AvailableTrainingTimes)
		assert.Equal(t, 1234567890, items[0].CreateTime)
		assert.Equal(t, 1234567891, items[0].UpdateTime)

		// Verify second voice
		assert.Equal(t, "voice2", items[1].VoiceID)
		assert.Equal(t, "Voice 2", items[1].Name)
		assert.True(t, items[1].IsSystemVoice)
		assert.Equal(t, "zh-CN", items[1].LanguageCode)
		assert.Equal(t, "Chinese (Simplified)", items[1].LanguageName)
		assert.Equal(t, "你好", items[1].PreviewText)
		assert.Equal(t, "url2", items[1].PreviewAudio)
		assert.Equal(t, 3, items[1].AvailableTrainingTimes)
		assert.Equal(t, 1234567892, items[1].CreateTime)
		assert.Equal(t, 1234567893, items[1].UpdateTime)
	})

	// Test List method with default pagination
	t.Run("List voices with default pagination", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify default pagination parameters
				assert.Equal(t, "1", req.URL.Query().Get("page_num"))
				assert.Equal(t, "20", req.URL.Query().Get("page_size"))

				// Return mock response with empty list
				return mockResponse(http.StatusOK, &ListAudioVoicesResp{
					Data: struct {
						VoiceList []*Voice `json:"voice_list"`
					}{
						VoiceList: []*Voice{},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		voices := newVoice(core)

		paged, err := voices.List(context.Background(), &ListAudioVoicesReq{})

		require.NoError(t, err)
		assert.False(t, paged.HasMore())
		assert.Empty(t, paged.Items())
	})

	// Test List method with error
	t.Run("List voices with error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		voices := newVoice(core)

		paged, err := voices.List(context.Background(), &ListAudioVoicesReq{
			PageNum:  -1, // Invalid page number
			PageSize: 20,
		})

		require.Error(t, err)
		assert.Nil(t, paged)
	})
}
