package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatasetsDocuments(t *testing.T) {
	// Test Create method
	t.Run("Create document success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/open_api/knowledge/document/create", req.URL.Path)

				// Verify headers
				assert.Equal(t, "str", req.Header.Get("Agw-Js-Conv"))

				// Return mock response
				return mockResponse(http.StatusOK, &createDatasetsDocumentsResp{
					CreateDatasetsDocumentsResp: &CreateDatasetsDocumentsResp{
						DocumentInfos: []*Document{
							{
								DocumentID: "doc1",
								Name:       "test.txt",
								CharCount:  100,
								Size:       1024,
								Type:       "txt",
								Status:     DocumentStatusCompleted,
								FormatType: DocumentFormatTypeDocument,
								SourceType: DocumentSourceTypeLocalFile,
								SliceCount: 1,
								CreateTime: 1234567890,
								UpdateTime: 1234567890,
								ChunkStrategy: &DocumentChunkStrategy{
									ChunkType: 0,
								},
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		documents := newDatasetsDocuments(core)

		resp, err := documents.Create(context.Background(), &CreateDatasetsDocumentsReq{
			DatasetID: 123,
			DocumentBases: []*DocumentBase{
				DocumentBaseBuildLocalFile("test.txt", "test content", "txt"),
			},
			ChunkStrategy: &DocumentChunkStrategy{
				ChunkType: 0,
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		require.Len(t, resp.DocumentInfos, 1)
		doc := resp.DocumentInfos[0]
		assert.Equal(t, "doc1", doc.DocumentID)
		assert.Equal(t, "test.txt", doc.Name)
		assert.Equal(t, DocumentStatusCompleted, doc.Status)
		assert.Equal(t, DocumentFormatTypeDocument, doc.FormatType)
		assert.Equal(t, DocumentSourceTypeLocalFile, doc.SourceType)
	})

	// Test Update method
	t.Run("Update document success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/open_api/knowledge/document/update", req.URL.Path)

				// Verify headers
				assert.Equal(t, "str", req.Header.Get("Agw-Js-Conv"))

				// Return mock response
				return mockResponse(http.StatusOK, &updateDatasetsDocumentsResp{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		documents := newDatasetsDocuments(core)

		resp, err := documents.Update(context.Background(), &UpdateDatasetsDocumentsReq{
			DocumentID:   123,
			DocumentName: "updated.txt",
			UpdateRule:   DocumentUpdateRuleBuildAutoUpdate(24),
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
	})

	// Test Delete method
	t.Run("Delete document success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/open_api/knowledge/document/delete", req.URL.Path)

				// Verify headers
				assert.Equal(t, "str", req.Header.Get("Agw-Js-Conv"))

				// Return mock response
				return mockResponse(http.StatusOK, &deleteDatasetsDocumentsResp{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		documents := newDatasetsDocuments(core)

		resp, err := documents.Delete(context.Background(), &DeleteDatasetsDocumentsReq{
			DocumentIDs: []int64{123, 456},
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
	})

	// Test List method
	t.Run("List documents success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/open_api/knowledge/document/list", req.URL.Path)

				// Verify headers
				assert.Equal(t, "str", req.Header.Get("Agw-Js-Conv"))

				// Return mock response
				return mockResponse(http.StatusOK, &listDatasetsDocumentsResp{
					ListDatasetsDocumentsResp: &ListDatasetsDocumentsResp{
						Total: 2,
						DocumentInfos: []*Document{
							{
								DocumentID: "doc1",
								Name:       "test1.txt",
								Status:     DocumentStatusCompleted,
								FormatType: DocumentFormatTypeDocument,
								SourceType: DocumentSourceTypeLocalFile,
								CreateTime: 1234567890,
								UpdateTime: 1234567890,
							},
							{
								DocumentID: "doc2",
								Name:       "test2.txt",
								Status:     DocumentStatusCompleted,
								FormatType: DocumentFormatTypeDocument,
								SourceType: DocumentSourceTypeLocalFile,
								CreateTime: 1234567891,
								UpdateTime: 1234567891,
							},
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		documents := newDatasetsDocuments(core)

		paged, err := documents.List(context.Background(), &ListDatasetsDocumentsReq{
			DatasetID: 123,
			Page:      1,
			Size:      20,
		})

		require.NoError(t, err)
		items := paged.Items()
		require.Len(t, items, 2)

		// Verify first document
		assert.Equal(t, "doc1", items[0].DocumentID)
		assert.Equal(t, "test1.txt", items[0].Name)
		assert.Equal(t, DocumentStatusCompleted, items[0].Status)
		assert.Equal(t, DocumentFormatTypeDocument, items[0].FormatType)
		assert.Equal(t, DocumentSourceTypeLocalFile, items[0].SourceType)

		// Verify second document
		assert.Equal(t, "doc2", items[1].DocumentID)
		assert.Equal(t, "test2.txt", items[1].Name)
		assert.Equal(t, DocumentStatusCompleted, items[1].Status)
		assert.Equal(t, DocumentFormatTypeDocument, items[1].FormatType)
		assert.Equal(t, DocumentSourceTypeLocalFile, items[1].SourceType)
	})

	// Test List method with default pagination
	t.Run("List documents with default pagination", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return mock response
				return mockResponse(http.StatusOK, &listDatasetsDocumentsResp{
					ListDatasetsDocumentsResp: &ListDatasetsDocumentsResp{
						Total:         0,
						DocumentInfos: []*Document{},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		documents := newDatasetsDocuments(core)

		paged, err := documents.List(context.Background(), &ListDatasetsDocumentsReq{
			DatasetID: 123,
		})

		require.NoError(t, err)
		assert.Empty(t, paged.Items())
	})

	// Test helper functions
	t.Run("Test helper functions", func(t *testing.T) {
		// Test BuildWebPage
		webPage := DocumentBaseBuildWebPage("test page", "https://example.com", nil)
		assert.Equal(t, "test page", webPage.Name)
		assert.Equal(t, "https://example.com", *webPage.SourceInfo.WebUrl)
		assert.Equal(t, 1, *webPage.SourceInfo.DocumentSource)
		assert.Equal(t, DocumentUpdateTypeNoAutoUpdate, webPage.UpdateRule.UpdateType)

		// Test BuildWebPageWithInterval
		webPageWithInterval := DocumentBaseBuildWebPage("test page", "https://example.com", ptr(24))
		assert.Equal(t, "test page", webPageWithInterval.Name)
		assert.Equal(t, "https://example.com", *webPageWithInterval.SourceInfo.WebUrl)
		assert.Equal(t, 1, *webPageWithInterval.SourceInfo.DocumentSource)
		assert.Equal(t, DocumentUpdateTypeAutoUpdate, webPageWithInterval.UpdateRule.UpdateType)
		assert.Equal(t, 24, webPageWithInterval.UpdateRule.UpdateInterval)

		// Test BuildLocalFile
		localFile := DocumentBaseBuildLocalFile("test.txt", "test content", "txt")
		assert.Equal(t, "test.txt", localFile.Name)
		assert.Equal(t, "txt", *localFile.SourceInfo.FileType)
		assert.NotEmpty(t, localFile.SourceInfo.FileBase64)

		// Test BuildAutoUpdateRule
		autoUpdateRule := DocumentUpdateRuleBuildAutoUpdate(24)
		assert.Equal(t, DocumentUpdateTypeAutoUpdate, autoUpdateRule.UpdateType)
		assert.Equal(t, 24, autoUpdateRule.UpdateInterval)

		// Test BuildNoAutoUpdateRule
		noAutoUpdateRule := DocumentUpdateRuleBuildNoAuto()
		assert.Equal(t, DocumentUpdateTypeNoAutoUpdate, noAutoUpdateRule.UpdateType)
	})
}

// Test document type constants
func TestDocumentTypes(t *testing.T) {
	t.Run("DocumentFormatType constants", func(t *testing.T) {
		assert.Equal(t, DocumentFormatType(0), DocumentFormatTypeDocument)
		assert.Equal(t, DocumentFormatType(1), DocumentFormatTypeSpreadsheet)
		assert.Equal(t, DocumentFormatType(2), DocumentFormatTypeImage)
	})

	t.Run("DocumentSourceType constants", func(t *testing.T) {
		assert.Equal(t, DocumentSourceType(0), DocumentSourceTypeLocalFile)
		assert.Equal(t, DocumentSourceType(1), DocumentSourceTypeOnlineWeb)
	})

	t.Run("DocumentStatus constants", func(t *testing.T) {
		assert.Equal(t, DocumentStatus(0), DocumentStatusProcessing)
		assert.Equal(t, DocumentStatus(1), DocumentStatusCompleted)
		assert.Equal(t, DocumentStatus(9), DocumentStatusFailed)
	})

	t.Run("DocumentUpdateType constants", func(t *testing.T) {
		assert.Equal(t, DocumentUpdateType(0), DocumentUpdateTypeNoAutoUpdate)
		assert.Equal(t, DocumentUpdateType(1), DocumentUpdateTypeAutoUpdate)
	})
}
