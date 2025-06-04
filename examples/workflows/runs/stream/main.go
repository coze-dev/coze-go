package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

// This examples describes how to use the workflow interface to stream chats.
func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	customClient := &http.Client{
		Timeout: time.Minute * 20,
	}

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli,
		coze.WithBaseURL(os.Getenv("COZE_API_BASE")),
		coze.WithHttpClient(customClient),
		coze.WithLogLevel(coze.LogLevelDebug),
	)

	ctx := context.Background()
	workflowID := os.Getenv("COZE_WORKFLOW_ID")
	appID := os.Getenv("COZE_APP_ID")

	// if your workflow need input params, you can send them by map
	data := map[string]interface{}{
		"date": "param values",
	}

	req := &coze.RunWorkflowsReq{
		WorkflowID: workflowID,
		Parameters: data,
		AppID:      appID,
	}

	resp, err := cozeCli.Workflows.Runs.Stream(ctx, req)
	if err != nil {
		fmt.Println("Error starting stream:", err)
		return
	}

	handleEvents(ctx, resp, cozeCli, workflowID)
}

// The stream interface will return an iterator of WorkflowEvent. Developers should iterate
// through this iterator to obtain WorkflowEvent and handle them separately according to
// the type of WorkflowEvent.
func handleEvents(ctx context.Context, resp coze.Stream[coze.WorkflowEvent], cozeCli coze.CozeAPI, workflowID string) {
	defer resp.Close()
	for {
		event, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished")
			break
		}
		if err != nil {
			fmt.Println("Error receiving event:", err)
			break
		}

		switch event.Event {
		case coze.WorkflowEventTypeMessage:
			fmt.Println("Got message:", event.Message)
		case coze.WorkflowEventTypeError:
			fmt.Println("Got error:", event.Error)
		case coze.WorkflowEventTypeDone:
			fmt.Println("Got message:", event.Message)
		case coze.WorkflowEventTypeInterrupt:
			resumeReq := &coze.ResumeRunWorkflowsReq{
				WorkflowID:    workflowID,
				EventID:       event.Interrupt.InterruptData.EventID,
				ResumeData:    "your data",
				InterruptType: event.Interrupt.InterruptData.Type,
			}
			newResp, err := cozeCli.Workflows.Runs.Resume(ctx, resumeReq)
			if err != nil {
				fmt.Println("Error resuming workflow:", err)
				return
			}
			fmt.Println("start resume workflow")
			handleEvents(ctx, newResp, cozeCli, workflowID)
		}
	}
	fmt.Printf("done, log:%s\n", resp.Response().LogID())
}
