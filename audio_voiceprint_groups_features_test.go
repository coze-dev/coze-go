package coze

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudioVoiceprintGroupFeature(t *testing.T) {
	as := assert.New(t)
	t.Run("create", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			groupID := randomString(10)
			featureID := randomString(10)
			features := newAudioVoiceprintGroupsFeatures(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodPost, req.Method)
				as.Equal("/v1/audio/voiceprint_groups/"+groupID+"/features", req.URL.Path)
				return mockResponse(http.StatusOK, &createVoicePrintGroupFeatureResp{
					Data: &CreateVoicePrintGroupFeatureResp{
						ID: featureID,
					},
				})
			})))
			resp, err := features.Create(context.Background(), &CreateVoicePrintGroupFeatureReq{
				GroupID:    groupID,
				Name:       "test_feature",
				File:       NewUploadFile(strings.NewReader("test_file"), "file.wav"),
				SampleRate: ptr(16000),
				Channel:    ptr(1),
				Desc:       ptr("test_desc"),
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
			as.Equal(featureID, resp.ID)
		})
		t.Run("error", func(t *testing.T) {
			features := newAudioVoiceprintGroupsFeatures(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := features.Create(context.Background(), &CreateVoicePrintGroupFeatureReq{
				GroupID: randomString(10),
				Name:    "test_feature",
				File:    NewUploadFile(strings.NewReader("test_file"), "file.wav"),
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			groupID := randomString(10)
			featureID := randomString(10)
			features := newAudioVoiceprintGroupsFeatures(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodPut, req.Method)
				as.Equal("/v1/audio/voiceprint_groups/"+groupID+"/features/"+featureID, req.URL.Path)
				return mockResponse(http.StatusOK, &updateVoicePrintGroupFeatureResp{
					Data: &UpdateVoicePrintGroupFeatureResp{},
				})
			})))
			resp, err := features.Update(context.Background(), &UpdateVoicePrintGroupFeatureReq{
				GroupID:   groupID,
				FeatureID: featureID,
				Name:      ptr("test_feature"),
				Desc:      ptr("test_desc"),
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
		})
		t.Run("error", func(t *testing.T) {
			features := newAudioVoiceprintGroupsFeatures(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := features.Update(context.Background(), &UpdateVoicePrintGroupFeatureReq{
				GroupID:   randomString(10),
				FeatureID: randomString(10),
				Name:      ptr("test_feature"),
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			groupID := randomString(10)
			featureID := randomString(10)
			features := newAudioVoiceprintGroupsFeatures(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodDelete, req.Method)
				as.Equal("/v1/audio/voiceprint_groups/"+groupID+"/features/"+featureID, req.URL.Path)
				return mockResponse(http.StatusOK, &deleteVoicePrintGroupFeatureResp{
					Data: &DeleteVoicePrintGroupFeatureResp{},
				})
			})))
			resp, err := features.Delete(context.Background(), &DeleteVoicePrintGroupFeatureReq{
				GroupID:   groupID,
				FeatureID: featureID,
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
		})
		t.Run("error", func(t *testing.T) {
			features := newAudioVoiceprintGroupsFeatures(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := features.Delete(context.Background(), &DeleteVoicePrintGroupFeatureReq{
				GroupID:   randomString(10),
				FeatureID: randomString(10),
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			groupID := randomString(10)
			featureID := randomString(10)
			features := newAudioVoiceprintGroupsFeatures(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodGet, req.Method)
				as.Equal("/v1/audio/voiceprint_groups/"+groupID+"/features", req.URL.Path)
				return mockResponse(http.StatusOK, &listVoicePrintGroupFeatureResp{
					Data: &listVoicePrintGroupFeatureRespData{
						Items: []*VoicePrintGroupFeature{
							{
								ID:   featureID,
								Name: "test_feature",
								Desc: "test_desc",
							},
						},
					},
				})
			})))
			resp, err := features.List(context.Background(), &ListVoicePrintGroupFeatureReq{
				GroupID:  groupID,
				PageNum:  1,
				PageSize: 10,
			})
			as.Nil(err)
			as.NotNil(resp)
			as.Len(resp.Items(), 1)
			as.Equal(featureID, resp.Items()[0].ID)
			// as.NotEmpty(resp.Response().LogID()) // todo
		})
		t.Run("error", func(t *testing.T) {
			groupID := randomString(10)
			features := newAudioVoiceprintGroupsFeatures(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := features.List(context.Background(), &ListVoicePrintGroupFeatureReq{
				GroupID:  groupID,
				PageNum:  1,
				PageSize: 10,
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})
}
