package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataset(t *testing.T) {
	as := assert.New(t)
	t.Run("create dataset success", func(t *testing.T) {
		datasetID := randomString(10)
		datasets := newDatasets(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/datasets", req.URL.Path)
			return mockResponse(http.StatusOK, &createDatasetResp{
				Data: &CreateDatasetResp{
					DatasetID: datasetID,
				},
			})
		})))
		req := &CreateDatasetsReq{
			Name:        "test_dataset",
			SpaceID:     "space_123",
			FormatType:  DocumentFormatTypeDocument,
			Description: "Test dataset description",
			IconFileID:  "icon_123",
		}
		resp, err := datasets.Create(context.Background(), req)
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal(datasetID, resp.DatasetID)
	})

	t.Run("list datasets success", func(t *testing.T) {
		datasets := newDatasets(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/datasets", req.URL.Path)
			return mockResponse(http.StatusOK, &listDatasetsResp{
				Data: &ListDatasetsResp{
					TotalCount: 2,
					DatasetList: []*Dataset{
						{
							ID:         "123",
							Name:       "dataset1",
							SpaceID:    "space_123",
							Status:     DatasetStatusEnabled,
							FormatType: DocumentFormatTypeDocument,
						},
						{
							ID:         "456",
							Name:       "dataset2",
							SpaceID:    "space_123",
							Status:     DatasetStatusEnabled,
							FormatType: DocumentFormatTypeDocument,
						},
					},
				},
			})
		})))
		req := NewListDatasetsReq("space_123")
		req.Name = "dataset"
		req.PageNum = 0  // 提高覆盖率
		req.PageSize = 0 // 提高覆盖率
		req.FormatType = DocumentFormatTypeDocument

		paged, err := datasets.List(context.Background(), req)
		as.Nil(err)
		as.NotNil(paged)
		// as.NotEmpty(paged.resp) // todo

		items := paged.Items()
		as.Len(items, 2)
		as.Equal("123", items[0].ID)
		as.Equal("456", items[1].ID)
		as.Equal(int(2), paged.Total())
		as.False(paged.HasMore())
	})

	t.Run("update dataset success", func(t *testing.T) {
		datasets := newDatasets(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPut, req.Method)
			as.Equal("/v1/datasets/123", req.URL.Path)
			return mockResponse(http.StatusOK, &updateDatasetResp{
				Data: &UpdateDatasetsResp{},
			})
		})))
		req := &UpdateDatasetsReq{
			DatasetID:   "123",
			Name:        "updated_dataset",
			Description: "Updated description",
			IconFileID:  "new_icon_123",
		}
		resp, err := datasets.Update(context.Background(), req)
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
	})

	t.Run("Delete dataset success", func(t *testing.T) {
		datasets := newDatasets(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodDelete, req.Method)
			as.Equal("/v1/datasets/123", req.URL.Path)
			return mockResponse(http.StatusOK, &updateDatasetResp{
				Data: &UpdateDatasetsResp{},
			})
		})))
		req := &DeleteDatasetsReq{
			DatasetID: "123",
		}
		resp, err := datasets.Delete(context.Background(), req)
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
	})

	t.Run("Process documents success", func(t *testing.T) {
		datasets := newDatasets(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/datasets/123/process", req.URL.Path)
			return mockResponse(http.StatusOK, &processDocumentsResp{
				Data: &ProcessDocumentsResp{
					Data: []*DocumentProgress{
						{
							DocumentID:   "doc_123",
							Status:       DocumentStatusCompleted,
							Progress:     100,
							DocumentName: "test.txt",
						},
					},
				},
			})
		})))
		req := &ProcessDocumentsReq{
			DatasetID:   "123",
			DocumentIDs: []string{"doc_123"},
		}
		resp, err := datasets.Process(context.Background(), req)
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Len(resp.Data, 1)

		progress := resp.Data[0]
		as.Equal("doc_123", progress.DocumentID)
		as.Equal(DocumentStatusCompleted, progress.Status)
		as.Equal(100, progress.Progress)
		as.Equal("test.txt", progress.DocumentName)
	})
}

func TestDatasetStatus(t *testing.T) {
	as := assert.New(t)
	t.Run("dataset status constants", func(t *testing.T) {
		as.Equal(DatasetStatus(1), DatasetStatusEnabled)
		as.Equal(DatasetStatus(3), DatasetStatusDisabled)
	})
}
