package coze

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatasets(t *testing.T) {
	t.Run("Create dataset success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/datasets", req.URL.Path)

				// Return mock response
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"data": {
							"dataset_id": "123"
						}
					}`)),
					Header: make(http.Header),
				}, nil
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		datasets := newDatasets(core)

		// Create test request
		req := &CreateDatasetsReq{
			Name:        "test_dataset",
			SpaceID:     "space_123",
			FormatType:  DocumentFormatTypeDocument,
			Description: "Test dataset description",
			IconFileID:  "icon_123",
		}

		// Test dataset creation
		resp, err := datasets.Create(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "123", resp.DatasetID)
	})

	t.Run("List datasets success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/datasets", req.URL.Path)

				// Return mock response
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"data": {
							"total_count": 2,
							"dataset_list": [
								{
									"dataset_id": "123",
									"name": "dataset1",
									"space_id": "space_123",
									"status": 1,
									"format_type": 0
								},
								{
									"dataset_id": "456",
									"name": "dataset2",
									"space_id": "space_123",
									"status": 1,
									"format_type": 0
								}
							]
						}
					}`)),
					Header: make(http.Header),
				}, nil
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		datasets := newDatasets(core)

		// Create test request
		req := NewListDatasetsReq("space_123")
		req.Name = "dataset"
		req.FormatType = DocumentFormatTypeDocument

		// Test dataset listing
		pager, err := datasets.List(context.Background(), req)
		require.NoError(t, err)

		// Verify pagination results
		items := pager.Items()
		assert.Len(t, items, 2)
		assert.Equal(t, "123", items[0].ID)
		assert.Equal(t, "456", items[1].ID)
		assert.Equal(t, int(2), pager.Total())
		assert.False(t, pager.HasMore())
	})

	t.Run("Update dataset success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPut, req.Method)
				assert.Equal(t, "/v1/datasets/123", req.URL.Path)

				// Return mock response
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"data": {}}`)),
					Header:     make(http.Header),
				}, nil
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		datasets := newDatasets(core)

		// Create test request
		req := &UpdateDatasetsReq{
			DatasetID:   "123",
			Name:        "updated_dataset",
			Description: "Updated description",
			IconFileID:  "new_icon_123",
		}

		// Test dataset update
		resp, err := datasets.Update(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("Delete dataset success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, "/v1/datasets/123", req.URL.Path)

				// Return mock response
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"data": {}}`)),
					Header:     make(http.Header),
				}, nil
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		datasets := newDatasets(core)

		// Create test request
		req := &DeleteDatasetsReq{
			DatasetID: "123",
		}

		// Test dataset deletion
		resp, err := datasets.Delete(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("Process documents success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/datasets/123/process", req.URL.Path)

				// Return mock response
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"data": {
							"data": [
								{
									"document_id": "doc_123",
									"status": 1,
									"progress": 100,
									"document_name": "test.txt"
								}
							]
						}
					}`)),
					Header: make(http.Header),
				}, nil
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		datasets := newDatasets(core)

		// Create test request
		req := &ProcessDocumentsReq{
			DatasetID:   "123",
			DocumentIDs: []string{"doc_123"},
		}

		// Test document processing
		resp, err := datasets.Process(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Data, 1)

		progress := resp.Data[0]
		assert.Equal(t, "doc_123", progress.DocumentID)
		assert.Equal(t, DocumentStatusCompleted, progress.Status)
		assert.Equal(t, 100, progress.Progress)
		assert.Equal(t, "test.txt", progress.DocumentName)
	})
}

// Test dataset status constants
func TestDatasetStatus(t *testing.T) {
	t.Run("DatasetStatus constants", func(t *testing.T) {
		assert.Equal(t, DatasetStatus(1), DatasetStatusEnabled)
		assert.Equal(t, DatasetStatus(3), DatasetStatusDisabled)
	})
}
