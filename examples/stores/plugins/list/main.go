package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/coze-dev/coze-go"
)

func main() {
	cozeAPIToken := os.Getenv("COZE_API_TOKEN")
	cozeAPIBase := os.Getenv("COZE_API_BASE")
	if cozeAPIBase == "" {
		cozeAPIBase = coze.CnBaseURL
	}

	// Init the Coze client through the access_token.
	authCli := coze.NewTokenAuth(cozeAPIToken)
	client := coze.NewCozeAPI(authCli, coze.WithBaseURL(cozeAPIBase))
	ctx := context.Background()

	listResp, err := client.Stores.Plugins.List(ctx, &coze.ListStoresPluginsReq{
		PageSize: 20,
	})
	if err != nil {
		log.Printf("Error listing plugins: %v", err)
		return
	}
	for listResp.Next() {
		plugin := listResp.Current()
		fmt.Printf("Plugin: %s %s\n", plugin.Metainfo.EntityID, plugin.Metainfo.Name)
	}
	err = listResp.Err()
	if err != nil {
		log.Fatalf("List plugins failed: %v", err)
		return
	}
}
