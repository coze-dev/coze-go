package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspaces(t *testing.T) {
	// Test List method
	t.Run("List workspaces success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/workspaces", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "1", req.URL.Query().Get("page_num"))
				assert.Equal(t, "20", req.URL.Query().Get("page_size"))

				// Return mock response
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
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		workspaces := newWorkspace(core)

		paged, err := workspaces.List(context.Background(), &ListWorkspaceReq{
			PageNum:  1,
			PageSize: 20,
		})

		require.NoError(t, err)
		assert.False(t, paged.HasMore())
		items := paged.Items()
		require.Len(t, items, 2)

		// Verify first workspace
		assert.Equal(t, "ws1", items[0].ID)
		assert.Equal(t, "Workspace 1", items[0].Name)
		assert.Equal(t, "https://example.com/icon1.png", items[0].IconUrl)
		assert.Equal(t, WorkspaceRoleTypeOwner, items[0].RoleType)
		assert.Equal(t, WorkspaceTypePersonal, items[0].WorkspaceType)

		// Verify second workspace
		assert.Equal(t, "ws2", items[1].ID)
		assert.Equal(t, "Workspace 2", items[1].Name)
		assert.Equal(t, "https://example.com/icon2.png", items[1].IconUrl)
		assert.Equal(t, WorkspaceRoleTypeAdmin, items[1].RoleType)
		assert.Equal(t, WorkspaceTypeTeam, items[1].WorkspaceType)
	})

	// Test List method with default pagination
	t.Run("List workspaces with default pagination", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify default pagination parameters
				assert.Equal(t, "1", req.URL.Query().Get("page_num"))
				assert.Equal(t, "20", req.URL.Query().Get("page_size"))

				// Return mock response
				return mockResponse(http.StatusOK, &listWorkspaceResp{
					Data: &ListWorkspaceResp{
						TotalCount: 0,
						Workspaces: []*Workspace{},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		workspaces := newWorkspace(core)

		paged, err := workspaces.List(context.Background(), NewListWorkspaceReq())

		require.NoError(t, err)
		assert.False(t, paged.HasMore())
		assert.Empty(t, paged.Items())
	})

	// Test List method with error
	t.Run("List workspaces with error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		workspaces := newWorkspace(core)

		paged, err := workspaces.List(context.Background(), &ListWorkspaceReq{
			PageNum:  1,
			PageSize: 20,
		})

		require.Error(t, err)
		assert.Nil(t, paged)
	})
}
