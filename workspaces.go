package coze

import (
	"context"
	"net/http"
	"strconv"
)

// 查看空间列表
//
// docs: https://www.coze.cn/open/docs/developer_guides/list_workspace
func (r *workspace) List(ctx context.Context, req *ListWorkspaceReq) (NumberPaged[Workspace], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged[Workspace](
		func(request *pageRequest) (*pageResponse[Workspace], error) {
			uri := "/v1/workspaces"
			resp := &listWorkspaceResp{}
			if err := r.core.Request(ctx, http.MethodGet, uri, nil, resp, req.toReq(request)...); err != nil {
				return nil, err
			}
			return &pageResponse[Workspace]{
				Total:   resp.Data.TotalCount,
				HasMore: len(resp.Data.Workspaces) >= request.PageSize,
				Data:    resp.Data.Workspaces,
				LogID:   resp.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

// ListWorkspaceReq represents the request parameters for listing workspaces
type ListWorkspaceReq struct {
	PageNum       int     `json:"page_num"`
	PageSize      int     `json:"page_size"`
	EnterpriseID  *string `json:"enterprise_id,omitempty"`
	UserID        *string `json:"user_id,omitempty"`
	CozeAccountID *string `json:"coze_account_id,omitempty"`
}

func NewListWorkspaceReq() *ListWorkspaceReq {
	return &ListWorkspaceReq{
		PageNum:  1,
		PageSize: 20,
	}
}

// listWorkspaceResp represents the response for listing workspaces
type listWorkspaceResp struct {
	baseResponse
	Data *ListWorkspaceResp
}

// ListWorkspaceResp represents the response for listing workspaces
type ListWorkspaceResp struct {
	baseModel
	TotalCount int          `json:"total_count"`
	Workspaces []*Workspace `json:"workspaces"`
}

// Workspace represents workspace information
type Workspace struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	IconUrl       string            `json:"icon_url"`
	RoleType      WorkspaceRoleType `json:"role_type"`
	WorkspaceType WorkspaceType     `json:"workspace_type"`
	EnterpriseID  string            `json:"enterprise_id"`
}

// WorkspaceRoleType represents the workspace role type
type WorkspaceRoleType string

const (
	WorkspaceRoleTypeOwner  WorkspaceRoleType = "owner"
	WorkspaceRoleTypeAdmin  WorkspaceRoleType = "admin"
	WorkspaceRoleTypeMember WorkspaceRoleType = "member"
)

// WorkspaceType represents the workspace type
type WorkspaceType string

const (
	WorkspaceTypePersonal WorkspaceType = "personal"
	WorkspaceTypeTeam     WorkspaceType = "team"
)

func (r ListWorkspaceReq) toReq(request *pageRequest) []RequestOption {
	res := []RequestOption{
		withHTTPQuery("page_num", strconv.Itoa(request.PageNum)),
		withHTTPQuery("page_size", strconv.Itoa(request.PageSize)),
	}
	if r.EnterpriseID != nil {
		res = append(res, withHTTPQuery("enterprise_id", *r.EnterpriseID))
	}
	if r.UserID != nil {
		res = append(res, withHTTPQuery("user_id", *r.UserID))
	}
	if r.CozeAccountID != nil {
		res = append(res, withHTTPQuery("coze_account_id", *r.CozeAccountID))
	}
	return res
}

type workspace struct {
	core *core
}

func newWorkspace(core *core) *workspace {
	return &workspace{core: core}
}
