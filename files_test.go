package coze

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFiles(t *testing.T) {
	// Test Upload method
	t.Run("Upload file success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/files/upload", req.URL.Path)

				// Return mock response
				return mockResponse(http.StatusOK, &uploadFilesResp{
					FileInfo: &UploadFilesResp{
						FileInfo: FileInfo{
							ID:        "file1",
							Bytes:     1024,
							CreatedAt: 1234567890,
							FileName:  "test.txt",
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		files := newFiles(core)

		// Create a test file content
		content := []byte("test file content")
		fileReader := bytes.NewReader(content)
		uploadReq := &UploadFilesReq{
			File: NewUploadFile(fileReader, "test.txt"),
		}

		resp, err := files.Upload(context.Background(), uploadReq)

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "file1", resp.ID)
		assert.Equal(t, 1024, resp.Bytes)
		assert.Equal(t, 1234567890, resp.CreatedAt)
		assert.Equal(t, "test.txt", resp.FileName)
	})

	// Test Retrieve method
	t.Run("Retrieve file success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "/v1/files/retrieve", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "file1", req.URL.Query().Get("file_id"))

				// Return mock response
				return mockResponse(http.StatusOK, &retrieveFilesResp{
					FileInfo: &RetrieveFilesResp{
						FileInfo: FileInfo{
							ID:        "file1",
							Bytes:     1024,
							CreatedAt: 1234567890,
							FileName:  "test.txt",
						},
					},
				})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		files := newFiles(core)

		resp, err := files.Retrieve(context.Background(), &RetrieveFilesReq{
			FileID: "file1",
		})

		require.NoError(t, err)
		assert.Equal(t, "test_log_id", resp.LogID())
		assert.Equal(t, "file1", resp.ID)
		assert.Equal(t, 1024, resp.Bytes)
		assert.Equal(t, 1234567890, resp.CreatedAt)
		assert.Equal(t, "test.txt", resp.FileName)
	})

	// Test Upload method with error
	t.Run("Upload file with error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		files := newFiles(core)

		content := []byte("test file content")
		fileReader := bytes.NewReader(content)
		uploadReq := &UploadFilesReq{
			File: NewUploadFile(fileReader, "test.txt"),
		}
		resp, err := files.Upload(context.Background(), uploadReq)

		require.Error(t, err)
		assert.Nil(t, resp)
	})

	// Test Retrieve method with error
	t.Run("Retrieve file with error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				return mockResponse(http.StatusBadRequest, &baseResponse{})
			},
		}

		core := newCore(&clientOption{baseURL: ComBaseURL, client: &http.Client{Transport: mockTransport}})
		files := newFiles(core)

		resp, err := files.Retrieve(context.Background(), &RetrieveFilesReq{
			FileID: "invalid_file_id",
		})

		require.Error(t, err)
		assert.Nil(t, resp)
	})

	// Test UploadFilesReq
	t.Run("Test UploadFilesReq", func(t *testing.T) {
		content := []byte("test file content")
		fileReader := bytes.NewReader(content)
		uploadReq := NewUploadFile(fileReader, "test.txt")

		assert.Equal(t, "test.txt", uploadReq.Name())

		// Test reading from the request
		buffer := make([]byte, len(content))
		n, err := uploadReq.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, len(content), n)
		assert.Equal(t, content, buffer)
	})
}
