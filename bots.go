package coze

import (
	"context"
	"net/http"
	"strconv"
)

func (r *bots) Create(ctx context.Context, req *CreateBotsReq) (*CreateBotsResp, error) {
	method := http.MethodPost
	uri := "/v1/bot/create"
	resp := &createBotsResp{}
	err := r.core.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	resp.Data.setHTTPResponse(resp.HTTPResponse)
	return resp.Data, nil
}

func (r *bots) Update(ctx context.Context, req *UpdateBotsReq) (*UpdateBotsResp, error) {
	method := http.MethodPost
	uri := "/v1/bot/update"
	resp := &updateBotsResp{}
	err := r.core.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	result := &UpdateBotsResp{}
	result.setHTTPResponse(resp.HTTPResponse)
	return result, nil
}

func (r *bots) Publish(ctx context.Context, req *PublishBotsReq) (*PublishBotsResp, error) {
	method := http.MethodPost
	uri := "/v1/bot/publish"
	resp := &publishBotsResp{}
	err := r.core.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	resp.Data.LogID = resp.HTTPResponse.GetLogID()
	return resp.Data, nil
}

func (r *bots) Retrieve(ctx context.Context, req *RetrieveBotsReq) (*RetrieveBotsResp, error) {
	method := http.MethodGet
	uri := "/v1/bot/get_online_info"
	resp := &retrieveBotsResp{}
	err := r.core.Request(ctx, method, uri, nil, resp, withHTTPQuery("bot_id", req.BotID))
	if err != nil {
		return nil, err
	}
	resp.Bot.LogID = resp.HTTPResponse.GetLogID()
	return resp.Bot, nil
}

func (r *bots) List(ctx context.Context, req *ListBotsReq) (*NumberPaged[SimpleBot], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged[SimpleBot](
		func(request *PageRequest) (*PageResponse[SimpleBot], error) {
			uri := "/v1/space/published_bots_list"
			resp := &listBotsResp{}
			err := r.core.Request(ctx, http.MethodGet, uri, nil, resp,
				withHTTPQuery("space_id", req.SpaceID),
				withHTTPQuery("page_index", strconv.Itoa(request.PageNum)),
				withHTTPQuery("page_size", strconv.Itoa(request.PageSize)))
			if err != nil {
				return nil, err
			}
			return &PageResponse[SimpleBot]{
				Total:   resp.Data.Total,
				HasMore: len(resp.Data.Bots) >= request.PageSize,
				Data:    resp.Data.Bots,
				LogID:   resp.HTTPResponse.GetLogID(),
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
	SpaceID        string             `json:"space_id"`        // Space ID
	Name           string             `json:"name"`            // Name
	Description    string             `json:"description"`     // Description
	IconFileID     string             `json:"icon_file_id"`    // Icon file ID
	PromptInfo     *BotPromptInfo     `json:"prompt_info"`     // Prompt information
	OnboardingInfo *BotOnboardingInfo `json:"onboarding_info"` // Onboarding information
}

// CreateBotsResp 创建机器人响应
type createBotsResp struct {
	baseResponse
	Data *CreateBotsResp `json:"data"`
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

// PublishBotsResp 发布机器人响应
type publishBotsResp struct {
	baseResponse
	Data *PublishBotsResp `json:"data"`
}

type PublishBotsResp struct {
	baseModel
	BotID      string `json:"bot_id"`
	BotVersion string `json:"version"`
}

// ListBotsReq represents the request structure for listing bots
type ListBotsReq struct {
	SpaceID  string `json:"space_id"`  // Space ID
	PageNum  int    `json:"page_num"`  // Page number
	PageSize int    `json:"page_size"` // Page size
}

// listBotsResp response structure for listing bots
type listBotsResp struct {
	baseResponse
	Data struct {
		Bots  []*SimpleBot `json:"space_bots"`
		Total int          `json:"total"`
	} `json:"data"`
}

// RetrieveBotsReq represents the request structure for retrieving a bot
type RetrieveBotsReq struct {
	BotID string `json:"bot_id"` // Bot ID
}

// RetrieveBotsResp response structure for retrieving a bot
type retrieveBotsResp struct {
	baseResponse
	Bot *RetrieveBotsResp `json:"data"`
}

type RetrieveBotsResp struct {
	Bot
	baseModel
}

// UpdateBotsReq represents the request structure for updating a bot
type UpdateBotsReq struct {
	BotID          string             `json:"bot_id"`          // Bot ID
	Name           string             `json:"name"`            // Name
	Description    string             `json:"description"`     // Description
	IconFileID     string             `json:"icon_file_id"`    // Icon file ID
	PromptInfo     *BotPromptInfo     `json:"prompt_info"`     // Prompt information
	OnboardingInfo *BotOnboardingInfo `json:"onboarding_info"` // Onboarding information
	Knowledge      *BotKnowledge      `json:"knowledge"`       // Knowledge
}

// UpdateBotsResp 更新机器人响应
type updateBotsResp struct {
	baseResponse
}

type UpdateBotsResp struct {
	baseModel
}

type bots struct {
	core *core
}

func newBots(core *core) *bots {
	return &bots{core: core}
}