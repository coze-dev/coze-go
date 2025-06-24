package coze

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudioSpeech(t *testing.T) {
	as := assert.New(t)
	t.Run("create", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			speech := newAudioSpeech(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{},
					Body:       io.NopCloser(strings.NewReader("mock audio data")),
				}
				resp.Header.Set(httpLogIDKey, "test_log_id")
				return resp, nil
			})))
			resp, err := speech.Create(context.Background(), &CreateAudioSpeechReq{
				Input:          randomString(10),
				VoiceID:        randomString(10),
				ResponseFormat: AudioFormatMP3.Ptr(),
				Speed:          ptr[float32](1.0),
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.Nil(resp.WriteToFile("/tmp/test.mp3"))
		})

		t.Run("write failed", func(t *testing.T) {
			speech := newAudioSpeech(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{},
					Body:       io.NopCloser(strings.NewReader("mock audio data")),
				}
				resp.Header.Set(httpLogIDKey, "test_log_id")
				return resp, nil
			})))
			resp, err := speech.Create(context.Background(), &CreateAudioSpeechReq{
				Input:          randomString(10),
				VoiceID:        randomString(10),
				ResponseFormat: AudioFormatMP3.Ptr(),
				Speed:          ptr[float32](1.0),
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.NotNil(resp.WriteToFile(""))
		})

		t.Run("failed", func(t *testing.T) {
			speech := newAudioSpeech(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("test error")
			})))
			_, err := speech.Create(context.Background(), &CreateAudioSpeechReq{
				Input:          "Hello, world!",
				VoiceID:        "invalid_voice",
				ResponseFormat: AudioFormatMP3.Ptr(),
				Speed:          ptr[float32](1.0),
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})
}
