package coze

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudioRooms(t *testing.T) {
	as := assert.New(t)
	t.Run("create", func(t *testing.T) {
		t.Run("with all fields", func(t *testing.T) {
			createRoomResp := &CreateAudioRoomsResp{
				RoomID: randomString(10),
				AppID:  randomString(10),
				Token:  randomString(10),
				UID:    randomString(10),
			}
			rooms := newRooms(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/audio/rooms", req.URL.Path)
				return mockResponse(http.StatusOK, &createAudioRoomsResp{
					Data: createRoomResp,
				})
			})))
			resp, err := rooms.Create(context.Background(), &CreateAudioRoomsReq{
				BotID:          randomString(10),
				ConversationID: randomString(10),
				VoiceID:        randomString(10),
				WorkflowID:     randomString(10),
				Config: &RoomConfig{
					AudioConfig: &RoomAudioConfig{
						Codec: AudioCodecOPUS,
					},
					VideoConfig: &RoomVideoConfig{
						Codec:           VideoCodecH264,
						StreamVideoType: StreamVideoTypeMain,
					},
					PrologueContent: randomString(10),
				},
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.Equal(createRoomResp.RoomID, resp.RoomID)
			as.Equal(createRoomResp.AppID, resp.AppID)
			as.Equal(createRoomResp.Token, resp.Token)
			as.Equal(createRoomResp.UID, resp.UID)
		})

		t.Run("minimal fields", func(t *testing.T) {
			createRoomResp := &CreateAudioRoomsResp{
				RoomID: randomString(10),
				AppID:  randomString(10),
				Token:  randomString(10),
				UID:    randomString(10),
			}
			rooms := newRooms(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/audio/rooms", req.URL.Path)
				return mockResponse(http.StatusOK, &createAudioRoomsResp{
					Data: createRoomResp,
				})
			})))
			resp, err := rooms.Create(context.Background(), &CreateAudioRoomsReq{
				BotID: randomString(10),
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.Equal(createRoomResp.RoomID, resp.RoomID)
			as.Equal(createRoomResp.AppID, resp.AppID)
			as.Equal(createRoomResp.Token, resp.Token)
			as.Equal(createRoomResp.UID, resp.UID)
		})

		t.Run("failed", func(t *testing.T) {
			rooms := newRooms(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("test error")
			})))

			_, err := rooms.Create(context.Background(), &CreateAudioRoomsReq{
				BotID: randomString(10),
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})
}

func TestAudioConst(t *testing.T) {
	as := assert.New(t)
	t.Run("AudioCodec", func(t *testing.T) {
		assert.Equal(t, AudioCodec("AACLC"), AudioCodecAACLC)
		as.Equal(AudioCodec("G711A"), AudioCodecG711A)
		as.Equal(AudioCodec("OPUS"), AudioCodecOPUS)
		as.Equal(AudioCodec("G722"), AudioCodecG722)
	})

	t.Run("VideoCodec", func(t *testing.T) {
		as.Equal(VideoCodec("H264"), VideoCodecH264)
		as.Equal(VideoCodec("BYTEVC1"), VideoCodecBYTEVC1)
	})
	t.Run("StreamVideoType", func(t *testing.T) {
		as.Equal(StreamVideoType("main"), StreamVideoTypeMain)
		as.Equal(StreamVideoType("screen"), StreamVideoTypeScreen)
	})
}
