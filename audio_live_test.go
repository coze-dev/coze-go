package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudioLive(t *testing.T) {
	as := assert.New(t)

	t.Run("retrieve", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			liveID := randomString(10)
			expectedAppID := randomString(10)
			expectedStreamInfos := []*StreamInfo{
				{
					StreamID: randomString(10),
					Name:     randomString(10),
					LiveType: LiveTypeOrigin,
				},
				{
					StreamID: randomString(10),
					Name:     randomString(10),
					LiveType: LiveTypeTranslation,
				},
			}

			live := newAudioLive(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/audio/live/"+liveID, req.URL.Path)
				return mockResponse(http.StatusOK, &retrieveAudioLiveResp{
					Data: &LiveInfo{
						AppID:       expectedAppID,
						StreamInfos: expectedStreamInfos,
					},
				})
			})))

			resp, err := live.Retrieve(context.Background(), &RetrieveAudioLiveReq{
				LiveID: liveID,
			})

			as.Nil(err)
			as.NotNil(resp)
			as.Equal(expectedAppID, resp.AppID)
			as.Equal(len(expectedStreamInfos), len(resp.StreamInfos))
			for i, streamInfo := range resp.StreamInfos {
				as.Equal(expectedStreamInfos[i].StreamID, streamInfo.StreamID)
				as.Equal(expectedStreamInfos[i].Name, streamInfo.Name)
				as.Equal(expectedStreamInfos[i].LiveType, streamInfo.LiveType)
			}
			as.NotEmpty(resp.Response().LogID())
		})

		t.Run("empty_response", func(t *testing.T) {
			liveID := randomString(10)

			live := newAudioLive(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/audio/live/"+liveID, req.URL.Path)
				return mockResponse(http.StatusOK, &retrieveAudioLiveResp{
					Data: &LiveInfo{
						AppID:       randomString(10),
						StreamInfos: []*StreamInfo{},
					},
				})
			})))

			resp, err := live.Retrieve(context.Background(), &RetrieveAudioLiveReq{
				LiveID: liveID,
			})

			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.AppID)
			as.Empty(resp.StreamInfos)
		})

		t.Run("error", func(t *testing.T) {
			liveID := randomString(10)

			live := newAudioLive(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/audio/live/"+liveID, req.URL.Path)
				return nil, errors.New("test error")
			})))

			_, err := live.Retrieve(context.Background(), &RetrieveAudioLiveReq{
				LiveID: liveID,
			})

			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})
}

func TestLiveType(t *testing.T) {
	as := assert.New(t)

	t.Run("constants", func(t *testing.T) {
		as.Equal("origin", LiveTypeOrigin.String())
		as.Equal("translation", LiveTypeTranslation.String())
	})

	t.Run("ptr", func(t *testing.T) {
		liveType := LiveTypeOrigin
		as.Equal(&liveType, liveType.Ptr())
	})
}
