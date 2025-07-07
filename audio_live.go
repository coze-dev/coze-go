package coze

import (
	"context"
	"net/http"
)

// Retrieve retrieves live stream information
func (r *audioLive) Retrieve(ctx context.Context, req *RetrieveAudioLiveReq) (*LiveInfo, error) {
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    "/v1/audio/live/:live_id",
		Body:   req,
	}
	response := new(retrieveAudioLiveResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// LiveType represents the type of live stream
type LiveType string

const (
	LiveTypeOrigin      LiveType = "origin"
	LiveTypeTranslation LiveType = "translation"
)

func (l LiveType) String() string {
	return string(l)
}

func (l LiveType) Ptr() *LiveType {
	return &l
}

// StreamInfo represents information about a stream
type StreamInfo struct {
	StreamID string   `json:"stream_id"`
	Name     string   `json:"name"`
	LiveType LiveType `json:"live_type"`
}

// LiveInfo represents information about a live session
type LiveInfo struct {
	baseModel
	AppID       string        `json:"app_id"`
	StreamInfos []*StreamInfo `json:"stream_infos"`
}

// RetrieveAudioLiveReq represents the request for retrieving live information
type RetrieveAudioLiveReq struct {
	LiveID string `path:"live_id" json:"-"`
}

type retrieveAudioLiveResp struct {
	baseResponse
	Data *LiveInfo `json:"data"`
}

// audioLive provides operations for live audio streams
type audioLive struct {
	core *core
}

func newAudioLive(core *core) *audioLive {
	return &audioLive{core: core}
}
