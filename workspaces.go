package coze

import (
	"context"
	"net/http"
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
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[Workspace], error) {
			response := new(listWorkspaceResp)
			if err := r.core.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/workspaces",
				Body:   req.toReq(request),
			}, response); err != nil {
				return nil, err
			}
			return &pageResponse[Workspace]{
				Total:   response.Data.TotalCount,
				HasMore: len(response.Data.Workspaces) >= request.PageSize,
				Data:    response.Data.Workspaces,
				LogID:   response.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

// ListWorkspaceReq represents the request parameters for listing workspaces
type ListWorkspaceReq struct {
	PageNum       int     `query:"page_num" json:"-"`
	PageSize      int     `query:"page_size" json:"-"`
	EnterpriseID  *string `query:"enterprise_id" json:"-"`
	UserID        *string `query:"user_id" json:"-"`
	CozeAccountID *string `query:"coze_account_id" json:"-"`
}

func NewListWorkspaceReq() *ListWorkspaceReq {
	return &ListWorkspaceReq{
		PageNum:  1,
		PageSize: 20,
	}
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

// listWorkspaceResp represents the response for listing workspaces
type listWorkspaceResp struct {
	baseResponse
	Data *ListWorkspaceResp
}

func (r ListWorkspaceReq) toReq(request *pageRequest) *ListWorkspaceReq {
	return &ListWorkspaceReq{
		PageNum:       request.PageNum,
		PageSize:      request.PageSize,
		EnterpriseID:  r.EnterpriseID,
		UserID:        r.UserID,
		CozeAccountID: r.CozeAccountID,
	}
}

type workspace struct {
	core    *core
	Members *workspacesMembers
}

func newWorkspace(core *core) *workspace {
	return &workspace{
		core:    core,
		Members: newWorkspacesMembers(core),
	}
}
