package coze

import (
	"context"
	"net/http"
)

// List 查看空间成员列表
//
// docs: https://www.coze.cn/open/docs/developer_guides/list_space_member
func (r *workspacesMembers) List(ctx context.Context, req *ListWorkspaceMemberReq) (NumberPaged[WorkspaceMember], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[WorkspaceMember], error) {
			response := new(listWorkspaceMemberResp)
			if err := r.core.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/workspaces/:workspace_id/members",
				Body:   req.toReq(request),
			}, response); err != nil {
				return nil, err
			}
			return &pageResponse[WorkspaceMember]{
				Total:   response.Data.TotalCount,
				HasMore: len(response.Data.Items) >= request.PageSize,
				Data:    response.Data.Items,
				LogID:   response.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

func (r *workspacesMembers) Create(ctx context.Context, req *CreateWorkspaceMemberReq) (*CreateWorkspaceMemberResp, error) {
	response := new(createWorkspaceMemberResp)
	err := r.core.rawRequest(ctx, &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/workspaces/:workspace_id/members",
		Body:   req,
	}, response)
	return response.Data, err
}

func (r *workspacesMembers) Delete(ctx context.Context, req *DeleteWorkspaceMemberReq) (*DeleteWorkspaceMemberResp, error) {
	response := new(deleteWorkspaceMemberResp)
	err := r.core.rawRequest(ctx, &RawRequestReq{
		Method: http.MethodDelete,
		URL:    "/v1/workspaces/:workspace_id/members",
		Body:   req,
	}, response)
	return response.Data, err
}

type WorkspaceMember struct {
	UserID         string            `json:"user_id"`                    // 用户ID
	RoleType       WorkspaceRoleType `json:"role_type"`                  // 当前用户角色
	UserNickname   string            `json:"user_nickname,omitempty"`    // 昵称（添加成员时不用传）
	UserUniqueName string            `json:"user_unique_name,omitempty"` // 用户名（添加成员时不用传）
	AvatarUrl      string            `json:"avatar_url,omitempty"`       // 头像 （添加成员时不用传）
}

// ListWorkspaceMemberReq ...
type ListWorkspaceMemberReq struct {
	WorkspaceID string `path:"workspace_id" json:"-"`
	PageNum     int    `query:"page_num" json:"-"`
	PageSize    int    `query:"page_size" json:"-"`
}

// ListWorkspaceMemberResp ...
type ListWorkspaceMemberResp struct {
	baseModel
	TotalCount int                `json:"total_count"`
	Items      []*WorkspaceMember `json:"items"`
}

type CreateWorkspaceMemberReq struct {
	WorkspaceID string             `path:"workspace_id" json:"-"`
	Users       []*WorkspaceMember `json:"users"`
}

type CreateWorkspaceMemberResp struct {
	baseModel
	AddedSuccessUserIDs   []string `json:"added_success_user_ids"`   // 团队或企业版成功添加的用户 ID 列表。
	InvitedSuccessUserIDs []string `json:"invited_success_user_ids"` // 个人版中，发起邀请且用户同意加入的用户 ID 列表。
	NotExistUserIDs       []string `json:"not_exist_user_ids"`       // 因用户不存在而导致添加失败的用户 ID 列表。
	AlreadyJoinedUserIDs  []string `json:"already_joined_user_ids"`  // 用户在该工作空间中已经存在，不重复添加。
	AlreadyInvitedUserIDs []string `json:"already_invited_user_ids"` // 已经发起邀请但用户还未同意加入的用户 ID 列表。
}

type DeleteWorkspaceMemberReq struct {
	WorkspaceID string   `path:"workspace_id" json:"-"`
	UserIDs     []string `json:"user_ids" sep:","`
}

type DeleteWorkspaceMemberResp struct {
	baseModel
	RemovedSuccessUserIDs        []string `json:"removed_success_user_ids"`          // 成功移除的成员列表。
	NotInWorkspaceUserIDs        []string `json:"not_in_workspace_user_ids"`         // 不在当前空间中的用户 ID 列表，这些用户不会被处理。
	OwnerNotSupportRemoveUserIDs []string `json:"owner_not_support_remove_user_ids"` // 移除失败，该用户为空间所有者。
}

type listWorkspaceMemberResp struct {
	baseResponse
	Data *ListWorkspaceMemberResp `json:"data"`
}

type createWorkspaceMemberResp struct {
	baseResponse
	Data *CreateWorkspaceMemberResp `json:"data"`
}

type deleteWorkspaceMemberResp struct {
	baseResponse
	Data *DeleteWorkspaceMemberResp `json:"data"`
}

func (r *ListWorkspaceMemberReq) toReq(request *pageRequest) *ListWorkspaceMemberReq {
	return &ListWorkspaceMemberReq{
		PageNum:     request.PageNum,
		PageSize:    request.PageSize,
		WorkspaceID: r.WorkspaceID,
	}
}

type workspacesMembers struct {
	core *core
}

func newWorkspacesMembers(core *core) *workspacesMembers {
	return &workspacesMembers{core: core}
}
