package coze

import (
	"context"
	"net/http"
)

func (r *apps) List(ctx context.Context, req *ListAppReq) (NumberPaged[SimpleApp], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[SimpleApp], error) {
			resp := new(listAppResp)
			err := r.core.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/apps",
				Body:   req.toReq(request),
			}, resp)
			if err != nil {
				return nil, err
			}
			return &pageResponse[SimpleApp]{
				response: resp.HTTPResponse,
				HasMore:  len(resp.Data.Items) >= request.PageSize,
				Data:     resp.Data.Items,
				LogID:    resp.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

type ListAppReq struct {
	WorkspaceID   string         `query:"workspace_id" json:"-"`
	PublishStatus *PublishStatus `query:"publish_status" json:"-"`
	ConnectorID   *string        `query:"connector_id" json:"-"`
	PageNum       int            `query:"page_num" json:"-"`
	PageSize      int            `query:"page_size" json:"-"`
}

type SimpleApp struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
	IsPublished bool   `json:"is_published,omitempty"`
	OwnerUserID string `json:"owner_user_id,omitempty"`
	UpdatedAt   int    `json:"updated_at,omitempty"`
	PublishedAt *int   `json:"published_at,omitempty"`
}

type ListAppResp struct {
	Total int          `json:"total"`
	Items []*SimpleApp `json:"items"`
}

type listAppResp struct {
	baseResponse
	Data *ListAppResp `json:"data"`
}

func (r ListAppReq) toReq(request *pageRequest) *ListAppReq {
	return &ListAppReq{
		WorkspaceID:   r.WorkspaceID,
		PublishStatus: r.PublishStatus,
		ConnectorID:   r.ConnectorID,
		PageNum:       request.PageNum,
		PageSize:      request.PageSize,
	}
}

type apps struct {
	core *core
}

func newApps(core *core) *apps {
	return &apps{
		core: core,
	}
}
