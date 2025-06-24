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
	}, response); err != nil {
		return nil, err
	}
	return response.Data, nil
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

type audioVoiceprintGroupsFeatures struct {
	core *core
}

func newAudioVoiceprintGroupsFeatures(core *core) *audioVoiceprintGroupsFeatures {
	return &audioVoiceprintGroupsFeatures{
		core: core,
	}
}
