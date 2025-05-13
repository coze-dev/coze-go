package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestVariablesService_Retrieve(t *testing.T) {
	mockClient := &http.Client{}
	httpmock.ActivateNonDefault(mockClient)
	defer httpmock.DeactivateAndReset()

	cli := NewCozeAPI(NewTokenAuth("test-token"), WithHttpClient(mockClient))

	ctx := context.Background()

	t.Run("Test with valid req", func(t *testing.T) {
		mockResp := `{
  "code": 0,
  "msg": "Success",
  "data": {
    "items": [
      {
        "value": "val1",
        "create_time": 0,
        "update_time": 0,
        "keyword": "key1"
      },
      {
        "update_time": 1744637812,
        "keyword": "key2",
        "value": "val2",
        "create_time": 1744637812
      }
    ]
  },
  "detail": {
    "logid": "20241210152726467C48D89D6DB2****"
  }
}
`

		httpmock.RegisterResponder("GET", "/v1/variables",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, mockResp)
				return resp, nil
			},
		)

		req := &RetrieveVariablesReq{
			ConnectorUID: "test-connector-uid",
			Keywords:     []string{"key1", "key2"},
			AppID:        ptr("test-app-id"),
		}
		res, err := cli.Variables.Retrieve(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Len(t, res.Items, 2)
		assert.Equal(t, "key1", res.Items[0].Keyword)
		assert.Equal(t, "val1", res.Items[0].Value)
		assert.Equal(t, int64(0), res.Items[0].CreateTime)
		assert.Equal(t, int64(0), res.Items[0].UpdateTime)

		assert.Equal(t, "key2", res.Items[1].Keyword)
		assert.Equal(t, "val2", res.Items[1].Value)
		assert.Equal(t, int64(1744637812), res.Items[1].CreateTime)
		assert.Equal(t, int64(1744637812), res.Items[1].UpdateTime)
	})

	t.Run("Test with nil req", func(t *testing.T) {
		_, err := cli.Variables.Retrieve(ctx, nil)
		assert.Error(t, err)
		assert.Equal(t, "invalid req", err.Error())
	})
}

func TestVariablesService_Update(t *testing.T) {
	mockClient := &http.Client{}
	httpmock.ActivateNonDefault(mockClient)
	defer httpmock.DeactivateAndReset()

	cli := NewCozeAPI(NewTokenAuth("test-token"), WithHttpClient(mockClient))

	ctx := context.Background()

	t.Run("Test with valid req", func(t *testing.T) {
		mockResp := `{
  "code": 0,
  "msg": "",
  "detail": {
    "logid": "20250416125552EE59A23A87AD80CA7051"
  }
}
`
		httpmock.RegisterResponder("PUT", "/v1/variables",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, mockResp)
				return resp, nil
			},
		)

		req := &UpdateVariablesReq{
			ConnectorUID: "test-connector-uid",
			Data: []VariableValue{
				{Keyword: "key1", Value: "new_value1"},
				{Keyword: "key2", Value: "new_value2"},
			},
			AppID: ptr("test-app-id"),
		}
		respData, err := cli.Variables.Update(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, respData)
	})

	t.Run("Test with nil req", func(t *testing.T) {
		_, err := cli.Variables.Update(ctx, nil)
		assert.Error(t, err)
		assert.Equal(t, "invalid req", err.Error())
	})
}
