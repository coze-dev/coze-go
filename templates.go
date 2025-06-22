package coze

import (
	"context"
	"net/http"
)

// Duplicate creates a copy of an existing template
func (r *templates) Duplicate(ctx context.Context, templateID string, req *DuplicateTemplateReq) (*TemplateDuplicateResp, error) {
	if req == nil {
		req = &DuplicateTemplateReq{}
	}
	req.TemplateID = templateID
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/templates/:template_id/duplicate",
		Body:   req,
	}
	response := new(templateDuplicateResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
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

// DuplicateTemplateReq represents the request to duplicate a template
type DuplicateTemplateReq struct {
	TemplateID  string  `path:"template_id" json:"-"`
	WorkspaceID string  `json:"workspace_id,omitempty"`
	Name        *string `json:"name,omitempty"`
}

// templateDuplicateResp represents response for creating document
type templateDuplicateResp struct {
	baseResponse
	Data *TemplateDuplicateResp `json:"data"`
}

type templates struct {
	core *core
}

func newTemplates(core *core) *templates {
	return &templates{core: core}
}
