package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkflowsList(t *testing.T) {
	as := assert.New(t)

	t.Run("list workflows success", func(t *testing.T) {
		workflows := newWorkflows(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/workflows", req.URL.Path)
			as.Equal("1", req.URL.Query().Get("page_num"))
			as.Equal("20", req.URL.Query().Get("page_size"))
			return mockResponse(http.StatusOK, &listWorkflowsResp{
				Data: &listWorkflowsData{
					Items: []*WorkflowInfo{
						{
							WorkflowID:   "w1",
							WorkflowName: "Test Workflow 1",
							Description:  "A test workflow",
							IconURL:      "https://example.com/icon1.png",
							AppID:        "app1",
						},
						{
							WorkflowID:   "w2",
							WorkflowName: "Test Workflow 2",
							Description:  "Another test workflow",
							IconURL:      "https://example.com/icon2.png",
							AppID:        "app2",
						},
					},
					HasMore: false,
				},
			})
		})))

		paged, err := workflows.List(context.Background(), &ListWorkflowsReq{
			PageNum:  1,
			PageSize: 20,
		})
		as.NoError(err)
		as.NotNil(paged)
		as.NotEmpty(paged.Response().LogID())

		// Test iteration
		var items []*WorkflowInfo
		for paged.Next() {
			items = append(items, paged.Current())
		}
		as.NoError(paged.Err())
		as.Len(items, 2)
		as.Equal("w1", items[0].WorkflowID)
		as.Equal("Test Workflow 1", items[0].WorkflowName)
		as.Equal("w2", items[1].WorkflowID)
		as.Equal("Test Workflow 2", items[1].WorkflowName)
	})

	t.Run("list workflows with filters", func(t *testing.T) {
		workflowMode := WorkflowModeWorkflow
		publishStatus := PublishStatusPublished
		workflows := newWorkflows(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/workflows", req.URL.Path)
			as.Equal("workspace123", req.URL.Query().Get("workspace_id"))
			as.Equal("workflow", req.URL.Query().Get("workflow_mode"))
			as.Equal("app123", req.URL.Query().Get("app_id"))
			as.Equal("published", req.URL.Query().Get("publish_status"))
			return mockResponse(http.StatusOK, &listWorkflowsResp{
				Data: &listWorkflowsData{
					Items:   []*WorkflowInfo{},
					HasMore: false,
				},
			})
		})))

		paged, err := workflows.List(context.Background(), &ListWorkflowsReq{
			WorkspaceID:   "workspace123",
			WorkflowMode:  &workflowMode,
			AppID:         "app123",
			PublishStatus: &publishStatus,
			PageNum:       1,
			PageSize:      20,
		})
		as.NoError(err)
		as.NotNil(paged)
	})

	t.Run("list workflows with defaults", func(t *testing.T) {
		workflows := newWorkflows(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/workflows", req.URL.Path)
			as.Equal("1", req.URL.Query().Get("page_num"))
			as.Equal("20", req.URL.Query().Get("page_size"))
			return mockResponse(http.StatusOK, &listWorkflowsResp{
				Data: &listWorkflowsData{
					Items:   []*WorkflowInfo{},
					HasMore: false,
				},
			})
		})))

		paged, err := workflows.List(context.Background(), nil)
		as.NoError(err)
		as.NotNil(paged)
	})

	t.Run("list workflows error", func(t *testing.T) {
		workflows := newWorkflows(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return nil, assert.AnError
		})))

		paged, err := workflows.List(context.Background(), &ListWorkflowsReq{})
		as.Error(err)
		as.Nil(paged)
	})
}