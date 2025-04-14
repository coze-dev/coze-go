package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBots(t *testing.T) {
	t.Run("Create bot success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/bot/create", req.URL.Path)
				return mockResponse(http.StatusOK, &createBotsResp{
					Data: &CreateBotsResp{
						BotID: "test_bot_id",
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		bots := newBots(core)

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

		require.NoError(t, err)
		assert.Equal(t, "test_bot_id", resp.BotID)
		assert.Equal(t, "test_log_id", resp.LogID())
	})

	t.Run("Update bot success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/bot/update", req.URL.Path)
				return mockResponse(http.StatusOK, &updateBotsResp{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		bots := newBots(core)

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

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
	})

	t.Run("Publish bot success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/bot/publish", req.URL.Path)
				return mockResponse(http.StatusOK, &publishBotsResp{
					Data: &PublishBotsResp{
						BotID:      "test_bot_id",
						BotVersion: "1.0.0",
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		bots := newBots(core)

		resp, err := bots.Publish(context.Background(), &PublishBotsReq{
			BotID:        "test_bot_id",
			ConnectorIDs: []string{"connector1", "connector2"},
		})

		require.NoError(t, err)
		assert.Equal(t, "test_bot_id", resp.BotID)
		assert.Equal(t, "1.0.0", resp.BotVersion)
		assert.Equal(t, "test_log_id", resp.LogID())
	})

	t.Run("Retrieve bot success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/bot/get_online_info", req.URL.Path)
				assert.Equal(t, "test_bot_id", req.URL.Query().Get("bot_id"))
				return mockResponse(http.StatusOK, &retrieveBotsResp{
					Bot: &RetrieveBotsResp{
						Bot: Bot{
							BotID:       "test_bot_id",
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
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		bots := newBots(core)

		resp, err := bots.Retrieve(context.Background(), &RetrieveBotsReq{
			BotID: "test_bot_id",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_bot_id", resp.Bot.BotID)
		assert.Equal(t, "Test Bot", resp.Bot.Name)
		assert.Equal(t, "1.0.0", resp.Bot.Version)
		assert.Equal(t, BotModeMultiAgent, resp.Bot.BotMode)
		assert.Equal(t, "test_log_id", resp.LogID())
	})

	t.Run("List bots success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/space/published_bots_list", req.URL.Path)
				assert.Equal(t, "test_space_id", req.URL.Query().Get("space_id"))
				assert.Equal(t, "1", req.URL.Query().Get("page_index"))
				assert.Equal(t, "20", req.URL.Query().Get("page_size"))
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
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		bots := newBots(core)

		paged, err := bots.List(context.Background(), &ListBotsReq{
			SpaceID:  "test_space_id",
			PageNum:  1,
			PageSize: 20,
		})

		require.NoError(t, err)
		assert.Len(t, paged.Items(), 2)
		items := paged.Items()
		assert.Equal(t, 2, len(items))
		assert.Equal(t, "bot1", items[0].BotID)
		assert.Equal(t, "Bot 1", items[0].BotName)
		assert.Equal(t, "bot2", items[1].BotID)
		assert.Equal(t, "Bot 2", items[1].BotName)
		assert.Nil(t, paged.Err())
	})

	t.Run("List bots with default pagination", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "1", req.URL.Query().Get("page_index"))
				assert.Equal(t, "20", req.URL.Query().Get("page_size"))
				return mockResponse(http.StatusOK, &listBotsResp{
					Data: struct {
						Bots  []*SimpleBot `json:"space_bots"`
						Total int          `json:"total"`
					}{
						Bots:  []*SimpleBot{},
						Total: 0,
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		bots := newBots(core)

		paged, err := bots.List(context.Background(), &ListBotsReq{
			SpaceID: "test_space_id",
		})

		require.NoError(t, err)
		assert.Empty(t, paged.Items())
	})
}

func TestBotMode(t *testing.T) {
	t.Run("BotMode constants", func(t *testing.T) {
		assert.Equal(t, BotMode(1), BotModeMultiAgent)
		assert.Equal(t, BotMode(0), BotModeSingleAgentWorkflow)
	})
}
