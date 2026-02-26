package coze

import (
	"context"
	"io"
	"net/http"
)

func (r *files) Upload(ctx context.Context, req *UploadFilesReq) (*UploadFilesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/files/upload",
		Body:   req,
		IsFile: true,
	}
	response := new(uploadFilesResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

func (r *files) Retrieve(ctx context.Context, req *RetrieveFilesReq) (*RetrieveFilesResp, error) {
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    "/v1/files/retrieve",
		Body:   req,
	}
	response := new(retrieveFilesResp)
	err := r.core.rawRequest(ctx, request, response)
	return response.Data, err
}

// FileInfo represents information about a file
type FileInfo struct {
	// The ID of the uploaded file.
	ID string `json:"id"`

	// The total byte size of the file.
	Bytes int `json:"bytes"`

	// The upload time of the file, in the format of a 10-digit Unix timestamp in seconds (s).
	CreatedAt int `json:"created_at"`

	// The name of the file.
	FileName string `json:"file_name"`
}

type FileTypes interface {
	io.Reader
	Name() string
}

type implFileInterface struct {
	io.Reader
	fileName string
}

func (r *implFileInterface) Name() string {
	return r.fileName
}

type UploadFilesReq struct {
	File FileTypes `json:"file"`
}

func NewUploadFile(reader io.Reader, fileName string) FileTypes {
	return &implFileInterface{
		Reader:   reader,
		fileName: fileName,
	}
}

// RetrieveFilesReq represents request for retrieving file
type RetrieveFilesReq struct {
	FileID string `query:"file_id" json:"-"`
}

// UploadFilesResp represents response for uploading file
type UploadFilesResp struct {
	baseModel
	FileInfo
}

// RetrieveFilesResp represents response for retrieving file
type RetrieveFilesResp struct {
	baseModel
	FileInfo
}

type uploadFilesResp struct {
	baseResponse
	Data *UploadFilesResp `json:"data"`
}

type retrieveFilesResp struct {
	baseResponse
	Data *RetrieveFilesResp `json:"data"`
}

type files struct {
	core *core
}

func newFiles(core *core) *files {
	return &files{core: core}
}
