package coze

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudioTranscription(t *testing.T) {
	as := assert.New(t)
	t.Run("Transcriptions with different text", func(t *testing.T) {
		text := randomString(10)
		transcription := newTranscriptions(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "/v1/audio/transcriptions", req.URL.Path)
			return mockResponse(http.StatusOK, &CreateAudioTranscriptionsResp{
				Data: AudioTranscriptionsData{
					Text: text,
				},
			})
		})))
		resp, err := transcription.Create(context.Background(), &AudioSpeechTranscriptionsReq{
			Filename: "test.mp3",
			Audio:    strings.NewReader(randomString(10)),
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal(text, resp.Data.Text)
	})

	t.Run("Transcription error", func(t *testing.T) {
		transcription := newTranscriptions(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "/v1/audio/transcriptions", req.URL.Path)
			return nil, fmt.Errorf("test error")
		})))
		_, err := transcription.Create(context.Background(), &AudioSpeechTranscriptionsReq{
			Filename: "test.mp3",
			Audio:    strings.NewReader(randomString(10)),
		})
		as.Error(err)
	})
}
