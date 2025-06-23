package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspaces(t *testing.T) {
	as := assert.New(t)
	t.Run("list workspaces success", func(t *testing.T) {
		workspaces := newWorkspace(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/workspaces", req.URL.Path)
			as.Equal("1", req.URL.Query().Get("page_num"))
			as.Equal("20", req.URL.Query().Get("page_size"))
			return mockResponse(http.StatusOK, &listWorkspaceResp{
				Data: &ListWorkspaceResp{
					TotalCount: 2,
					Workspaces: []*Workspace{
						{
							ID:            "ws1",
							Name:          "Workspace 1",
							IconUrl:       "https://example.com/icon1.png",
							RoleType:      WorkspaceRoleTypeOwner,
							WorkspaceType: WorkspaceTypePersonal,
						},
						{
							ID:            "ws2",
							Name:          "Workspace 2",
							IconUrl:       "https://example.com/icon2.png",
							RoleType:      WorkspaceRoleTypeAdmin,
							WorkspaceType: WorkspaceTypeTeam,
						},
					},
				},
			})
		})))
		paged, err := workspaces.List(context.Background(), &ListWorkspaceReq{})
		as.Nil(err)
		as.NotNil(paged)
		// as.NotEmpty(paged.LogID()) // todo
		as.False(paged.HasMore())

		items := paged.Items()
		as.Len(items, 2)

		as.Equal("ws1", items[0].ID)
		as.Equal("Workspace 1", items[0].Name)
		as.Equal("https://example.com/icon1.png", items[0].IconUrl)
		as.Equal(WorkspaceRoleTypeOwner, items[0].RoleType)
		as.Equal(WorkspaceTypePersonal, items[0].WorkspaceType)

		as.Equal("ws2", items[1].ID)
		as.Equal("Workspace 2", items[1].Name)
		as.Equal("https://example.com/icon2.png", items[1].IconUrl)
		as.Equal(WorkspaceRoleTypeAdmin, items[1].RoleType)
		as.Equal(WorkspaceTypeTeam, items[1].WorkspaceType)
	})

	t.Run("list workspaces with default pagination", func(t *testing.T) {
		workspaces := newWorkspace(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal("1", req.URL.Query().Get("page_num"))
			as.Equal("20", req.URL.Query().Get("page_size"))
			return mockResponse(http.StatusOK, &listWorkspaceResp{
				Data: &ListWorkspaceResp{
					TotalCount: 0,
					Workspaces: []*Workspace{},
				},
			})
		})))
		paged, err := workspaces.List(context.Background(), NewListWorkspaceReq())
		as.Nil(err)
		as.NotNil(paged)
		// as.NotEmpty(paged.LogID()) // todo
		as.False(paged.HasMore())
		as.Empty(paged.Items())
	})

	t.Run("list workspaces with error", func(t *testing.T) {
		workspaces := newWorkspace(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("test error")
		})))
		_, err := workspaces.List(context.Background(), &ListWorkspaceReq{
			PageNum:  1,
			PageSize: 20,
		})
		as.NotNil(err)
	})
}
