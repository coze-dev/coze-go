package coze

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Variables(t *testing.T) {
	as := assert.New(t)

	t.Run("retrieve", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			connectorUID := randomString(10)
			keywords := []string{randomString(10), randomString(10)}
			appID := randomString(10)
			variables := newVariables(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				as.Equal(connectorUID, req.URL.Query().Get("connector_uid"))
				as.Equal(strings.Join(keywords, ","), req.URL.Query().Get("keywords"))
				as.Equal(appID, req.URL.Query().Get("app_id"))
				return mockResponse(http.StatusOK, &retrieveVariablesResp{
					Data: &RetrieveVariablesResp{
						Items: []*VariableValue{
							{Keyword: keywords[0], Value: keywords[0]},
						},
					},
				})
			})))
			res, err := variables.Retrieve(context.Background(), &RetrieveVariablesReq{
				ConnectorUID: connectorUID,
				Keywords:     keywords,
				AppID:        ptr(appID),
			})
			as.Nil(err)
			as.NotNil(res)
			as.NotEmpty(res.Response().LogID())
			as.Len(res.Items, 1)
			as.Equal(keywords[0], res.Items[0].Keyword)
			as.Equal(keywords[0], res.Items[0].Value)
			as.Equal(int64(0), res.Items[0].CreateTime)
			as.Equal(int64(0), res.Items[0].UpdateTime)
		})

		t.Run("req nil", func(t *testing.T) {
			variables := newVariables(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				if req.URL.Query().Get("connector_uid") == "" {
					return nil, fmt.Errorf("invalid req")
				}
				return mockResponse(http.StatusOK, &retrieveVariablesResp{
					Data: &RetrieveVariablesResp{
						Items: []*VariableValue{},
					},
				})
			})))
			_, err := variables.Retrieve(context.Background(), nil)
			as.NotNil(err)
			as.Contains(err.Error(), "invalid req")
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			connectorUID := randomString(10)
			keywords := []string{randomString(10), randomString(10)}
			variables := newVariables(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusOK, &updateVariablesResp{})
			})))
			res, err := variables.Update(context.Background(), &UpdateVariablesReq{
				ConnectorUID: connectorUID,
				Data:         []VariableValue{{Keyword: keywords[0], Value: keywords[0]}},
			})
			as.Nil(err)
			as.NotNil(res)
			as.NotEmpty(res.Response().LogID())
		})
	})
}
