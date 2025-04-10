package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

// This examples is for describing how to create a bot, update a bot and publish a bot to the API.
func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	// step one, create a bot
	workspaceID := os.Getenv("WORKSPACE_ID")

	// Call the upload file interface to get the avatar id.
	avatarPath := os.Getenv("IMAGE_FILE_PATH")
	ctx := context.Background()
	file, err := os.Open(avatarPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	uploadReq := &coze.UploadFilesReq{
		File: file,
	}
	avatarInfo, err := cozeCli.Files.Upload(ctx, uploadReq)
	if err != nil {
		fmt.Println("Error uploading avatar:", err)
		return
	}
	fmt.Println(avatarInfo)

	// build the request
	createResp, err := cozeCli.Bots.Create(ctx, &coze.CreateBotsReq{
		SpaceID:     workspaceID,
		Description: "the description of your bot",
		Name:        "the name of your bot",
		PromptInfo: &coze.BotPromptInfo{
			Prompt: "your prompt",
		},
		OnboardingInfo: &coze.BotOnboardingInfo{
			Prologue:           "the prologue of your bot",
			SuggestedQuestions: []string{"question 1", "question 2"},
		},
		IconFileID: avatarInfo.FileInfo.ID,
		ModelInfoConfig: &coze.BotModelInfoConfig{
			ModelID: "1738675210",
		},
		WorkflowIDList: &coze.WorkflowIDList{
			IDs: []coze.WorkflowIDInfo{
				{
					ID: os.Getenv("WORKFLOW_ID"), // 工作流ID
				},
			},
		},
	})
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}
	botID := createResp.BotID
	fmt.Println(createResp)
	fmt.Println(createResp.LogID())

	//
	// step two, update the bot, you can update the bot after being created
	// in this examples, we will update the avatar of the bot

	publishResp, err := cozeCli.Bots.Publish(ctx, &coze.PublishBotsReq{
		BotID:        botID,
		ConnectorIDs: []string{"1024"},
	})
	if err != nil {
		fmt.Println("Error publishing bot:", err)
		return
	}
	fmt.Println(publishResp)
	fmt.Println(publishResp.LogID())

	//
	// step three, you can also modify the bot configuration and republish it.
	// in this examples, we will update the avatar of the bot

	newFile, err := os.Open(os.Getenv("IMAGE_FILE_PATH"))
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	newUploadReq := &coze.UploadFilesReq{
		File: newFile,
	}
	newAvatarInfo, err := cozeCli.Files.Upload(ctx, newUploadReq)
	if err != nil {
		fmt.Println("Error uploading new avatar:", err)
		return
	}
	fmt.Println(newAvatarInfo)
	fmt.Println(newAvatarInfo.LogID())

	// Update bot
	updateResp, err := cozeCli.Bots.Update(ctx, &coze.UpdateBotsReq{
		BotID:      botID,
		Name:       "the name of your bot",
		IconFileID: newAvatarInfo.FileInfo.ID,
		ModelInfoConfig: &coze.BotModelInfoConfig{
			ModelID: "1738675210",
		},
		WorkflowIDList: &coze.WorkflowIDList{
			IDs: []coze.WorkflowIDInfo{
				{
					ID: os.Getenv("WORKFLOW_ID"),
				},
			},
		},
	})
	if err != nil {
		fmt.Println("Error updating bot:", err)
		return
	}
	fmt.Println(updateResp.LogID())

	// Republish bot
	publishResp, err = cozeCli.Bots.Publish(ctx, &coze.PublishBotsReq{
		BotID:        botID,
		ConnectorIDs: []string{"1024"},
	})
	if err != nil {
		fmt.Println("Error republishing bot:", err)
		return
	}
	fmt.Println(publishResp)
	fmt.Println(publishResp.LogID())
}
