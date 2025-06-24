package coze

import (
	"context"
	"net/http"
)

func (r *audioVoiceprintGroups) Create(ctx context.Context, req *CreateVoicePrintGroupReq) (*CreateVoicePrintGroupResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/audio/voiceprint_groups",
		Body:   req,
	}
	response := new(createVoicePrintGroupResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

func (r *audioVoiceprintGroups) Update(ctx context.Context, req *UpdateVoicePrintGroupReq) (*UpdateVoicePrintGroupResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPut,
		URL:    "/v1/audio/voiceprint_groups/:group_id",
		Body:   req,
	}
	response := new(updateVoicePrintGroupResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

func (r *audioVoiceprintGroups) Delete(ctx context.Context, req *DeleteVoicePrintGroupReq) (*DeleteVoicePrintGroupResp, error) {
	request := &RawRequestReq{
		Method: http.MethodDelete,
		URL:    "/v1/audio/voiceprint_groups/:group_id",
		Body:   req,
	}
	response := new(deleteVoicePrintGroupResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

func (r *audioVoiceprintGroups) List(ctx context.Context, req *ListVoicePrintGroupReq) (NumberPaged[VoicePrintGroup], error) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[VoicePrintGroup], error) {
			response := new(listVoicePrintGroupResp)
			if err := r.core.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/audio/voiceprint_groups",
				Body:   req.toReq(request),
			}, response); err != nil {
				return nil, err
			}
			return &pageResponse[VoicePrintGroup]{
				Total:   response.Data.Total,
				HasMore: len(response.Data.Items) >= request.PageSize,
				Data:    response.Data.Items,
				LogID:   response.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

type VoicePrintGroup struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Desc         string    `json:"desc"`
	CreatedAt    int       `json:"created_at"`
	UpdatedAt    int       `json:"updated_at"`
	IconURL      string    `json:"icon_url"`
	FeatureCount int       `json:"feature_count"`
	UserInfo     *UserInfo `json:"user_info"`
}

type CreateVoicePrintGroupReq struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type CreateVoicePrintGroupResp struct {
	baseModel
	ID string `json:"id"`
}

type UpdateVoicePrintGroupReq struct {
	GroupID string  `path:"group_id" json:"-"`
	Name    *string `json:"name,omitempty"`
	Desc    *string `json:"desc,omitempty"`
}

type UpdateVoicePrintGroupResp struct {
	baseModel
}

type DeleteVoicePrintGroupReq struct {
	GroupID string `path:"group_id" json:"-"`
}

type DeleteVoicePrintGroupResp struct {
	baseModel
}

type ListVoicePrintGroupReq struct {
	Name     *string `query:"name" json:"-"`
	GroupID  *string `query:"group_id" json:"-"`
	UserID   *string `query:"user_id" json:"-"`
	PageSize int     `query:"page_size" json:"-"`
	PageNum  int     `query:"page_num" json:"-"`
}

type ListVoicePrintGroupResp struct {
	Total int                `json:"total"`
	Items []*VoicePrintGroup `json:"items"`
}

type createVoicePrintGroupResp struct {
	baseResponse
	Data *CreateVoicePrintGroupResp `json:"data"`
}

type updateVoicePrintGroupResp struct {
	baseResponse
	Data *UpdateVoicePrintGroupResp `json:"data"`
}

type deleteVoicePrintGroupResp struct {
	baseResponse
	Data *DeleteVoicePrintGroupResp `json:"data"`
}

type listVoicePrintGroupResp struct {
	baseResponse
	Data *ListVoicePrintGroupResp `json:"data"`
}

func (r *ListVoicePrintGroupReq) toReq(request *pageRequest) *ListVoicePrintGroupReq {
	return &ListVoicePrintGroupReq{
		Name:     r.Name,
		GroupID:  r.GroupID,
		UserID:   r.UserID,
		PageSize: request.PageSize,
		PageNum:  request.PageNum,
	}
}

type audioVoiceprintGroups struct {
	core     *core
	Features *audioVoiceprintGroupsFeatures
}

func newAudioVoiceprintGroups(core *core) *audioVoiceprintGroups {
	return &audioVoiceprintGroups{
		core:     core,
		Features: newAudioVoiceprintGroupsFeatures(core),
	}
}
