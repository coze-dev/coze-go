package coze

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspacesMembers(t *testing.T) {
	as := assert.New(t)
	t.Run("list", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			workspaceID := randomString(10)
			members := newWorkspacesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal(fmt.Sprintf("/v1/workspaces/%s/members", workspaceID), req.URL.Path)
				as.Equal("1", req.URL.Query().Get("page_num"))
				as.Equal("20", req.URL.Query().Get("page_size"))
				return mockResponse(http.StatusOK, &listWorkspaceMemberResp{
					Data: &ListWorkspaceMemberResp{
						TotalCount: 2,
						Items: []*WorkspaceMember{
							{
								UserID:         "ws1",
								UserNickname:   "User 1",
								UserUniqueName: "user1",
								AvatarUrl:      "https://example.com/icon1.png",
								RoleType:       WorkspaceRoleTypeOwner,
							},
							{
								UserID:         "ws2",
								UserNickname:   "User 2",
								UserUniqueName: "user2",
								AvatarUrl:      "https://example.com/icon2.png",
								RoleType:       WorkspaceRoleTypeAdmin,
							},
						},
					},
				})
			})))
			resp, err := members.List(context.Background(), &ListWorkspaceMemberReq{
				WorkspaceID: workspaceID,
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.False(resp.HasMore())

			items := resp.Items()
			as.Len(items, 2)

			as.Equal("ws1", items[0].UserID)
			as.Equal("User 1", items[0].UserNickname)
			as.Equal("https://example.com/icon1.png", items[0].AvatarUrl)
			as.Equal(WorkspaceRoleTypeOwner, items[0].RoleType)

			as.Equal("ws2", items[1].UserID)
			as.Equal("User 2", items[1].UserNickname)
			as.Equal("https://example.com/icon2.png", items[1].AvatarUrl)
			as.Equal(WorkspaceRoleTypeAdmin, items[1].RoleType)
		})

		t.Run("error", func(t *testing.T) {
			workspaceID := randomString(10)
			members := newWorkspacesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := members.List(context.Background(), &ListWorkspaceMemberReq{
				WorkspaceID: workspaceID,
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("create", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			workspaceID := randomString(10)
			members := newWorkspacesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodPost, req.Method)
				as.Equal(fmt.Sprintf("/v1/workspaces/%s/members", workspaceID), req.URL.Path)
				return mockResponse(http.StatusOK, &createWorkspaceMemberResp{
					Data: &CreateWorkspaceMemberResp{
						AddedSuccessUserIDs:   []string{"ws1", "ws2"},
						InvitedSuccessUserIDs: []string{"ws3", "ws4"},
						NotExistUserIDs:       []string{"ws5", "ws6"},
						AlreadyJoinedUserIDs:  []string{"ws7", "ws8"},
						AlreadyInvitedUserIDs: []string{"ws9", "ws10"},
					},
				})
			})))
			resp, err := members.Create(context.Background(), &CreateWorkspaceMemberReq{
				WorkspaceID: workspaceID,
				Users: []*WorkspaceMember{
					{
						UserID:         "ws1",
						UserNickname:   "User 1",
						UserUniqueName: "user1",
						AvatarUrl:      "https://example.com/icon1.png",
						RoleType:       WorkspaceRoleTypeOwner,
					},
					{
						UserID:         "ws2",
						UserNickname:   "User 2",
						UserUniqueName: "user2",
						AvatarUrl:      "https://example.com/icon2.png",
						RoleType:       WorkspaceRoleTypeAdmin,
					},
				},
			})
			as.Nil(err)
			as.NotNil(resp)
			as.Equal("ws1", resp.AddedSuccessUserIDs[0])
			as.Equal("ws2", resp.AddedSuccessUserIDs[1])
			as.Equal("ws3", resp.InvitedSuccessUserIDs[0])
			as.Equal("ws4", resp.InvitedSuccessUserIDs[1])
			as.Equal("ws5", resp.NotExistUserIDs[0])
			as.Equal("ws6", resp.NotExistUserIDs[1])
			as.Equal("ws7", resp.AlreadyJoinedUserIDs[0])
			as.Equal("ws8", resp.AlreadyJoinedUserIDs[1])
			as.Equal("ws9", resp.AlreadyInvitedUserIDs[0])
			as.Equal("ws10", resp.AlreadyInvitedUserIDs[1])
		})

		t.Run("error", func(t *testing.T) {
			workspaceID := randomString(10)
			members := newWorkspacesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := members.Create(context.Background(), &CreateWorkspaceMemberReq{
				WorkspaceID: workspaceID,
				Users: []*WorkspaceMember{
					{
						UserID:         "ws1",
						UserNickname:   "User 1",
						UserUniqueName: "user1",
						AvatarUrl:      "https://example.com/icon1.png",
						RoleType:       WorkspaceRoleTypeOwner,
					},
				},
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			workspaceID := randomString(10)
			members := newWorkspacesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodDelete, req.Method)
				as.Equal(fmt.Sprintf("/v1/workspaces/%s/members", workspaceID), req.URL.Path)
				return mockResponse(http.StatusOK, &deleteWorkspaceMemberResp{
					Data: &DeleteWorkspaceMemberResp{
						RemovedSuccessUserIDs:        []string{"ws1", "ws2"},
						NotInWorkspaceUserIDs:        []string{"ws3", "ws4"},
						OwnerNotSupportRemoveUserIDs: []string{"ws5", "ws6"},
					},
				})
			})))
			resp, err := members.Delete(context.Background(), &DeleteWorkspaceMemberReq{
				WorkspaceID: workspaceID,
				UserIDs:     []string{"ws1", "ws2"},
			})
			as.Nil(err)
			as.NotNil(resp)
			as.Equal("ws1", resp.RemovedSuccessUserIDs[0])
			as.Equal("ws2", resp.RemovedSuccessUserIDs[1])
			as.Equal("ws3", resp.NotInWorkspaceUserIDs[0])
			as.Equal("ws4", resp.NotInWorkspaceUserIDs[1])
			as.Equal("ws5", resp.OwnerNotSupportRemoveUserIDs[0])
			as.Equal("ws6", resp.OwnerNotSupportRemoveUserIDs[1])
		})

		t.Run("error", func(t *testing.T) {
			workspaceID := randomString(10)
			members := newWorkspacesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := members.Delete(context.Background(), &DeleteWorkspaceMemberReq{
				WorkspaceID: workspaceID,
				UserIDs:     []string{"ws1", "ws2"},
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})
}
