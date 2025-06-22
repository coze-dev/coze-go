package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplates(t *testing.T) {
	as := assert.New(t)

	t.Run("create success", func(t *testing.T) {
		templateID := randomString(10)
		workspaceID := randomString(10)
		templates := newTemplates(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/templates/"+templateID+"/duplicate", req.URL.Path)
			return mockResponse(http.StatusOK, &templateDuplicateResp{
				Data: &TemplateDuplicateResp{
					EntityID:   templateID,
					EntityType: TemplateEntityTypeAgent,
				},
			})
		})))

		resp, err := templates.Duplicate(context.Background(), templateID, &DuplicateTemplateReq{
			WorkspaceID: workspaceID,
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal(templateID, resp.EntityID)
		as.Equal(TemplateEntityTypeAgent, resp.EntityType)
	})
}
