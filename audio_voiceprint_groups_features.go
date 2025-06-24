package coze

import (
	"context"
	"net/http"
)

func (r *audioVoiceprintGroupsFeatures) Create(ctx context.Context, req *CreateVoicePrintGroupFeatureReq) (*CreateVoicePrintGroupFeatureResp, error) {
	response := new(createVoicePrintGroupFeatureResp)
	if err := r.core.rawRequest(ctx, &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/audio/voiceprint_groups/:group_id/features",
		Body:   req,
		IsFile: true,
	}, response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (r *audioVoiceprintGroupsFeatures) Update(ctx context.Context, req *UpdateVoicePrintGroupFeatureReq) (*UpdateVoicePrintGroupFeatureResp, error) {
	response := new(updateVoicePrintGroupFeatureResp)
	if err := r.core.rawRequest(ctx, &RawRequestReq{
		Method: http.MethodPut,
		URL:    "/v1/audio/voiceprint_groups/:group_id/features/:feature_id",
		Body:   req,
		IsFile: true,
	}, response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (r *audioVoiceprintGroupsFeatures) Delete(ctx context.Context, req *DeleteVoicePrintGroupFeatureReq) (*DeleteVoicePrintGroupFeatureResp, error) {
	response := new(deleteVoicePrintGroupFeatureResp)
	if err := r.core.rawRequest(ctx, &RawRequestReq{
		Method: http.MethodDelete,
		URL:    "/v1/audio/voiceprint_groups/:group_id/features/:feature_id",
		Body:   req,
	}, response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (r *audioVoiceprintGroupsFeatures) List(ctx context.Context, req *ListVoicePrintGroupFeatureReq) (NumberPaged[VoicePrintGroupFeature], error) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[VoicePrintGroupFeature], error) {
			response := new(listVoicePrintGroupFeatureResp)
			if err := r.core.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/audio/voiceprint_groups/:group_id/features",
				Body:   req.toReq(request),
			}, response); err != nil {
				return nil, err
			}
			return &pageResponse[VoicePrintGroupFeature]{
				Total:   response.Data.Total,
				HasMore: len(response.Data.Items) >= request.PageSize,
				Data:    response.Data.Items,
				LogID:   response.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

type UserInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

type VoicePrintGroupFeature struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	Name      string    `json:"name"`
	AudioURL  string    `json:"audio_url"`
	CreatedAt int       `json:"created_at"`
	UpdatedAt int       `json:"updated_at"`
	Desc      string    `json:"desc"`
	IconURL   string    `json:"icon_url"`
	UserInfo  *UserInfo `json:"user_info"`
}

type FeatureScore struct {
	FeatureID   string  `json:"feature_id"`
	FeatureName string  `json:"feature_name"`
	FeatureDesc string  `json:"feature_desc"`
	Score       float64 `json:"score"`
}

type CreateVoicePrintGroupFeatureReq struct {
	GroupID    string    `path:"group_id" json:"-"`
	Name       string    `json:"name,omitempty"`
	File       FileTypes `json:"file,omitempty"`
	Desc       *string   `json:"desc,omitempty"`
	SampleRate *int      `json:"sample_rate,omitempty"`
	Channel    *int      `json:"channel,omitempty"`
}

type CreateVoicePrintGroupFeatureResp struct {
	baseModel
	ID string `json:"id"`
}

type UpdateVoicePrintGroupFeatureReq struct {
	GroupID    string     `path:"group_id" json:"-"`
	FeatureID  string     `path:"feature_id" json:"-"`
	Name       *string    `json:"name,omitempty"`
	Desc       *string    `json:"desc,omitempty"`
	File       *FileTypes `json:"file,omitempty"`
	SampleRate *int       `json:"sample_rate,omitempty"`
	Channel    *int       `json:"channel,omitempty"`
}

type UpdateVoicePrintGroupFeatureResp struct {
	baseModel
}

type DeleteVoicePrintGroupFeatureReq struct {
	GroupID   string `path:"group_id" json:"-"`
	FeatureID string `path:"feature_id" json:"-"`
}

type DeleteVoicePrintGroupFeatureResp struct {
	baseModel
}

type ListVoicePrintGroupFeatureReq struct {
	GroupID  string `path:"group_id" json:"-"`
	PageSize int    `query:"page_size" json:"-"`
	PageNum  int    `query:"page_num" json:"-"`
}

type createVoicePrintGroupFeatureResp struct {
	baseResponse
	Data *CreateVoicePrintGroupFeatureResp `json:"data"`
}

type updateVoicePrintGroupFeatureResp struct {
	baseResponse
	Data *UpdateVoicePrintGroupFeatureResp `json:"data"`
}

type deleteVoicePrintGroupFeatureResp struct {
	baseResponse
	Data *DeleteVoicePrintGroupFeatureResp `json:"data"`
}

func (r *ListVoicePrintGroupFeatureReq) toReq(request *pageRequest) *ListVoicePrintGroupFeatureReq {
	return &ListVoicePrintGroupFeatureReq{
		GroupID:  r.GroupID,
		PageSize: request.PageSize,
		PageNum:  request.PageNum,
	}
}

type listVoicePrintGroupFeatureRespData struct {
	Total int                       `json:"total"`
	Items []*VoicePrintGroupFeature `json:"items"`
}

type listVoicePrintGroupFeatureResp struct {
	baseResponse
	Data *listVoicePrintGroupFeatureRespData `json:"data"`
}

type audioVoiceprintGroupsFeatures struct {
	core *core
}

func newAudioVoiceprintGroupsFeatures(core *core) *audioVoiceprintGroupsFeatures {
	return &audioVoiceprintGroupsFeatures{
		core: core,
	}
}
