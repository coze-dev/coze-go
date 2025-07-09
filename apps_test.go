package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockApp(id string) *SimpleApp {
	return &SimpleApp{
		ID:          id,
		Name:        randomString(10),
		Description: randomString(10),
		IconURL:     randomString(10),
		IsPublished: true,
		OwnerUserID: randomString(10),
		UpdatedAt:   1,
	}
}

func TestApp(t *testing.T) {
	as := assert.New(t)
	t.Run("list", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			apps := newApps(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/apps", req.URL.Path)
				as.Equal("1", req.URL.Query().Get("page_num"))
				as.Equal("20", req.URL.Query().Get("page_size"))
				return mockResponse(http.StatusOK, &listAppResp{
					Data: &ListAppResp{
						Total: 2,
						Items: []*SimpleApp{
							mockApp("1"),
							mockApp("2"),
						},
					},
				})
			})))
			resp, err := apps.List(context.Background(), &ListAppReq{})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.False(resp.HasMore())

			items := resp.Items()
			as.Len(items, 2)

			as.Equal("1", items[0].ID)
			as.Equal("2", items[1].ID)
		})

		t.Run("error", func(t *testing.T) {
			apps := newApps(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := apps.List(context.Background(), &ListAppReq{})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})
}
