package coze

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// variables manages variable-related API interactions.
//
// Update API docs: https://www.coze.cn/open/docs/developer_guides/update_variable
// Retrieve API docs: https://www.coze.cn/open/docs/developer_guides/read_variable
type variables struct {
	core *core
}

func newVariables(core *core) *variables {
	return &variables{core: core}
}

// VariableValue represents a single variable with its keyword and value.
type VariableValue struct {
	Keyword    string `json:"keyword"`
	Value      string `json:"value"`
	UpdateTime int64  `json:"update_time,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
}

// RetrieveVariablesReq represents the parameters for retrieving variables.
type RetrieveVariablesReq struct {
	ConnectorUID string   `json:"connector_uid"`          // Required: Unique identifier for the connector
	Keywords     []string `json:"keywords"`               // Required: List of variable keywords to retrieve
	AppID        *string  `json:"app_id,omitempty"`       // Optional: Application ID filter
	BotID        *string  `json:"bot_id,omitempty"`       // Optional: Bot ID filter
	ConnectorID  *string  `json:"connector_id,omitempty"` // Optional: Connector ID filter
}

type retrieveVariablesResp struct {
	baseResponse
	Data *RetrieveVariablesResp `json:"data"`
}

type RetrieveVariablesResp struct {
	baseModel
	Items []*VariableValue `json:"items"`
}

// Retrieve retrieves variables matching the specified criteria.
func (s *variables) Retrieve(ctx context.Context, req *RetrieveVariablesReq) (*RetrieveVariablesResp, error) {
	if req == nil {
		return nil, errors.New("invalid req")
	}
	method := http.MethodGet
	path := "/v1/variables"
	baseOpts := []RequestOption{
		withHTTPQuery("connector_uid", req.ConnectorUID),
		withHTTPQuery("keywords", strings.Join(req.Keywords, ",")),
	}
	if req.AppID != nil {
		baseOpts = append(baseOpts, withHTTPQuery("app_id", *req.AppID))
	}
	if req.BotID != nil {
		baseOpts = append(baseOpts, withHTTPQuery("bot_id", *req.BotID))
	}
	if req.ConnectorID != nil {
		baseOpts = append(baseOpts, withHTTPQuery("connector_id", *req.ConnectorID))
	}

	resp := &retrieveVariablesResp{}
	err := s.core.Request(ctx, method, path, nil, resp, baseOpts...)
	if err != nil {
		return nil, err
	}
	result := &RetrieveVariablesResp{
		baseModel: baseModel{
			httpResponse: resp.HTTPResponse,
		},
		Items: resp.Data.Items,
	}
	return result, nil
}

// UpdateVariablesReq represents the request body for updating variables.
type UpdateVariablesReq struct {
	ConnectorUID string          `json:"connector_uid"`          // Required: Unique identifier for the connector
	Data         []VariableValue `json:"data"`                   // Required: List of variable values to update
	AppID        *string         `json:"app_id,omitempty"`       // Optional: Application ID filter
	BotID        *string         `json:"bot_id,omitempty"`       // Optional: Bot ID filter
	ConnectorID  *string         `json:"connector_id,omitempty"` // Optional: Connector ID filter
}

type updateVariablesResp struct {
	baseResponse
	Data *UpdateVariablesResp `json:"data"`
}

type UpdateVariablesResp struct {
	baseModel
}

// Update updates variables with the provided data.
func (s *variables) Update(ctx context.Context, req *UpdateVariablesReq) (*UpdateVariablesResp, error) {
	if req == nil {
		return nil, errors.New("invalid req")
	}
	method := http.MethodPut
	uri := "/v1/variables"
	resp := &updateVariablesResp{}
	err := s.core.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	if resp.Data == nil {
		resp.Data = new(UpdateVariablesResp)
	}
	resp.Data.setHTTPResponse(resp.HTTPResponse)
	return resp.Data, nil
}
