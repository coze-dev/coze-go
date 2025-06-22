package coze

import (
	"context"
	"net/http"
)

// Retrieve 获取用户变量的值
//
// docs: https://www.coze.cn/open/docs/developer_guides/read_variable
func (r *variables) Retrieve(ctx context.Context, req *RetrieveVariablesReq) (*RetrieveVariablesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    "/v1/variables",
		Body:   req,
	}
	response := new(retrieveVariablesResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// Update 设置用户变量的值
//
// docs: https://www.coze.cn/open/docs/developer_guides/update_variable
func (r *variables) Update(ctx context.Context, req *UpdateVariablesReq) (*UpdateVariablesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPut,
		URL:    "/v1/variables",
		Body:   req,
	}
	response := new(updateVariablesResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
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
	ConnectorUID string   `query:"connector_uid" json:"-"`    // Required: Unique identifier for the connector
	Keywords     []string `query:"keywords" sep:"," json:"-"` // Required: List of variable keywords to retrieve
	AppID        *string  `query:"app_id" json:"-"`           // Optional: Application ID filter
	BotID        *string  `query:"bot_id" json:"-"`           // Optional: Bot ID filter
	ConnectorID  *string  `query:"connector_id" json:"-"`     // Optional: Connector ID filter
}

type RetrieveVariablesResp struct {
	baseModel
	Items []*VariableValue `json:"items"`
}

// UpdateVariablesReq represents the request body for updating variables.
type UpdateVariablesReq struct {
	ConnectorUID string          `json:"connector_uid"`          // Required: Unique identifier for the connector
	Data         []VariableValue `json:"data"`                   // Required: List of variable values to update
	AppID        *string         `json:"app_id,omitempty"`       // Optional: Application ID filter
	BotID        *string         `json:"bot_id,omitempty"`       // Optional: Bot ID filter
	ConnectorID  *string         `json:"connector_id,omitempty"` // Optional: Connector ID filter
}

type UpdateVariablesResp struct {
	baseModel
}

type retrieveVariablesResp struct {
	baseResponse
	Data *RetrieveVariablesResp `json:"data"`
}

type updateVariablesResp struct {
	baseResponse
	Data *UpdateVariablesResp `json:"data"`
}

type variables struct {
	core *core
}

func newVariables(core *core) *variables {
	return &variables{core: core}
}
