package coze

import (
	"context"
	"net/http"
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
	request := &RawRequestReq{
		Method: http.MethodPut,
		URL:    "/v1/datasets/:dataset_id/images/:document_id",
		Body:   req,
	}
	response := new(updateImageResp)
	err := r.client.rawRequest(ctx, request, response)
	return response.Data, err
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
			response := new(listImagesResp)
			if err := r.client.rawRequest(ctx, &RawRequestReq{
				Method: http.MethodGet,
				URL:    "/v1/datasets/:dataset_id/images",
				Body:   req.toReq(request),
			}, response); err != nil {
				return nil, err
			}
			return &pageResponse[Image]{
				Total:   response.Data.TotalCount,
				HasMore: len(response.Data.ImagesInfos) >= request.PageSize,
				Data:    response.Data.ImagesInfos,
				LogID:   response.HTTPResponse.LogID(),
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
	DatasetID  string  `path:"dataset_id" json:"-"`
	DocumentID string  `path:"document_id" json:"-"`
	Caption    *string `json:"caption"` // 图片描述
}

type UpdateDatasetImageResp struct {
	baseModel
}

// ListDatasetsImagesReq 表示列出图片的请求
type ListDatasetsImagesReq struct {
	DatasetID  string  `path:"dataset_id" json:"-"`
	Keyword    *string `query:"keyword" json:"-"`
	HasCaption *bool   `query:"has_caption" json:"-"`
	PageNum    int     `query:"page_num" json:"-"`
	PageSize   int     `query:"page_size" json:"-"`
}

type ListImagesResp struct {
	baseModel
	ImagesInfos []*Image `json:"photo_infos"`
	TotalCount  int      `json:"total_count"`
}

type updateImageResp struct {
	baseResponse
	Data *UpdateDatasetImageResp `json:"data"`
}

type listImagesResp struct {
	baseResponse
	Data *ListImagesResp `json:"data"`
}

func (r ListDatasetsImagesReq) toReq(request *pageRequest) *ListDatasetsImagesReq {
	return &ListDatasetsImagesReq{
		DatasetID:  r.DatasetID,
		Keyword:    r.Keyword,
		HasCaption: r.HasCaption,
		PageNum:    request.PageNum,
		PageSize:   request.PageSize,
	}
}
