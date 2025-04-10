package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

// This examples demonstrates how to handle different types of errors from the Coze API.
func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()

	// Example 1: Handle API error
	_, err := cozeCli.Bots.Retrieve(ctx, &coze.RetrieveBotsReq{
		BotID: "invalid_bot_id",
	})
	if err != nil {
		if cozeErr, ok := coze.AsCozeError(err); ok {
			// Handle Coze API error
			fmt.Printf("Coze API error: %s (code: %d)\n", cozeErr.Message, cozeErr.Code)
			return
		}
		// Handle other errors
		fmt.Printf("Other error: %v\n", err)
		return
	}

	// Example 2: Handle auth error
	invalidToken := "invalid_token"
	invalidAuthCli := coze.NewTokenAuth(invalidToken)
	invalidCozeCli := coze.NewCozeAPI(invalidAuthCli)

	_, err = invalidCozeCli.Bots.List(ctx, &coze.ListBotsReq{
		PageNum:  1,
		PageSize: 10,
	})
	if err != nil {
		if cozeErr, ok := coze.AsAuthError(err); ok {
			// Handle auth error
			if cozeErr.Code == "unauthorized" {
				fmt.Println("Authentication failed. Please check your token.")
				return
			}
			fmt.Printf("Coze API error: %s (code: %s)\n", cozeErr.ErrorMessage, cozeErr.Code)
			return
		}
		fmt.Printf("Other error: %v\n", err)
		return
	}
}
