package coze

import (
	"context"
	"net/http"
)

// Create 创建智能体
//
// docs: https://www.coze.cn/open/docs/developer_guides/create_bot
func (r *bots) Create(ctx context.Context, req *CreateBotsReq) (*CreateBotsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/bot/create",
		Body:   req,
	}
	response := new(createBotsResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// Update 更新智能体
//
// docs: https://www.coze.cn/open/docs/developer_guides/update_bot
func (r *bots) Update(ctx context.Context, req *UpdateBotsReq) (*UpdateBotsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/bot/update",
		Body:   req,
	}
	response := new(updateBotsResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// Publish 发布智能体
//
// docs: https://www.coze.cn/open/docs/developer_guides/publish_bot
func (r *bots) Publish(ctx context.Context, req *PublishBotsReq) (*PublishBotsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/bot/publish",
		Body:   req,
	}
	response := new(publishBotsResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// Retrieve 获取已发布智能体配置（即将下线）
//
// docs: https://www.coze.cn/open/docs/developer_guides/get_metadata
func (r *bots) Retrieve(ctx context.Context, req *RetrieveBotsReq) (*RetrieveBotsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    "/v1/bot/get_online_info",
		Body:   req,
	}
	response := new(retrieveBotsResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// List 查看已发布智能体列表（即将下线）
//
// docs: https://www.coze.cn/open/docs/developer_guides/published_bots_list
func (r *bots) List(ctx context.Context, req *ListBotsReq) (NumberPaged[SimpleBot], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged[SimpleBot](
		func(request *pageRequest) (*pageResponse[SimpleBot], error) {
			response := new(listBotsResp)
			err := r.core.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/space/published_bots_list",
				Body: &ListBotsReq{
					SpaceID:  req.SpaceID,
					PageNum:  request.PageNum,
					PageSize: request.PageSize,
				},
			}, response)
			if err != nil {
				return nil, err
			}
			return &pageResponse[SimpleBot]{
				Total:   response.Data.Total,
				HasMore: len(response.Data.Bots) >= request.PageSize,
				Data:    response.Data.Bots,
				LogID:   response.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

// BotMode represents the bot mode
type BotMode int

const (
	BotModeMultiAgent          BotMode = 1
	BotModeSingleAgentWorkflow BotMode = 0
)

// Bot represents complete bot information
type Bot struct {
	BotID          string             `json:"bot_id"`
	Name           string             `json:"name"`
	Description    string             `json:"description,omitempty"`
	IconURL        string             `json:"icon_url,omitempty"`
	CreateTime     int64              `json:"create_time"`
	UpdateTime     int64              `json:"update_time"`
	Version        string             `json:"version,omitempty"`
	PromptInfo     *BotPromptInfo     `json:"prompt_info,omitempty"`
	OnboardingInfo *BotOnboardingInfo `json:"onboarding_info,omitempty"`
	BotMode        BotMode            `json:"bot_mode"`
	PluginInfoList []*BotPluginInfo   `json:"plugin_info_list,omitempty"`
	ModelInfo      *BotModelInfo      `json:"model_info,omitempty"`
}

// SimpleBot represents simplified bot information
type SimpleBot struct {
	BotID       string `json:"bot_id"`
	BotName     string `json:"bot_name"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
	PublishTime string `json:"publish_time,omitempty"`
}

// BotKnowledge represents bot knowledge base configuration
type BotKnowledge struct {
	DatasetIDs     []string `json:"dataset_ids"`
	AutoCall       bool     `json:"auto_call"`
	SearchStrategy int      `json:"search_strategy"`
}

// BotModelInfo represents bot model information
type BotModelInfo struct {
	ModelID   string `json:"model_id"`
	ModelName string `json:"model_name"`
}

type BotModelInfoConfig struct {
	TopK             int     `json:"top_k,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	ModelID          string  `json:"model_id"`
	MaxTokens        int     `json:"max_tokens,omitempty"`
	Temperature      float64 `json:"temperature,omitempty"`
	ContextRound     int     `json:"context_round,omitempty"`
	ResponseFormat   string  `json:"response_format,omitempty"` // text,markdown,json
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
}

// WorkflowIDList represents workflow ID information
type WorkflowIDList struct {
	IDs []WorkflowIDInfo `json:"ids"`
}

type WorkflowIDInfo struct {
	ID string `json:"id"`
}

// BotOnboardingInfo represents bot onboarding information
type BotOnboardingInfo struct {
	Prologue           string   `json:"prologue,omitempty"`
	SuggestedQuestions []string `json:"suggested_questions,omitempty"`
}

// BotPluginAPIInfo represents bot plugin API information
type BotPluginAPIInfo struct {
	APIID       string `json:"api_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// BotPluginInfo represents bot plugin information
type BotPluginInfo struct {
	PluginID    string              `json:"plugin_id"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	IconURL     string              `json:"icon_url,omitempty"`
	APIInfoList []*BotPluginAPIInfo `json:"api_info_list,omitempty"`
}

// BotPromptInfo represents bot prompt information
type BotPromptInfo struct {
	Prompt string `json:"prompt"`
}

type CreateBotsReq struct {
	SpaceID         string              `json:"space_id"`          // Space ID
	Name            string              `json:"name"`              // Name
	Description     string              `json:"description"`       // Description
	IconFileID      string              `json:"icon_file_id"`      // Icon file ID
	PromptInfo      *BotPromptInfo      `json:"prompt_info"`       // Prompt information
	OnboardingInfo  *BotOnboardingInfo  `json:"onboarding_info"`   // Onboarding information
	ModelInfoConfig *BotModelInfoConfig `json:"model_info_config"` // ModelInfoConfig information
	WorkflowIDList  *WorkflowIDList     `json:"workflow_id_list"`  // WorkflowIDList information
}

type CreateBotsResp struct {
	baseModel
	BotID string `json:"bot_id"`
}

// PublishBotsReq represents the request structure for publishing a bot
type PublishBotsReq struct {
	BotID        string   `json:"bot_id"`        // Bot ID
	ConnectorIDs []string `json:"connector_ids"` // Connector ID list
}

type PublishBotsResp struct {
	baseModel
	BotID      string `json:"bot_id"`
	BotVersion string `json:"version"`
}

// ListBotsReq represents the request structure for listing bots
type ListBotsReq struct {
	SpaceID  string `query:"space_id" json:"-"`   // Space ID
	PageNum  int    `query:"page_index" json:"-"` // Page number
	PageSize int    `query:"page_size" json:"-"`  // Page size
}

// RetrieveBotsReq represents the request structure for retrieving a bot
type RetrieveBotsReq struct {
	BotID string `query:"bot_id" json:"-"` // Bot ID
}

// RetrieveBotsResp response structure for retrieving a bot
type retrieveBotsResp struct {
	baseResponse
	Data *RetrieveBotsResp `json:"data"`
}

type RetrieveBotsResp struct {
	Bot
	baseModel
}

// UpdateBotsReq represents the request structure for updating a bot
type UpdateBotsReq struct {
	BotID           string              `json:"bot_id"`            // Bot ID
	Name            string              `json:"name"`              // Name
	Description     string              `json:"description"`       // Description
	IconFileID      string              `json:"icon_file_id"`      // Icon file ID
	PromptInfo      *BotPromptInfo      `json:"prompt_info"`       // Prompt information
	OnboardingInfo  *BotOnboardingInfo  `json:"onboarding_info"`   // Onboarding information
	Knowledge       *BotKnowledge       `json:"knowledge"`         // Knowledge
	ModelInfoConfig *BotModelInfoConfig `json:"model_info_config"` // ModelInfoConfig information
	WorkflowIDList  *WorkflowIDList     `json:"workflow_id_list"`  // WorkflowIDList information
}

type UpdateBotsResp struct {
	baseModel
}

type createBotsResp struct {
	baseResponse
	Data *CreateBotsResp `json:"data"`
}

type publishBotsResp struct {
	baseResponse
	Data *PublishBotsResp `json:"data"`
}

type updateBotsResp struct {
	baseResponse
	Data *UpdateBotsResp `json:"data"`
}

type listBotsResp struct {
	baseResponse
	Data struct {
		Bots  []*SimpleBot `json:"space_bots"`
		Total int          `json:"total"`
	} `json:"data"`
}

type bots struct {
	core *core
}

func newBots(core *core) *bots {
	return &bots{core: core}
}
