package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockWorkflow(id string) *WorkflowInfo {
	return &WorkflowInfo{
		WorkflowID:   id,
		WorkflowName: randomString(10),
		Description:  randomString(10),
		IconURL:      randomString(10),
		AppID:        randomString(10),
	}
}

func TestWorkflow(t *testing.T) {
	as := assert.New(t)
	t.Run("list", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			workflows := newWorkflows(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/workflows", req.URL.Path)
				as.Equal("1", req.URL.Query().Get("page_num"))
				as.Equal("20", req.URL.Query().Get("page_size"))
				return mockResponse(http.StatusOK, &listWorkflowResp{
					Data: &ListWorkflowResp{
						HasMore: false,
						Items: []*WorkflowInfo{
							mockWorkflow("1"),
							mockWorkflow("2"),
						},
					},
				})
			})))
			resp, err := workflows.List(context.Background(), &ListWorkflowReq{})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.False(resp.HasMore())

			items := resp.Items()
			as.Len(items, 2)

			as.Equal("1", items[0].WorkflowID)
			as.Equal("2", items[1].WorkflowID)
		})

		t.Run("error", func(t *testing.T) {
			workflows := newWorkflows(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := workflows.List(context.Background(), &ListWorkflowReq{})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})
}
