package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAudioRooms(t *testing.T) {
	// Test Create method
	t.Run("Create audio room success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/audio/rooms", req.URL.Path)

				// Return mock response
				return mockResponse(http.StatusOK, &createAudioRoomsResp{
					Data: &CreateAudioRoomsResp{
						RoomID: "room1",
						AppID:  "app1",
						Token:  "token1",
						UID:    "uid1",
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		rooms := newRooms(core)

		// Test with all optional fields
		resp, err := rooms.Create(context.Background(), &CreateAudioRoomsReq{
			BotID:          "bot1",
			ConversationID: "conv1",
			VoiceID:        "voice1",
			Config: &RoomConfig{
				AudioConfig: &RoomAudioConfig{
					Codec: AudioCodecOPUS,
				},
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "room1", resp.RoomID)
		assert.Equal(t, "app1", resp.AppID)
		assert.Equal(t, "token1", resp.Token)
		assert.Equal(t, "uid1", resp.UID)
	})

	// Test Create method with minimal fields
	t.Run("Create audio room with minimal fields", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return mock response
				return mockResponse(http.StatusOK, &createAudioRoomsResp{
					Data: &CreateAudioRoomsResp{
						RoomID: "room1",
						AppID:  "app1",
						Token:  "token1",
						UID:    "uid1",
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		rooms := newRooms(core)

		// Test with only required fields
		resp, err := rooms.Create(context.Background(), &CreateAudioRoomsReq{
			BotID: "bot1",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "room1", resp.RoomID)
	})

	// Test Create method with error
	t.Run("Create audio room with error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		rooms := newRooms(core)

		resp, err := rooms.Create(context.Background(), &CreateAudioRoomsReq{
			BotID: "invalid_bot",
		})

		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestAudioCodec(t *testing.T) {
	t.Run("AudioCodec constants", func(t *testing.T) {
		assert.Equal(t, AudioCodec("AACLC"), AudioCodecAACLC)
		assert.Equal(t, AudioCodec("G711A"), AudioCodecG711A)
		assert.Equal(t, AudioCodec("OPUS"), AudioCodecOPUS)
		assert.Equal(t, AudioCodec("G722"), AudioCodecG722)
	})
}
