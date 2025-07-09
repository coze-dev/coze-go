package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudioVoiceprintGroup(t *testing.T) {
	as := assert.New(t)
	t.Run("create", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			groupID := randomString(10)
			groups := newAudioVoiceprintGroups(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodPost, req.Method)
				as.Equal("/v1/audio/voiceprint_groups", req.URL.Path)
				return mockResponse(http.StatusOK, &createVoicePrintGroupResp{
					Data: &CreateVoicePrintGroupResp{
						ID: groupID,
					},
				})
			})))
			resp, err := groups.Create(context.Background(), &CreateVoicePrintGroupReq{
				Name: "test_group",
				Desc: "test_desc",
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.Equal(groupID, resp.ID)
		})
		t.Run("error", func(t *testing.T) {
			groups := newAudioVoiceprintGroups(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := groups.Create(context.Background(), &CreateVoicePrintGroupReq{
				Name: "test_group",
				Desc: "test_desc",
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			groupID := randomString(10)
			groups := newAudioVoiceprintGroups(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodPut, req.Method)
				as.Equal("/v1/audio/voiceprint_groups/"+groupID, req.URL.Path)
				return mockResponse(http.StatusOK, &updateVoicePrintGroupResp{
					Data: &UpdateVoicePrintGroupResp{},
				})
			})))
			resp, err := groups.Update(context.Background(), &UpdateVoicePrintGroupReq{
				GroupID: groupID,
				Name:    ptr("test_group"),
				Desc:    ptr("test_desc"),
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
		})
		t.Run("error", func(t *testing.T) {
			groups := newAudioVoiceprintGroups(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := groups.Update(context.Background(), &UpdateVoicePrintGroupReq{
				GroupID: randomString(10),
				Name:    ptr("test_group"),
				Desc:    ptr("test_desc"),
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			groupID := randomString(10)
			groups := newAudioVoiceprintGroups(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodDelete, req.Method)
				as.Equal("/v1/audio/voiceprint_groups/"+groupID, req.URL.Path)
				return mockResponse(http.StatusOK, &deleteVoicePrintGroupResp{
					Data: &DeleteVoicePrintGroupResp{},
				})
			})))
			resp, err := groups.Delete(context.Background(), &DeleteVoicePrintGroupReq{
				GroupID: groupID,
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
		})
		t.Run("error", func(t *testing.T) {
			groups := newAudioVoiceprintGroups(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := groups.Delete(context.Background(), &DeleteVoicePrintGroupReq{
				GroupID: randomString(10),
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			groupID := randomString(10)
			groups := newAudioVoiceprintGroups(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/audio/voiceprint_groups", req.URL.Path)
				return mockResponse(http.StatusOK, &listVoicePrintGroupResp{
					Data: &ListVoicePrintGroupResp{
						Items: []*VoicePrintGroup{
							{
								ID:   groupID,
								Name: "test_group",
								Desc: "test_desc",
							},
						},
					},
				})
			})))
			resp, err := groups.List(context.Background(), &ListVoicePrintGroupReq{
				PageNum:  1,
				PageSize: 10,
			})
			as.Nil(err)
			as.NotNil(resp)
			as.Len(resp.Items(), 1)
			as.Equal(groupID, resp.Items()[0].ID)
			as.NotEmpty(resp.Response().LogID())
		})
		t.Run("empty req", func(t *testing.T) {
			groupID := randomString(10)
			groups := newAudioVoiceprintGroups(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/audio/voiceprint_groups", req.URL.Path)
				return mockResponse(http.StatusOK, &listVoicePrintGroupResp{
					Data: &ListVoicePrintGroupResp{
						Items: []*VoicePrintGroup{
							{
								ID:   groupID,
								Name: "test_group",
								Desc: "test_desc",
							},
						},
					},
				})
			})))
			resp, err := groups.List(context.Background(), &ListVoicePrintGroupReq{})
			as.Nil(err)
			as.NotNil(resp)
			as.Len(resp.Items(), 1)
			as.Equal(groupID, resp.Items()[0].ID)
			as.NotEmpty(resp.Response().LogID())
		})
		t.Run("error", func(t *testing.T) {
			groups := newAudioVoiceprintGroups(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := groups.List(context.Background(), &ListVoicePrintGroupReq{
				PageNum:  1,
				PageSize: 10,
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})
}
