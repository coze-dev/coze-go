package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockFolder(id string) *SimpleFolder {
	return &SimpleFolder{
		ID:          id,
		Name:        randomString(10),
		Description: randomString(10),
		FolderType:  FolderTypeLibrary,
	}
}

func TestFolders(t *testing.T) {
	as := assert.New(t)
	t.Run("list", func(t *testing.T) {
		t.Run("list folder success", func(t *testing.T) {
			folders := newFolders(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/folders", req.URL.Path)
				as.Equal("1", req.URL.Query().Get("page_num"))
				as.Equal("20", req.URL.Query().Get("page_size"))
				return mockResponse(http.StatusOK, &listFoldersResp{
					Data: &ListFoldersResp{
						TotalCount: 2,
						Items: []*SimpleFolder{
							mockFolder("f1"),
							mockFolder("f2"),
						},
					},
				})
			})))
			resp, err := folders.List(context.Background(), &ListFoldersReq{})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.False(resp.HasMore())

			items := resp.Items()
			as.Len(items, 2)

			as.Equal("f1", items[0].ID)
			as.Equal("f2", items[1].ID)
		})

		t.Run("list folders with default pagination", func(t *testing.T) {
			folders := newFolders(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal("1", req.URL.Query().Get("page_num"))
				as.Equal("20", req.URL.Query().Get("page_size"))
				return mockResponse(http.StatusOK, &listFoldersResp{
					Data: &ListFoldersResp{
						TotalCount: 0,
						Items:      []*SimpleFolder{},
					},
				})
			})))
			resp, err := folders.List(context.Background(), &ListFoldersReq{})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.False(resp.HasMore())
			as.Empty(resp.Items())
		})

		t.Run("list folders with error", func(t *testing.T) {
			folders := newFolders(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := folders.List(context.Background(), &ListFoldersReq{
				PageNum:  1,
				PageSize: 20,
			})
			as.NotNil(err)
		})
	})

	t.Run("retrieve", func(t *testing.T) {
		t.Run("retrieve success", func(t *testing.T) {
			id := randomString(10)
			folders := newFolders(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/folders/"+id, req.URL.Path)
				return mockResponse(http.StatusOK, &retrieveFolderResp{
					Data: mockFolder(id),
				})
			})))
			resp, err := folders.Retrieve(context.Background(), &RetrieveFolderReq{
				FolderID: id,
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.Equal(id, resp.ID)
			as.NotEmpty(resp.Name)
			as.NotEmpty(resp.Description)
			as.Equal(FolderTypeLibrary, resp.FolderType)
		})
	})
}
