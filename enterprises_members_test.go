package coze

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnterprisesMembers(t *testing.T) {
	as := assert.New(t)

	t.Run("create", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			enterpriseID := randomString(10)
			members := newEnterprisesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodPost, req.Method)
				as.Equal(fmt.Sprintf("/v1/enterprises/%s/members", enterpriseID), req.URL.Path)
				return mockResponse(http.StatusOK, &createEnterpriseMemberResp{
					Data: &CreateEnterpriseMemberResp{},
				})
			})))
			resp, err := members.Create(context.Background(), &CreateEnterpriseMemberReq{
				EnterpriseID: enterpriseID,
				Users: []*EnterpriseMember{
					{
						UserID: "user1",
						Role:   EnterpriseMemberRoleMember,
					},
					{
						UserID: "user2",
						Role:   EnterpriseMemberRoleAdmin,
					},
				},
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
		})

		t.Run("error", func(t *testing.T) {
			enterpriseID := randomString(10)
			members := newEnterprisesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := members.Create(context.Background(), &CreateEnterpriseMemberReq{
				EnterpriseID: enterpriseID,
				Users: []*EnterpriseMember{
					{
						UserID: "user1",
						Role:   EnterpriseMemberRoleMember,
					},
				},
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			enterpriseID := randomString(10)
			userID := randomString(10)
			receiverUserID := randomString(10)
			members := newEnterprisesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodDelete, req.Method)
				as.Equal(fmt.Sprintf("/v1/enterprises/%s/members/%s", enterpriseID, userID), req.URL.Path)
				return mockResponse(http.StatusOK, &deleteEnterpriseMemberResp{
					Data: &DeleteEnterpriseMemberResp{},
				})
			})))
			resp, err := members.Delete(context.Background(), &DeleteEnterpriseMemberReq{
				EnterpriseID:   enterpriseID,
				UserID:         userID,
				ReceiverUserID: receiverUserID,
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
		})

		t.Run("error", func(t *testing.T) {
			enterpriseID := randomString(10)
			userID := randomString(10)
			receiverUserID := randomString(10)
			members := newEnterprisesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := members.Delete(context.Background(), &DeleteEnterpriseMemberReq{
				EnterpriseID:   enterpriseID,
				UserID:         userID,
				ReceiverUserID: receiverUserID,
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			enterpriseID := randomString(10)
			userID := randomString(10)
			members := newEnterprisesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(http.MethodPut, req.Method)
				as.Equal(fmt.Sprintf("/v1/enterprises/%s/members/%s", enterpriseID, userID), req.URL.Path)
				return mockResponse(http.StatusOK, &updateEnterpriseMemberResp{
					Data: &UpdateEnterpriseMemberResp{},
				})
			})))
			resp, err := members.Update(context.Background(), &UpdateEnterpriseMemberReq{
				EnterpriseID: enterpriseID,
				UserID:       userID,
				Role:         EnterpriseMemberRoleAdmin,
			})
			as.Nil(err)
			as.NotNil(resp)
			as.NotEmpty(resp.Response().LogID())
		})

		t.Run("error", func(t *testing.T) {
			enterpriseID := randomString(10)
			userID := randomString(10)
			members := newEnterprisesMembers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
			_, err := members.Update(context.Background(), &UpdateEnterpriseMemberReq{
				EnterpriseID: enterpriseID,
				UserID:       userID,
				Role:         EnterpriseMemberRoleAdmin,
			})
			as.NotNil(err)
			as.Contains(err.Error(), "test error")
		})
	})
}

func TestEnterpriseMemberRole(t *testing.T) {
	as := assert.New(t)
	t.Run("constants", func(t *testing.T) {
		as.Equal("enterprise_admin", string(EnterpriseMemberRoleAdmin))
		as.Equal("enterprise_member", string(EnterpriseMemberRoleMember))
	})
}
