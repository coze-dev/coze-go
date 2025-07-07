package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

// This example demonstrates how to create and update bots with thinking_type support
func main() {
	// Get an access_token through personal access token or oauth
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	workspaceID := os.Getenv("WORKSPACE_ID")
	ctx := context.Background()

	// Example 1: Create a bot with thinking enabled for Claude models
	fmt.Println("=== Creating bot with Claude thinking enabled ===")
	createResp, err := cozeCli.Bots.Create(ctx, &coze.CreateBotsReq{
		SpaceID:     workspaceID,
		Name:        "Claude Thinking Bot",
		Description: "A bot with thinking capabilities using Claude model",
		PromptInfo: &coze.BotPromptInfo{
			Prompt: "You are a helpful assistant with thinking capabilities.",
		},
		ModelInfoConfig: &coze.BotModelInfoConfig{
			ModelID:        "claude-3-5-sonnet",
			ResponseFormat: coze.ResponseFormatMarkdown,
			Temperature:    0.7,
			MaxTokens:      4000,
			Parameters: map[string]string{
				"thinking_type":          "enable",
				"thinking_budget_tokens": "2000",
			},
		},
	})
	if err != nil {
		fmt.Println("Error creating Claude bot:", err)
		return
	}
	claudeBotID := createResp.BotID
	fmt.Printf("Created Claude bot with ID: %s\n", claudeBotID)
	fmt.Printf("Log ID: %s\n\n", createResp.LogID())

	// Example 2: Create a bot with auto thinking for Doubao models
	fmt.Println("=== Creating bot with Doubao auto thinking ===")
	createResp2, err := cozeCli.Bots.Create(ctx, &coze.CreateBotsReq{
		SpaceID:     workspaceID,
		Name:        "Doubao Auto Thinking Bot",
		Description: "A bot with auto thinking capabilities using Doubao model",
		PromptInfo: &coze.BotPromptInfo{
			Prompt: "You are an intelligent assistant that can think automatically when needed.",
		},
		ModelInfoConfig: &coze.BotModelInfoConfig{
			ModelID:        "doubao-pro-128k",
			ResponseFormat: coze.ResponseFormatJSON,
			Temperature:    0.5,
			MaxTokens:      8000,
			Parameters: map[string]string{
				"thinking_type": "auto",
			},
		},
	})
	if err != nil {
		fmt.Println("Error creating Doubao bot:", err)
		return
	}
	doubaoBotID := createResp2.BotID
	fmt.Printf("Created Doubao bot with ID: %s\n", doubaoBotID)
	fmt.Printf("Log ID: %s\n\n", createResp2.LogID())

	// Example 3: Update the Claude bot to disable thinking
	fmt.Println("=== Updating Claude bot to disable thinking ===")
	updateResp, err := cozeCli.Bots.Update(ctx, &coze.UpdateBotsReq{
		BotID:       claudeBotID,
		Name:        "Claude Bot (Thinking Disabled)",
		Description: "Claude bot with thinking disabled",
		ModelInfoConfig: &coze.BotModelInfoConfig{
			ModelID:        "claude-3-5-sonnet",
			ResponseFormat: coze.ResponseFormatText,
			Parameters: map[string]string{
				"thinking_type": "disabled",
			},
		},
	})
	if err != nil {
		fmt.Println("Error updating Claude bot:", err)
		return
	}
	fmt.Printf("Updated Claude bot successfully\n")
	fmt.Printf("Log ID: %s\n\n", updateResp.LogID())

	// Example 4: Create a Gemini bot with thinking enabled
	fmt.Println("=== Creating bot with Gemini thinking enabled ===")
	createResp3, err := cozeCli.Bots.Create(ctx, &coze.CreateBotsReq{
		SpaceID:     workspaceID,
		Name:        "Gemini Thinking Bot",
		Description: "A bot with thinking capabilities using Gemini model",
		PromptInfo: &coze.BotPromptInfo{
			Prompt: "You are a helpful assistant powered by Gemini with thinking capabilities.",
		},
		ModelInfoConfig: &coze.BotModelInfoConfig{
			ModelID:        "gemini-1.5-pro",
			ResponseFormat: coze.ResponseFormatMarkdown,
			Temperature:    0.8,
			Parameters: map[string]string{
				"thinking_type": "enable",
			},
		},
	})
	if err != nil {
		fmt.Println("Error creating Gemini bot:", err)
		return
	}
	geminiBotID := createResp3.BotID
	fmt.Printf("Created Gemini bot with ID: %s\n", geminiBotID)
	fmt.Printf("Log ID: %s\n\n", createResp3.LogID())

	fmt.Println("=== Summary ===")
	fmt.Printf("Claude Bot ID: %s (thinking disabled)\n", claudeBotID)
	fmt.Printf("Doubao Bot ID: %s (auto thinking)\n", doubaoBotID)
	fmt.Printf("Gemini Bot ID: %s (thinking enabled)\n", geminiBotID)
}