package coze

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type datasetsImages struct {
	client *core
}

func newDatasetsImages(core *core) *datasetsImages {
	return &datasetsImages{
		client: core,
	}
}

func (r *datasetsImages) Update(ctx context.Context, req *UpdateDatasetImageReq) (*UpdateDatasetImageResp, error) {
	method := http.MethodPut
	uri := fmt.Sprintf("/v1/datasets/%s/images/%s", req.DatasetID, req.DocumentID)
	resp := &updateImageResp{}
	err := r.client.Request(ctx, method, uri, req, resp)
	if err != nil {
		return nil, err
	}
	result := &UpdateDatasetImageResp{}
	result.setHTTPResponse(resp.HTTPResponse)
	return result, nil
}

func (r *datasetsImages) List(ctx context.Context, req *ListDatasetsImagesReq) (NumberPaged[Image], error) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}

	return NewNumberPaged[Image](
		func(request *pageRequest) (*pageResponse[Image], error) {
			uri := fmt.Sprintf("/v1/datasets/%s/images", req.DatasetID)
			resp := &listImagesResp{}
			var queries []RequestOption
			if req.Keyword != nil {
				queries = append(queries, withHTTPQuery("keyword", *req.Keyword))
			}
			if req.HasCaption != nil {
				queries = append(queries, withHTTPQuery("has_caption", strconv.FormatBool(*req.HasCaption)))
			}
			queries = append(queries,
				withHTTPQuery("page_num", strconv.Itoa(request.PageNum)),
				withHTTPQuery("page_size", strconv.Itoa(request.PageSize)),
			)
			err := r.client.Request(ctx, http.MethodGet, uri, nil, resp, queries...)
			if err != nil {
				return nil, err
			}
			return &pageResponse[Image]{
				Total:   resp.Data.TotalCount,
				HasMore: len(resp.Data.ImagesInfos) >= request.PageSize,
				Data:    resp.Data.ImagesInfos,
				LogID:   resp.HTTPResponse.LogID(),
			}, nil
		}, req.PageSize, req.PageNum)
}

// ImageStatus 表示图片状态
type ImageStatus int

const (
	ImageStatusInProcessing     ImageStatus = 0 // 处理中
	ImageStatusCompleted        ImageStatus = 1 // 已完成
	ImageStatusProcessingFailed ImageStatus = 9 // 处理失败
)

// Image 表示图片信息
type Image struct {
	// The ID of the file.
	DocumentID string `json:"document_id"`

	// The total character count of the file content.
	CharCount int `json:"char_count"`

	// The chunking rules. For detailed instructions, refer to the ChunkStrategy object.
	ChunkStrategy *DocumentChunkStrategy `json:"chunk_strategy"`

	// The upload time of the file, in the format of a 10-digit Unix timestamp.
	CreateTime int `json:"create_time"`

	// The last modified time of the file, in the format of a 10-digit Unix timestamp.
	UpdateTime int `json:"update_time"`

	// The type of file format. Values include:
	// 0: Document type, such as txt, pdf, online web pages, etc.
	// 1: Spreadsheet type, such as xls spreadsheets, etc.
	// 2: Images type, such as png images, etc.
	FormatType DocumentFormatType `json:"format_type"`

	// The number of times the file has been hit in conversations.
	HitCount int `json:"hit_count"`

	// The name of the file.
	Name string `json:"name"`

	// The size of the file in bytes.
	Size int `json:"size"`

	// The number of slices the file has been divided into.
	SliceCount int `json:"slice_count"`

	// The method of uploading the file. Values include:
	// 0: Upload local files.
	// 1: Upload online web pages.
	SourceType DocumentSourceType `json:"source_type"`

	// The processing status of the file. Values include:
	// 0: Processing
	// 1: Completed
	// 9: Processing failed, it is recommended to re-upload
	Status ImageStatus `json:"status"`

	// The caption of the image.
	Caption string `json:"caption"`

	// The ID of the creator.
	CreatorID string `json:"creator_id"`
}

// UpdateDatasetImageReq 表示更新图片的请求
type UpdateDatasetImageReq struct {
	DatasetID  string  `json:"-"`
	DocumentID string  `json:"-"`
	Caption    *string `json:"caption"` // 图片描述
}

// UpdateImageResp 表示更新图片的响应
type updateImageResp struct {
	baseResponse
	Data *UpdateDatasetImageResp `json:"data"`
}

type UpdateDatasetImageResp struct {
	baseModel
}

// ListDatasetsImagesReq 表示列出图片的请求
type ListDatasetsImagesReq struct {
	DatasetID  string  `json:"-"`
	Keyword    *string `json:"keyword,omitempty"`
	HasCaption *bool   `json:"has_caption,omitempty"`
	PageNum    int     `json:"page_num"`
	PageSize   int     `json:"page_size"`
}

// ListImagesResp 表示列出图片的响应
type listImagesResp struct {
	baseResponse
	Data *ListImagesResp `json:"data"`
}

type ListImagesResp struct {
	baseModel
	ImagesInfos []*Image `json:"photo_infos"`
	TotalCount  int      `json:"total_count"`
}
