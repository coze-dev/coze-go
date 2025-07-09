package coze

import (
	"context"
	"net/http"
)

func (r *folders) List(ctx context.Context, req *ListFoldersReq, options ...CozeAPIOption) (NumberPaged[SimpleFolder], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[SimpleFolder], error) {
			response := new(listFoldersResp)
			if err := r.core.rawRequest(ctx, &RawRequestReq{
				Method:  http.MethodGet,
				URL:     "/v1/folders",
				Body:    req.toReq(request),
				options: options,
			}, response); err != nil {
				return nil, err
			}
			return &pageResponse[SimpleFolder]{
				response: response.HTTPResponse,
				Total:    response.Data.TotalCount,
				HasMore:  len(response.Data.Items) >= request.PageSize,
				Data:     response.Data.Items,
				LogID:    response.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

func (r *folders) Retrieve(ctx context.Context, req *RetrieveFolderReq, options ...CozeAPIOption) (*SimpleFolder, error) {
	resp := new(retrieveFolderResp)
	err := r.core.rawRequest(ctx, &RawRequestReq{
		Method:  http.MethodGet,
		URL:     "/v1/folders/:folder_id",
		Body:    req,
		options: options,
	}, resp)
	return resp.Data, err
}

type FolderType string

const (
	FolderTypeDevelopment FolderType = "development" // 项目开发
	FolderTypeLibrary     FolderType = "library"     // 资源库
)

type SimpleFolder struct {
	baseModel
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	WorkspaceID    string     `json:"workspace_id"`
	CreatorUserID  string     `json:"creator_user_id"`
	FolderType     FolderType `json:"folder_type"`
	ParentFolderID *string    `json:"parent_folder_id,omitempty"`
	ChildrenCount  *int64     `json:"children_count,omitempty"`
}

type ListFoldersReq struct {
	WorkspaceID    string     `query:"workspace_id" json:"-"`
	FolderType     FolderType `query:"folder_type" json:"-"`
	ParentFolderID *string    `query:"parent_folder_id" json:"-"`
	PageNum        int        `query:"page_num" json:"-"`
	PageSize       int        `query:"page_size" json:"-"`
}

type ListFoldersResp struct {
	Items      []*SimpleFolder `json:"items"`
	TotalCount int             `json:"total_count"`
}

type RetrieveFolderReq struct {
	FolderID string `path:"folder_id" json:"-"`
}

func (r ListFoldersReq) toReq(page *pageRequest) *ListFoldersReq {
	return &ListFoldersReq{
		WorkspaceID:    r.WorkspaceID,
		FolderType:     r.FolderType,
		ParentFolderID: r.ParentFolderID,
		PageNum:        page.PageNum,
		PageSize:       page.PageSize,
	}
}

type listFoldersResp struct {
	baseResponse
	Data *ListFoldersResp `json:"data"`
}

type retrieveFolderResp struct {
	baseResponse
	Data *SimpleFolder `json:"data"`
}

type folders struct {
	core *core
}

func newFolders(core *core) *folders {
	return &folders{core: core}
}
