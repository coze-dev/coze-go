package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatasetsDocuments(t *testing.T) {
	as := assert.New(t)
	t.Run("Create document success", func(t *testing.T) {
		documents := newDatasetsDocuments(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/open_api/knowledge/document/create", req.URL.Path)
			as.Equal("str", req.Header.Get("Agw-Js-Conv"))
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
		})))
		resp, err := documents.Create(context.Background(), &CreateDatasetsDocumentsReq{
			DatasetID: 123,
			DocumentBases: []*DocumentBase{
				DocumentBaseBuildLocalFile("test.txt", "test content", "txt"),
			},
			ChunkStrategy: &DocumentChunkStrategy{
				ChunkType: 0,
			},
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Len(resp.DocumentInfos, 1)
		doc := resp.DocumentInfos[0]
		as.Equal("doc1", doc.DocumentID)
		as.Equal("test.txt", doc.Name)
		as.Equal(DocumentStatusCompleted, doc.Status)
		as.Equal(DocumentFormatTypeDocument, doc.FormatType)
		as.Equal(DocumentSourceTypeLocalFile, doc.SourceType)
	})

	t.Run("update document success", func(t *testing.T) {
		documents := newDatasetsDocuments(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/open_api/knowledge/document/update", req.URL.Path)
			as.Equal("str", req.Header.Get("Agw-Js-Conv"))
			return mockResponse(http.StatusOK, &updateDatasetsDocumentsResp{})
		})))
		resp, err := documents.Update(context.Background(), &UpdateDatasetsDocumentsReq{
			DocumentID:   123,
			DocumentName: "updated.txt",
			UpdateRule:   DocumentUpdateRuleBuildAutoUpdate(24),
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
	})

	t.Run("delete document success", func(t *testing.T) {
		documents := newDatasetsDocuments(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/open_api/knowledge/document/delete", req.URL.Path)
			as.Equal("str", req.Header.Get("Agw-Js-Conv"))
			return mockResponse(http.StatusOK, &deleteDatasetsDocumentsResp{})
		})))
		resp, err := documents.Delete(context.Background(), &DeleteDatasetsDocumentsReq{
			DocumentIDs: []int64{123, 456},
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
	})

	t.Run("list documents success", func(t *testing.T) {
		documents := newDatasetsDocuments(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/open_api/knowledge/document/list", req.URL.Path)
			as.Equal("str", req.Header.Get("Agw-Js-Conv"))
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
		})))
		paged, err := documents.List(context.Background(), &ListDatasetsDocumentsReq{
			DatasetID: 123,
			Page:      1,
			Size:      20,
		})
		as.Nil(err)
		as.NotNil(paged)
		// as.NotEmpty(paged.Response().LogID()) // todo
		items := paged.Items()
		as.Len(items, 2)

		as.Equal("doc1", items[0].DocumentID)
		as.Equal("test1.txt", items[0].Name)
		as.Equal(DocumentStatusCompleted, items[0].Status)
		as.Equal(DocumentFormatTypeDocument, items[0].FormatType)
		as.Equal(DocumentSourceTypeLocalFile, items[0].SourceType)

		as.Equal("doc2", items[1].DocumentID)
		as.Equal("test2.txt", items[1].Name)
		as.Equal(DocumentStatusCompleted, items[1].Status)
		as.Equal(DocumentFormatTypeDocument, items[1].FormatType)
		as.Equal(DocumentSourceTypeLocalFile, items[1].SourceType)
	})

	// Test List method with default pagination
	t.Run("List documents with default pagination", func(t *testing.T) {
		documents := newDatasetsDocuments(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return mockResponse(http.StatusOK, &listDatasetsDocumentsResp{
				ListDatasetsDocumentsResp: &ListDatasetsDocumentsResp{
					Total:         0,
					DocumentInfos: []*Document{},
				},
			})
		})))
		paged, err := documents.List(context.Background(), &ListDatasetsDocumentsReq{
			DatasetID: 123,
		})
		as.Nil(err)
		as.NotNil(paged)
		// as.NotEmpty(paged.Response().LogID()) // todo
		as.Empty(paged.Items())
	})

	t.Run("helper functions", func(t *testing.T) {
		webPage := DocumentBaseBuildWebPage("test page", "https://example.com", nil)
		as.Equal("test page", webPage.Name)
		as.Equal("https://example.com", *webPage.SourceInfo.WebUrl)
		as.Equal(1, *webPage.SourceInfo.DocumentSource)
		as.Equal(DocumentUpdateTypeNoAutoUpdate, webPage.UpdateRule.UpdateType)

		// Test BuildWebPageWithInterval
		webPageWithInterval := DocumentBaseBuildWebPage("test page", "https://example.com", ptr(24))
		as.Equal("test page", webPageWithInterval.Name)
		as.Equal("https://example.com", *webPageWithInterval.SourceInfo.WebUrl)
		as.Equal(1, *webPageWithInterval.SourceInfo.DocumentSource)
		as.Equal(DocumentUpdateTypeAutoUpdate, webPageWithInterval.UpdateRule.UpdateType)
		as.Equal(24, webPageWithInterval.UpdateRule.UpdateInterval)

		// Test BuildLocalFile
		localFile := DocumentBaseBuildLocalFile("test.txt", "test content", "txt")
		as.Equal("test.txt", localFile.Name)
		as.Equal("txt", *localFile.SourceInfo.FileType)
		as.NotEmpty(localFile.SourceInfo.FileBase64)

		// Test BuildAutoUpdateRule
		autoUpdateRule := DocumentUpdateRuleBuildAutoUpdate(24)
		as.Equal(DocumentUpdateTypeAutoUpdate, autoUpdateRule.UpdateType)
		as.Equal(24, autoUpdateRule.UpdateInterval)

		// Test BuildNoAutoUpdateRule
		noAutoUpdateRule := DocumentUpdateRuleBuildNoAuto()
		as.Equal(DocumentUpdateTypeNoAutoUpdate, noAutoUpdateRule.UpdateType)
	})
}

func TestDocumentTypes(t *testing.T) {
	as := assert.New(t)
	t.Run("DocumentFormatType constants", func(t *testing.T) {
		as.Equal(DocumentFormatType(0), DocumentFormatTypeDocument)
		as.Equal(DocumentFormatType(1), DocumentFormatTypeSpreadsheet)
		as.Equal(DocumentFormatType(2), DocumentFormatTypeImage)
	})

	t.Run("DocumentSourceType constants", func(t *testing.T) {
		as.Equal(DocumentSourceType(0), DocumentSourceTypeLocalFile)
		as.Equal(DocumentSourceType(1), DocumentSourceTypeOnlineWeb)
	})

	t.Run("DocumentStatus constants", func(t *testing.T) {
		as.Equal(DocumentStatus(0), DocumentStatusProcessing)
		as.Equal(DocumentStatus(1), DocumentStatusCompleted)
		as.Equal(DocumentStatus(9), DocumentStatusFailed)
	})

	t.Run("DocumentUpdateType constants", func(t *testing.T) {
		as.Equal(DocumentUpdateType(0), DocumentUpdateTypeNoAutoUpdate)
		as.Equal(DocumentUpdateType(1), DocumentUpdateTypeAutoUpdate)
	})
}
