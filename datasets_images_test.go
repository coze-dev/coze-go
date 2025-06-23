package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatasetsImages(t *testing.T) {
	as := assert.New(t)
	t.Run("update image success", func(t *testing.T) {
		images := newDatasetsImages(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodPut, req.Method)
			as.Equal("/v1/datasets/123/images/456", req.URL.Path)
			return mockResponse(http.StatusOK, &updateImageResp{})
		})))
		resp, err := images.Update(context.Background(), &UpdateDatasetImageReq{
			DatasetID:  "123",
			DocumentID: "456",
			Caption:    ptr("test caption"),
		})
		as.Nil(err)
		as.NotNil(resp)
		as.NotEmpty(resp.Response().LogID())
	})

	t.Run("list images success", func(t *testing.T) {
		imageList := []*Image{
			{
				DocumentID:    randomString(10),
				CharCount:     0,
				ChunkStrategy: nil,
				CreateTime:    0,
				UpdateTime:    0,
				FormatType:    DocumentFormatTypeImage,
				HitCount:      0,
				Name:          randomString(10),
				Size:          0,
				SliceCount:    0,
				SourceType:    DocumentSourceTypeLocalFile,
				Status:        ImageStatusCompleted,
				Caption:       randomString(10),
				CreatorID:     randomString(10),
			},
			{
				DocumentID:    randomString(10),
				CharCount:     0,
				ChunkStrategy: nil,
				CreateTime:    0,
				UpdateTime:    0,
				FormatType:    DocumentFormatTypeImage,
				HitCount:      0,
				Name:          randomString(10),
				Size:          0,
				SliceCount:    0,
				SourceType:    DocumentSourceTypeLocalFile,
				Status:        ImageStatusCompleted,
				Caption:       randomString(10),
				CreatorID:     randomString(10),
			},
		}
		images := newDatasetsImages(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			as.Equal(http.MethodGet, req.Method)
			as.Equal("/v1/datasets/123/images", req.URL.Path)
			as.Equal("test", req.URL.Query().Get("keyword"))
			as.Equal("true", req.URL.Query().Get("has_caption"))
			as.Equal("1", req.URL.Query().Get("page_num"))
			as.Equal("10", req.URL.Query().Get("page_size"))
			return mockResponse(http.StatusOK, &listImagesResp{
				Data: &ListImagesResp{
					ImagesInfos: imageList,
					TotalCount:  len(imageList),
				},
			})
		})))
		paged, err := images.List(context.Background(), &ListDatasetsImagesReq{
			DatasetID:  "123",
			Keyword:    ptr("test"),
			HasCaption: ptr(true),
		})
		as.Nil(err)
		as.NotNil(paged)
		// as.NotEmpty(paged.Response().LogID())// todo

		items := paged.Items()
		as.Len(items, len(imageList))
		as.Equal(imageList[0].DocumentID, items[0].DocumentID)
		as.Equal(imageList[0].Name, items[0].Name)
		as.Equal(imageList[0].Caption, items[0].Caption)
		as.Equal(ImageStatusCompleted, items[0].Status)
		as.Equal(DocumentFormatTypeImage, items[0].FormatType)
		as.Equal(DocumentSourceTypeLocalFile, items[0].SourceType)

		as.Equal(imageList[1].DocumentID, items[1].DocumentID)
		as.Equal(imageList[1].Name, items[1].Name)
		as.Equal(imageList[1].Caption, items[1].Caption)
		as.Equal(ImageStatusCompleted, items[1].Status)
		as.Equal(DocumentFormatTypeImage, items[1].FormatType)
		as.Equal(DocumentSourceTypeLocalFile, items[1].SourceType)

		as.Equal(2, paged.Total())
		as.False(paged.HasMore())
	})
}
