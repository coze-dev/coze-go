package coze

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudioVoices(t *testing.T) {
	as := assert.New(t)
	t.Run("Clone voice success", func(t *testing.T) {
		voices := newVoice(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/audio/voices/clone", req.URL.Path)
			return mockResponse(http.StatusOK, &cloneAudioVoicesResp{
				Data: &CloneAudioVoicesResp{
					VoiceID: "voice1",
				},
			})
		})))
		resp, err := voices.Clone(context.Background(), &CloneAudioVoicesReq{
			VoiceName:   "test_voice",
			File:        strings.NewReader("mock audio data"),
			AudioFormat: AudioFormatMP3,
			Language:    ptr(LanguageCodeEN),
			VoiceID:     ptr("base_voice"),
			PreviewText: ptr("Hello"),
			Text:        ptr("Sample text"),
			Description: ptr("Test voice"),
			SpaceID:     ptr("test_space"),
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal("voice1", resp.VoiceID)
	})

	t.Run("clone voice with error", func(t *testing.T) {
		voices := newVoice(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("test error")
		})))
		_, err := voices.Clone(context.Background(), &CloneAudioVoicesReq{
			VoiceName: "test_voice",
			File:      strings.NewReader("invalid audio data"),
		})
		as.NotNil(err)
	})

	t.Run("list voices success", func(t *testing.T) {
		voices := newVoice(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/audio/voices", req.URL.Path)
			as.Equal("1", req.URL.Query().Get("page_num"))
			as.Equal("20", req.URL.Query().Get("page_size"))
			as.Equal("true", req.URL.Query().Get("filter_system_voice"))
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
		})))
		paged, err := voices.List(context.Background(), &ListAudioVoicesReq{
			FilterSystemVoice: true,
			PageNum:           1,
			PageSize:          20,
		})
		as.Nil(err)
		as.NotNil(paged)
		// as.NotEmpty(paged.logid) // todo
		as.False(paged.HasMore())
		items := paged.Items()
		as.Len(items, 2)

		// Verify first voice
		as.Equal("voice1", items[0].VoiceID)
		as.Equal("Voice 1", items[0].Name)
		as.False(items[0].IsSystemVoice)
		as.Equal("en-US", items[0].LanguageCode)
		as.Equal("English (US)", items[0].LanguageName)
		as.Equal("Hello", items[0].PreviewText)
		as.Equal("url1", items[0].PreviewAudio)
		as.Equal(5, items[0].AvailableTrainingTimes)
		as.Equal(1234567890, items[0].CreateTime)
		as.Equal(1234567891, items[0].UpdateTime)

		// Verify second voice
		as.Equal("voice2", items[1].VoiceID)
		as.Equal("Voice 2", items[1].Name)
		as.True(items[1].IsSystemVoice)
		as.Equal("zh-CN", items[1].LanguageCode)
		as.Equal("Chinese (Simplified)", items[1].LanguageName)
		as.Equal("你好", items[1].PreviewText)
		as.Equal("url2", items[1].PreviewAudio)
		as.Equal(3, items[1].AvailableTrainingTimes)
		as.Equal(1234567892, items[1].CreateTime)
		as.Equal(1234567893, items[1].UpdateTime)
	})

	t.Run("list voices with default pagination", func(t *testing.T) {
		voices := newVoice(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal("1", req.URL.Query().Get("page_num"))
			as.Equal("20", req.URL.Query().Get("page_size"))
			return mockResponse(http.StatusOK, &ListAudioVoicesResp{
				Data: struct {
					VoiceList []*Voice `json:"voice_list"`
				}{
					VoiceList: []*Voice{},
				},
			})
		})))
		paged, err := voices.List(context.Background(), &ListAudioVoicesReq{})
		as.Nil(err)
		as.NotNil(paged)
		// as.NotEmpty(paged.Response().LogID()) // todo
		as.False(paged.HasMore())
		as.Empty(paged.Items())
	})

	t.Run("list voices with error", func(t *testing.T) {
		voices := newVoice(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("test error")
		})))
		_, err := voices.List(context.Background(), &ListAudioVoicesReq{
			PageNum:  -1, // Invalid page number
			PageSize: 20,
		})
		as.NotNil(err)
	})
}
