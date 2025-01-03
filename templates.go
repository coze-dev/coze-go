package coze

import (
	"context"
	"fmt"
	"net/http"
)

// templates provides access to template-related operations
type templates struct {
	core *core
}

// newTemplates creates a new Templates client
func newTemplates(core *core) *templates {
	return &templates{core: core}
}

// Duplicate creates a copy of an existing template
func (c *templates) Duplicate(ctx context.Context, templateID string, req *DuplicateTemplateReq) (*TemplateDuplicateResp, error) {
	url := fmt.Sprintf("/v1/templates/%s/duplicate", templateID)

	var resp templateDuplicateResp
	err := c.core.Request(ctx, http.MethodPost, url, req, &resp)
	if err != nil {
		return nil, err
	}
	result := resp.Data
	result.setHTTPResponse(resp.HTTPResponse)
	return result, nil
}

// TemplateEntityType represents the type of template entity
type TemplateEntityType string

const (
	// TemplateEntityTypeAgent represents an agent template
	TemplateEntityTypeAgent TemplateEntityType = "agent"
)

// TemplateDuplicateResp represents the response from duplicating a template
type TemplateDuplicateResp struct {
	baseModel
	EntityID   string             `json:"entity_id"`
	EntityType TemplateEntityType `json:"entity_type"`
}

// templateDuplicateResp represents response for creating document
type templateDuplicateResp struct {
	baseResponse
	Data *TemplateDuplicateResp `json:"data"`
}

// DuplicateTemplateReq represents the request to duplicate a template
type DuplicateTemplateReq struct {
	WorkspaceID string  `json:"workspace_id"`
	Name        *string `json:"name,omitempty"`
}
