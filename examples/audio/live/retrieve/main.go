package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coze-dev/coze-go"
)

func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	// Init the coze client.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()
	liveID := os.Getenv("LIVE_ID")

	// Retrieve live information
	liveInfo, err := cozeCli.Audio.Live.Retrieve(ctx, &coze.RetrieveAudioLiveReq{
		LiveID: liveID,
	})
	if err != nil {
		fmt.Printf("Error retrieving live info: %v\n", err)
		return
	}

	fmt.Printf("Live Info:\n")
	fmt.Printf("App ID: %s\n", liveInfo.AppID)
	fmt.Printf("Number of streams: %d\n", len(liveInfo.StreamInfos))

	for i, stream := range liveInfo.StreamInfos {
		fmt.Printf("Stream %d:\n", i+1)
		fmt.Printf("  Stream ID: %s\n", stream.StreamID)
		fmt.Printf("  Name: %s\n", stream.Name)
		fmt.Printf("  Type: %s\n", stream.LiveType)
	}
}
