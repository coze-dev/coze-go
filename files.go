package coze

import (
	"context"
	"io"
	"net/http"
)

func (r *files) Upload(ctx context.Context, req *UploadFilesReq) (*UploadFilesResp, error) {
	path := "/v1/files/upload"
	resp := &uploadFilesResp{}
	err := r.core.UploadFile(ctx, path, req.File, req.File.Name(), nil, resp)
	if err != nil {
		return nil, err
	}

	resp.FileInfo.setHTTPResponse(resp.HTTPResponse)
	return resp.FileInfo, nil
}

func (r *files) Retrieve(ctx context.Context, req *RetrieveFilesReq) (*RetrieveFilesResp, error) {
	method := http.MethodPost
	uri := "/v1/files/retrieve"
	resp := &retrieveFilesResp{}
	err := r.core.Request(ctx, method, uri, nil, resp, withHTTPQuery("file_id", req.FileID))
	if err != nil {
		return nil, err
	}
	resp.FileInfo.setHTTPResponse(resp.HTTPResponse)
	return resp.FileInfo, nil
}

type files struct {
	core *core
}

func newFiles(core *core) *files {
	return &files{core: core}
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
	File FileTypes
}

func NewUploadFile(reader io.Reader, fileName string) FileTypes {
	return &implFileInterface{
		Reader:   reader,
		fileName: fileName,
	}
}

// RetrieveFilesReq represents request for retrieving file
type RetrieveFilesReq struct {
	FileID string `json:"file_id"`
}

// uploadFilesResp represents response for uploading file
type uploadFilesResp struct {
	baseResponse
	FileInfo *UploadFilesResp `json:"data"`
}

// UploadFilesResp represents response for uploading file
type UploadFilesResp struct {
	baseModel
	FileInfo
}

// retrieveFilesResp represents response for retrieving file
type retrieveFilesResp struct {
	baseResponse
	FileInfo *RetrieveFilesResp `json:"data"`
}

// RetrieveFilesResp represents response for retrieving file
type RetrieveFilesResp struct {
	baseModel
	FileInfo
}
