package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBots(t *testing.T) {
	as := assert.New(t)
	t.Run("Create bot success", func(t *testing.T) {
		botID := randomString(10)
		bots := newBots(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/bot/create", req.URL.Path)
			return mockResponse(http.StatusOK, &createBotsResp{
				Data: &CreateBotsResp{
					BotID: botID,
				},
			})
		})))
		resp, err := bots.Create(context.Background(), &CreateBotsReq{
			SpaceID:     "test_space_id",
			Name:        "Test Bot",
			Description: "Test Description",
			IconFileID:  "test_icon_id",
			PromptInfo: &BotPromptInfo{
				Prompt: "Test Prompt",
			},
			OnboardingInfo: &BotOnboardingInfo{
				Prologue:           "Test Prologue",
				SuggestedQuestions: []string{"Q1", "Q2"},
			},
			ModelInfoConfig: &BotModelInfoConfig{
				ModelID: "test_model_id",
			},
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal(botID, resp.BotID)
	})

	t.Run("update bot", func(t *testing.T) {
		bots := newBots(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/bot/update", req.URL.Path)
			return mockResponse(http.StatusOK, &updateBotsResp{})
		})))
		resp, err := bots.Update(context.Background(), &UpdateBotsReq{
			BotID:       "test_bot_id",
			Name:        "Updated Bot",
			Description: "Updated Description",
			IconFileID:  "updated_icon_id",
			PromptInfo: &BotPromptInfo{
				Prompt: "Updated Prompt",
			},
			OnboardingInfo: &BotOnboardingInfo{
				Prologue:           "Updated Prologue",
				SuggestedQuestions: []string{"Q3", "Q4"},
			},
			Knowledge: &BotKnowledge{
				DatasetIDs:     []string{"dataset1", "dataset2"},
				AutoCall:       true,
				SearchStrategy: 1,
			},
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
	})

	t.Run("publish bot", func(t *testing.T) {
		botID := randomString(10)
		bots := newBots(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/bot/publish", req.URL.Path)
			return mockResponse(http.StatusOK, &publishBotsResp{
				Data: &PublishBotsResp{
					BotID:      botID,
					BotVersion: "1.0.0",
				},
			})
		})))
		resp, err := bots.Publish(context.Background(), &PublishBotsReq{
			BotID:        botID,
			ConnectorIDs: []string{"connector1", "connector2"},
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal(botID, resp.BotID)
		as.Equal("1.0.0", resp.BotVersion)
	})

	t.Run("retrieve bot", func(t *testing.T) {
		botID := randomString(10)
		bots := newBots(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/bot/get_online_info", req.URL.Path)
			as.Equal(botID, req.URL.Query().Get("bot_id"))
			return mockResponse(http.StatusOK, &retrieveBotsResp{
				Data: &RetrieveBotsResp{
					Bot: Bot{
						BotID:       botID,
						Name:        "Test Bot",
						Description: "Test Description",
						IconURL:     "https://example.com/icon.png",
						CreateTime:  1234567890,
						UpdateTime:  1234567891,
						Version:     "1.0.0",
						BotMode:     BotModeMultiAgent,
						PromptInfo: &BotPromptInfo{
							Prompt: "Test Prompt",
						},
						OnboardingInfo: &BotOnboardingInfo{
							Prologue:           "Test Prologue",
							SuggestedQuestions: []string{"Q1", "Q2"},
						},
						PluginInfoList: []*BotPluginInfo{
							{
								PluginID:    "plugin1",
								Name:        "Plugin 1",
								Description: "Plugin Description",
								IconURL:     "https://example.com/plugin-icon.png",
								APIInfoList: []*BotPluginAPIInfo{
									{
										APIID:       "api1",
										Name:        "API 1",
										Description: "API Description",
									},
								},
							},
						},
						ModelInfo: &BotModelInfo{
							ModelID:   "model1",
							ModelName: "Model 1",
						},
					},
				},
			})
		})))
		resp, err := bots.Retrieve(context.Background(), &RetrieveBotsReq{
			BotID: botID,
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal(botID, resp.Bot.BotID)
		as.Equal("Test Bot", resp.Bot.Name)
		as.Equal("1.0.0", resp.Bot.Version)
		as.Equal(BotModeMultiAgent, resp.Bot.BotMode)
	})

	t.Run("list bots", func(t *testing.T) {
		workspaceID := randomString(10)
		bots := newBots(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/space/published_bots_list", req.URL.Path)
			as.Equal(workspaceID, req.URL.Query().Get("space_id"))
			as.Equal("1", req.URL.Query().Get("page_index"))
			as.Equal("20", req.URL.Query().Get("page_size"))
			return mockResponse(http.StatusOK, &listBotsResp{
				Data: struct {
					Bots  []*SimpleBot `json:"space_bots"`
					Total int          `json:"total"`
				}{
					Bots: []*SimpleBot{
						{
							BotID:       "bot1",
							BotName:     "Bot 1",
							Description: "Description 1",
							IconURL:     "https://example.com/icon1.png",
							PublishTime: "2024-01-01",
						},
						{
							BotID:       "bot2",
							BotName:     "Bot 2",
							Description: "Description 2",
							IconURL:     "https://example.com/icon2.png",
							PublishTime: "2024-01-02",
						},
					},
					Total: 2,
				},
			})
		})))
		paged, err := bots.List(context.Background(), &ListBotsReq{
			SpaceID:  workspaceID,
			PageNum:  1,
			PageSize: 20,
		})
		as.Nil(err)
		as.NotNil(paged)
		as.Equal(2, len(paged.Items()))
		items := paged.Items()
		as.Equal("bot1", items[0].BotID)
		as.Equal("Bot 1", items[0].BotName)
		as.Equal("bot2", items[1].BotID)
		as.Equal("Bot 2", items[1].BotName)
		as.Nil(paged.Err())
	})

	t.Run("List bots with default pagination", func(t *testing.T) {
		bots := newBots(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal("1", req.URL.Query().Get("page_index"))
			as.Equal("20", req.URL.Query().Get("page_size"))
			return mockResponse(http.StatusOK, &listBotsResp{
				Data: struct {
					Bots  []*SimpleBot `json:"space_bots"`
					Total int          `json:"total"`
				}{
					Bots:  []*SimpleBot{},
					Total: 0,
				},
			})
		})))
		paged, err := bots.List(context.Background(), &ListBotsReq{
			SpaceID: "test_space_id",
		})

		require.NoError(t, err)
		as.Empty(paged.Items())
	})
}

func TestBotMode(t *testing.T) {
	as := assert.New(t)
	t.Run("BotMode constants", func(t *testing.T) {
		as.Equal(BotMode(1), BotModeMultiAgent)
		as.Equal(BotMode(0), BotModeSingleAgentWorkflow)
	})
}
