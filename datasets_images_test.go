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

func TestDatasetsImages(t *testing.T) {
	t.Run("Update image success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPut, req.Method)
				assert.Equal(t, "/v1/datasets/123/images/456", req.URL.Path)

				// Return mock response
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"data": {}}`)),
					Header:     make(http.Header),
				}, nil
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		images := newDatasetsImages(core)

		// Create test request
		caption := "test caption"
		req := &UpdateDatasetImageReq{
			DatasetID:  "123",
			DocumentID: "456",
			Caption:    &caption,
		}

		// Test image update
		resp, err := images.Update(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("List images success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/datasets/123/images", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "test", req.URL.Query().Get("keyword"))
				assert.Equal(t, "true", req.URL.Query().Get("has_caption"))
				assert.Equal(t, "1", req.URL.Query().Get("page_num"))
				assert.Equal(t, "10", req.URL.Query().Get("page_size"))

				// Return mock response
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"data": {
							"total_count": 2,
							"photo_infos": [
								{
									"document_id": "img1",
									"name": "image1.png",
									"status": 1,
									"format_type": 2,
									"source_type": 0,
									"caption": "test image 1"
								},
								{
									"document_id": "img2",
									"name": "image2.png",
									"status": 1,
									"format_type": 2,
									"source_type": 0,
									"caption": "test image 2"
								}
							]
						}
					}`)),
					Header: make(http.Header),
				}, nil
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		images := newDatasetsImages(core)

		// Create test request
		keyword := "test"
		hasCaption := true
		req := &ListDatasetsImagesReq{
			DatasetID:  "123",
			Keyword:    &keyword,
			HasCaption: &hasCaption,
			PageNum:    1,
			PageSize:   10,
		}

		// Test image listing
		pager, err := images.List(context.Background(), req)
		require.NoError(t, err)

		// Verify pagination results
		items := pager.Items()
		require.Len(t, items, 2)
		assert.Equal(t, "img1", items[0].DocumentID)
		assert.Equal(t, "image1.png", items[0].Name)
		assert.Equal(t, "test image 1", items[0].Caption)
		assert.Equal(t, ImageStatusCompleted, items[0].Status)
		assert.Equal(t, DocumentFormatTypeImage, items[0].FormatType)
		assert.Equal(t, DocumentSourceTypeLocalFile, items[0].SourceType)

		assert.Equal(t, "img2", items[1].DocumentID)
		assert.Equal(t, "image2.png", items[1].Name)
		assert.Equal(t, "test image 2", items[1].Caption)
		assert.Equal(t, ImageStatusCompleted, items[1].Status)
		assert.Equal(t, DocumentFormatTypeImage, items[1].FormatType)
		assert.Equal(t, DocumentSourceTypeLocalFile, items[1].SourceType)

		assert.Equal(t, 2, pager.Total())
		assert.False(t, pager.HasMore())
	})
}
