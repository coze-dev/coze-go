package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/coze-dev/coze-go"
)

func main() {
	// Get an access_token through personal access token or oauth.
	token := os.Getenv("COZE_API_TOKEN")
	authCli := coze.NewTokenAuth(token)

	// Init the Coze client through the access_token.
	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(os.Getenv("COZE_API_BASE")))

	ctx := context.Background()
	datasetID, _ := strconv.ParseInt(os.Getenv("DATASET_ID"), 10, 64)

	//
	// create document in to specific dataset

	createReq := &coze.CreateDatasetsDocumentsReq{
		DatasetID: datasetID,
		DocumentBases: []*coze.DocumentBase{
			coze.DocumentBaseBuildLocalFile("file doc examples", "your file content", "txt"),
		},
	}
	createResp, err := cozeCli.Datasets.Documents.Create(ctx, createReq)
	if err != nil {
		fmt.Println("Error creating documents:", err)
		return
	}
	fmt.Println(createResp)
	fmt.Println(createResp.LogID())

	var documentIDs []int64
	for _, doc := range createResp.DocumentInfos {
		id, _ := strconv.ParseInt(doc.DocumentID, 10, 64)
		documentIDs = append(documentIDs, id)
	}

	//
	// update document. It means success that no exception has been thrown

	updateReq := &coze.UpdateDatasetsDocumentsReq{
		DocumentID:   documentIDs[0],
		DocumentName: "new name",
	}
	updateResp, err := cozeCli.Datasets.Documents.Update(ctx, updateReq)
	if err != nil {
		fmt.Println("Error updating document:", err)
		return
	}
	fmt.Println(updateResp)
	fmt.Println(updateResp.LogID())

	//
	// delete document. It means success that no exception has been thrown

	deleteResp, err := cozeCli.Datasets.Documents.Delete(ctx, &coze.DeleteDatasetsDocumentsReq{
		DocumentIDs: []int64{documentIDs[0]},
	})
	if err != nil {
		fmt.Println("Error deleting document:", err)
		return
	}
	fmt.Println(deleteResp)
	fmt.Println(deleteResp.LogID())
}
