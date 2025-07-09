package coze

import (
	"context"
	"net/http"
)

// Create adds enterprise members
func (r *enterprisesMembers) Create(ctx context.Context, req *CreateEnterpriseMemberReq) (*CreateEnterpriseMemberResp, error) {
	response := new(createEnterpriseMemberResp)
	err := r.core.rawRequest(ctx, &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/enterprises/:enterprise_id/members",
		Body:   req,
	}, response)
	return response.Data, err
}

// Delete removes an enterprise member
func (r *enterprisesMembers) Delete(ctx context.Context, req *DeleteEnterpriseMemberReq) (*DeleteEnterpriseMemberResp, error) {
	response := new(deleteEnterpriseMemberResp)
	err := r.core.rawRequest(ctx, &RawRequestReq{
		Method: http.MethodDelete,
		URL:    "/v1/enterprises/:enterprise_id/members/:user_id",
		Body:   req,
	}, response)
	return response.Data, err
}

// Update modifies an enterprise member's role
func (r *enterprisesMembers) Update(ctx context.Context, req *UpdateEnterpriseMemberReq) (*UpdateEnterpriseMemberResp, error) {
	response := new(updateEnterpriseMemberResp)
	err := r.core.rawRequest(ctx, &RawRequestReq{
		Method: http.MethodPut,
		URL:    "/v1/enterprises/:enterprise_id/members/:user_id",
		Body:   req,
	}, response)
	return response.Data, err
}

// EnterpriseMemberRole represents the role of an enterprise member
type EnterpriseMemberRole string

const (
	// EnterpriseMemberRoleAdmin represents an enterprise administrator
	EnterpriseMemberRoleAdmin EnterpriseMemberRole = "enterprise_admin"
	// EnterpriseMemberRoleMember represents an enterprise member
	EnterpriseMemberRoleMember EnterpriseMemberRole = "enterprise_member"
)

// EnterpriseMember represents an enterprise member
type EnterpriseMember struct {
	UserID string               `json:"user_id"` // 用户ID
	Role   EnterpriseMemberRole `json:"role"`    // 当前用户角色
}

// CreateEnterpriseMemberReq represents the request to create enterprise members
type CreateEnterpriseMemberReq struct {
	EnterpriseID string              `path:"enterprise_id" json:"-"`
	Users        []*EnterpriseMember `json:"users"`
}

// CreateEnterpriseMemberResp represents the response from creating enterprise members
type CreateEnterpriseMemberResp struct {
	baseModel
}

// DeleteEnterpriseMemberReq represents the request to delete an enterprise member
type DeleteEnterpriseMemberReq struct {
	EnterpriseID   string `path:"enterprise_id" json:"-"`
	UserID         string `path:"user_id" json:"-"`
	ReceiverUserID string `json:"receiver_user_id"`
}

// DeleteEnterpriseMemberResp represents the response from deleting an enterprise member
type DeleteEnterpriseMemberResp struct {
	baseModel
}

// UpdateEnterpriseMemberReq represents the request to update an enterprise member
type UpdateEnterpriseMemberReq struct {
	EnterpriseID string               `path:"enterprise_id" json:"-"`
	UserID       string               `path:"user_id" json:"-"`
	Role         EnterpriseMemberRole `json:"role"`
}

// UpdateEnterpriseMemberResp represents the response from updating an enterprise member
type UpdateEnterpriseMemberResp struct {
	baseModel
}

// createEnterpriseMemberResp represents the response for creating enterprise members
type createEnterpriseMemberResp struct {
	baseResponse
	Data *CreateEnterpriseMemberResp `json:"data"`
}

// deleteEnterpriseMemberResp represents the response for deleting enterprise members
type deleteEnterpriseMemberResp struct {
	baseResponse
	Data *DeleteEnterpriseMemberResp `json:"data"`
}

// updateEnterpriseMemberResp represents the response for updating enterprise members
type updateEnterpriseMemberResp struct {
	baseResponse
	Data *UpdateEnterpriseMemberResp `json:"data"`
}

type enterprisesMembers struct {
	core *core
}

func newEnterprisesMembers(core *core) *enterprisesMembers {
	return &enterprisesMembers{core: core}
}
