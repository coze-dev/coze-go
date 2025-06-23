package coze

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiles(t *testing.T) {
	as := assert.New(t)
	t.Run("upload file success", func(t *testing.T) {
		files := newFiles(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPost, req.Method)
			as.Equal("/v1/files/upload", req.URL.Path)
			return mockResponse(http.StatusOK, &uploadFilesResp{
				Data: &UploadFilesResp{
					FileInfo: FileInfo{
						ID:        "file1",
						Bytes:     1024,
						CreatedAt: 1234567890,
						FileName:  "test.txt",
					},
				},
			})
		})))
		resp, err := files.Upload(context.Background(), &UploadFilesReq{
			File: NewUploadFile(strings.NewReader("test file content"), "test.txt"),
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal("file1", resp.ID)
		as.Equal(1024, resp.Bytes)
		as.Equal(1234567890, resp.CreatedAt)
		as.Equal("test.txt", resp.FileName)
	})

	t.Run("retrieve file success", func(t *testing.T) {
		files := newFiles(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/files/retrieve", req.URL.Path)
			as.Equal("file1", req.URL.Query().Get("file_id"))
			return mockResponse(http.StatusOK, &retrieveFilesResp{
				Data: &RetrieveFilesResp{
					FileInfo: FileInfo{
						ID:        "file1",
						Bytes:     1024,
						CreatedAt: 1234567890,
						FileName:  "test.txt",
					},
				},
			})
		})))
		resp, err := files.Retrieve(context.Background(), &RetrieveFilesReq{
			FileID: "file1",
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
		as.Equal("file1", resp.ID)
		as.Equal(1024, resp.Bytes)
		as.Equal(1234567890, resp.CreatedAt)
		as.Equal("test.txt", resp.FileName)
	})

	t.Run("upload file with error", func(t *testing.T) {
		files := newFiles(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("test error")
		})))
		_, err := files.Upload(context.Background(), &UploadFilesReq{
			File: NewUploadFile(strings.NewReader("test file content"), "test.txt"),
		})
		as.NotNil(err)
	})

	t.Run("retrieve file with error", func(t *testing.T) {
		files := newFiles(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("test error")
		})))
		_, err := files.Retrieve(context.Background(), &RetrieveFilesReq{
			FileID: "invalid_file_id",
		})
		as.NotNil(err)
	})

	t.Run("test upload files req", func(t *testing.T) {
		uploadReq := NewUploadFile(strings.NewReader("test file content"), "test.txt")
		as.Equal("test.txt", uploadReq.Name())

		buffer := make([]byte, 1024)
		_, err := uploadReq.Read(buffer)
		as.Nil(err)
	})
}
