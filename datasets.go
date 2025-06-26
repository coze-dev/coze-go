package coze

import (
	"context"
	"net/http"
)

func (r *datasets) Create(ctx context.Context, req *CreateDatasetsReq) (*CreateDatasetResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/datasets",
		Body:   req,
	}
	response := new(createDatasetResp)
	err := r.client.rawRequest(ctx, request, response)
	return response.Data, err
}

func (r *datasets) List(ctx context.Context, req *ListDatasetsReq) (NumberPaged[Dataset], error) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged(
		func(request *pageRequest) (*pageResponse[Dataset], error) {
			response := new(listDatasetsResp)
			err := r.client.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/datasets",
				Body:   req.toReq(request),
			}, response)
			if err != nil {
				return nil, err
			}
			return &pageResponse[Dataset]{
				response: response.HTTPResponse,
				Total:    response.Data.TotalCount,
				HasMore:  len(response.Data.DatasetList) >= request.PageSize,
				Data:     response.Data.DatasetList,
				LogID:    response.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

func (r *datasets) Update(ctx context.Context, req *UpdateDatasetsReq) (*UpdateDatasetsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPut,
		URL:    "/v1/datasets/:dataset_id",
		Body:   req,
	}
	response := new(updateDatasetResp)
	err := r.client.rawRequest(ctx, request, response)
	return response.Data, err
}

func (r *datasets) Delete(ctx context.Context, req *DeleteDatasetsReq) (*DeleteDatasetsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodDelete,
		URL:    "/v1/datasets/:dataset_id",
		Body:   req,
	}
	response := new(deleteDatasetResp)
	err := r.client.rawRequest(ctx, request, response)
	return response.Data, err
}

func (r *datasets) Process(ctx context.Context, req *ProcessDocumentsReq) (*ProcessDocumentsResp, error) {
	request := &RawRequestReq{
		Method: http.MethodPost,
		URL:    "/v1/datasets/:dataset_id/process",
		Body:   req,
	}
	response := new(processDocumentsResp)
	err := r.client.rawRequest(ctx, request, response)
	return response.Data, err
}

// DatasetStatus 表示数据集状态
type DatasetStatus int

const (
	DatasetStatusEnabled  DatasetStatus = 1
	DatasetStatusDisabled DatasetStatus = 3
)

// Dataset 表示数据集信息
type Dataset struct {
	ID                   string                 `json:"dataset_id"`
	Name                 string                 `json:"name"`
	Description          string                 `json:"description"`
	SpaceID              string                 `json:"space_id"`
	Status               DatasetStatus          `json:"status"`
	FormatType           DocumentFormatType     `json:"format_type"`
	CanEdit              bool                   `json:"can_edit"`
	IconURL              string                 `json:"icon_url"`
	DocCount             int                    `json:"doc_count"`
	FileList             []string               `json:"file_list"`
	HitCount             int                    `json:"hit_count"`
	BotUsedCount         int                    `json:"bot_used_count"`
	SliceCount           int                    `json:"slice_count"`
	AllFileSize          string                 `json:"all_file_size"`
	ChunkStrategy        *DocumentChunkStrategy `json:"chunk_strategy,omitempty"`
	FailedFileList       []string               `json:"failed_file_list"`
	ProcessingFileList   []string               `json:"processing_file_list"`
	ProcessingFileIDList []string               `json:"processing_file_id_list"`
	AvatarURL            string                 `json:"avatar_url"`
	CreatorID            string                 `json:"creator_id"`
	CreatorName          string                 `json:"creator_name"`
	CreateTime           int                    `json:"create_time"`
	UpdateTime           int                    `json:"update_time"`
}

// CreateDatasetsReq 表示创建数据集的请求
type CreateDatasetsReq struct {
	Name        string             `json:"name"`
	SpaceID     string             `json:"space_id"`
	FormatType  DocumentFormatType `json:"format_type"`
	Description string             `json:"description,omitempty"`
	IconFileID  string             `json:"file_id,omitempty"`
}

type CreateDatasetResp struct {
	baseModel
	DatasetID string `json:"dataset_id"`
}

// ListDatasetsReq 表示列出数据集的请求
type ListDatasetsReq struct {
	SpaceID    string             `query:"space_id" json:"-"`
	Name       string             `query:"name,omitempty" json:"-"`
	FormatType DocumentFormatType `query:"format_type,omitempty" json:"-"`
	PageNum    int                `query:"page_num" json:"-"`
	PageSize   int                `query:"page_size" json:"-"`
}

func NewListDatasetsReq(spaceID string) *ListDatasetsReq {
	return &ListDatasetsReq{
		SpaceID:  spaceID,
		PageNum:  1,
		PageSize: 10,
	}
}

type ListDatasetsResp struct {
	baseModel
	TotalCount  int        `json:"total_count"`
	DatasetList []*Dataset `json:"dataset_list"`
}

// UpdateDatasetsReq 表示更新数据集的请求
type UpdateDatasetsReq struct {
	DatasetID   string `path:"dataset_id" json:"-"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IconFileID  string `json:"file_id,omitempty"`
}

type UpdateDatasetsResp struct {
	baseModel
}

// DeleteDatasetsReq 表示删除数据集的请求
type DeleteDatasetsReq struct {
	DatasetID string `path:"dataset_id" json:"-"`
}

type DeleteDatasetsResp struct {
	baseModel
}

// DocumentProgress 表示文档处理进度
type DocumentProgress struct {
	DocumentID     string             `json:"document_id"`
	URL            string             `json:"url"`
	Size           int                `json:"size"`
	Type           string             `json:"type"`
	Status         DocumentStatus     `json:"status"`
	Progress       int                `json:"progress"`
	UpdateType     DocumentUpdateType `json:"update_type"`
	DocumentName   string             `json:"document_name"`
	RemainingTime  int                `json:"remaining_time"`
	StatusDescript string             `json:"status_descript"`
	UpdateInterval int                `json:"update_interval"`
}

// ProcessDocumentsReq 表示处理文档的请求
type ProcessDocumentsReq struct {
	DatasetID   string   `path:"dataset_id" json:"-"`
	DocumentIDs []string `json:"document_ids"`
}

type ProcessDocumentsResp struct {
	baseModel
	Data []*DocumentProgress `json:"data"`
}

type processDocumentsResp struct {
	baseResponse
	Data *ProcessDocumentsResp `json:"data"`
}

type deleteDatasetResp struct {
	baseResponse
	Data *DeleteDatasetsResp `json:"data"`
}

type createDatasetResp struct {
	baseResponse
	Data *CreateDatasetResp `json:"data"`
}

type listDatasetsResp struct {
	baseResponse
	Data *ListDatasetsResp `json:"data"`
}

type updateDatasetResp struct {
	baseResponse
	Data *UpdateDatasetsResp `json:"data"`
}

func (r ListDatasetsReq) toReq(request *pageRequest) *ListDatasetsReq {
	return &ListDatasetsReq{
		SpaceID:    r.SpaceID,
		Name:       r.Name,
		FormatType: r.FormatType,
		PageNum:    request.PageNum,
		PageSize:   request.PageSize,
	}
}

type datasets struct {
	client    *core
	Documents *datasetsDocuments
	Images    *datasetsImages
}

func newDatasets(core *core) *datasets {
	return &datasets{
		client:    core,
		Documents: newDatasetsDocuments(core),
		Images:    newDatasetsImages(core),
	}
}
