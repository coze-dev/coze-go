package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/coze-dev/coze-go"
)

func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()
	filePath := os.Getenv("IMAGE_FILE_PATH")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// the developers can directly use file to upload
	uploadResp, err := cozeCli.Files.Upload(ctx, &coze.UploadFilesReq{
		File: file,
	})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}
	fileInfo := uploadResp.FileInfo
	// The developers can also upload file by bytes stream
	// uploadResp, err = cozeCli.Files.Upload(ctx, coze.NewUploadFileReq(bytes.NewReader([]byte("hello world"))," your file name"))
	// if err != nil{
	//	 fmt.Println("Error uploading file:", err)
	//	 return
	// }

	// wait the server to process the file
	time.Sleep(time.Second)

	// retrieve file
	retrievedResp, err := cozeCli.Files.Retrieve(ctx, &coze.RetrieveFilesReq{
		FileID: fileInfo.ID,
	})
	if err != nil {
		fmt.Println("Error retrieving file:", err)
		return
	}
	retrievedInfo := retrievedResp.FileInfo
	fmt.Println(retrievedInfo)
	fmt.Println(retrievedResp.LogID())
}
