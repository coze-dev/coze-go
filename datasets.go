package coze

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

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

func (r *datasets) Create(ctx context.Context, req *CreateDatasetsReq) (*CreateDatasetResp, error) {
	method := http.MethodPost
	uri := "/v1/datasets"
	resp := &createDatasetResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	resp.Data.setHTTPResponse(resp.HTTPResponse)
	return resp.Data, nil
}

func (r *datasets) List(ctx context.Context, req *ListDatasetsReq) (NumberPaged[Dataset], error) {
	if req.PageSize == 0 {
		req.PageSize = 10 // 设置默认值为10
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged[Dataset](
		func(request *pageRequest) (*pageResponse[Dataset], error) {
			uri := "/v1/datasets"
			resp := &listDatasetsResp{}
			queries := []RequestOption{}
			if req.SpaceID != "" {
				queries = append(queries, withHTTPQuery("space_id", req.SpaceID))
			}
			if req.Name != "" {
				queries = append(queries, withHTTPQuery("name", req.Name))
			}
			if req.FormatType != 0 {
				queries = append(queries, withHTTPQuery("format_type", strconv.Itoa(int(req.FormatType))))
			}
			queries = append(queries,
				withHTTPQuery("page_num", strconv.Itoa(request.PageNum)),
				withHTTPQuery("page_size", strconv.Itoa(request.PageSize)),
			)
			err := r.client.Request(ctx, http.MethodGet, uri, nil, resp, queries...)
			if err != nil {
				return nil, err
			}
			return &pageResponse[Dataset]{
				Total:   resp.Data.TotalCount,
				HasMore: len(resp.Data.DatasetList) >= request.PageSize,
				Data:    resp.Data.DatasetList,
				LogID:   resp.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

func (r *datasets) Update(ctx context.Context, req *UpdateDatasetsReq) (*UpdateDatasetsResp, error) {
	method := http.MethodPut
	uri := fmt.Sprintf("/v1/datasets/%s", req.DatasetID)
	resp := &updateDatasetResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	result := &UpdateDatasetsResp{}
	result.setHTTPResponse(resp.HTTPResponse)
	return result, nil
}

func (r *datasets) Delete(ctx context.Context, req *DeleteDatasetsReq) (*DeleteDatasetsResp, error) {
	method := http.MethodDelete
	uri := fmt.Sprintf("/v1/datasets/%s", req.DatasetID)
	resp := &deleteDatasetResp{}
	err := r.client.Request(ctx, method, uri, nil, resp)
	if err != nil {
		return nil, err
	}
	result := &DeleteDatasetsResp{}
	result.setHTTPResponse(resp.HTTPResponse)
	return result, nil
}

func (r *datasets) Process(ctx context.Context, req *ProcessDocumentsReq) (*ProcessDocumentsResp, error) {
	method := http.MethodPost
	uri := fmt.Sprintf("/v1/datasets/%s/process", req.DatasetID)
	resp := &processDocumentsResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	resp.Data.setHTTPResponse(resp.HTTPResponse)
	return resp.Data, nil
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

// CreateDatasetResp 表示创建数据集的响应
type createDatasetResp struct {
	baseResponse
	Data *CreateDatasetResp `json:"data"`
}

type CreateDatasetResp struct {
	baseModel
	DatasetID string `json:"dataset_id"`
}

// ListDatasetsReq 表示列出数据集的请求
type ListDatasetsReq struct {
	SpaceID    string             `json:"space_id"`
	Name       string             `json:"name,omitempty"`
	FormatType DocumentFormatType `json:"format_type,omitempty"`
	PageNum    int                `json:"page_num"`
	PageSize   int                `json:"page_size"`
}

func NewListDatasetsReq(spaceID string) *ListDatasetsReq {
	return &ListDatasetsReq{
		SpaceID:  spaceID,
		PageNum:  1,
		PageSize: 10,
	}
}

// ListDatasetsResp 表示列出数据集的响应
type listDatasetsResp struct {
	baseResponse
	Data *ListDatasetsResp `json:"data"`
}

type ListDatasetsResp struct {
	baseModel
	TotalCount  int        `json:"total_count"`
	DatasetList []*Dataset `json:"dataset_list"`
}

// UpdateDatasetsReq 表示更新数据集的请求
type UpdateDatasetsReq struct {
	DatasetID   string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IconFileID  string `json:"file_id,omitempty"`
}

// UpdateDatasetResp 表示更新数据集的响应
type updateDatasetResp struct {
	baseResponse
	Data *UpdateDatasetsResp `json:"data"`
}

type UpdateDatasetsResp struct {
	baseModel
}

// DeleteDatasetsReq 表示删除数据集的请求
type DeleteDatasetsReq struct {
	DatasetID string `json:"-"`
}

// DeleteDatasetResp 表示删除数据集的响应
type deleteDatasetResp struct {
	baseResponse
	Data *DeleteDatasetsResp `json:"data"`
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
	DatasetID   string   `json:"-"`
	DocumentIDs []string `json:"document_ids"`
}

// ProcessDocumentsResp 表示处理文档的响应
type processDocumentsResp struct {
	baseResponse
	Data *ProcessDocumentsResp `json:"data"`
}

type ProcessDocumentsResp struct {
	baseModel
	Data []*DocumentProgress `json:"data"`
}
