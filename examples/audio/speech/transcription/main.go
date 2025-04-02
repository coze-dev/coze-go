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

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	filename := "/Users/u/Downloads/hello.mp3"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	resp, err := cozeCli.Audio.Speech.Transcription(context.Background(), file, filename)
	if err != nil {
		fmt.Println("Error creating speech:", err)
		return
	}

	fmt.Println(resp)
	fmt.Println(resp.LogID())
}
