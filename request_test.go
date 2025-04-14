package coze

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResponse 用于测试的响应结构
type TestResponse struct {
	Data struct {
		Name string `json:"name"`
	} `json:"data"`
	baseResponse
}

type TestReq struct {
	Test string `json:"test"`
	Data string `json:"data"`
}

func TestClient_Request_Success(t *testing.T) {
	// 准备测试数据
	expectedResp := &TestResponse{
		Data: struct {
			Name string `json:"name"`
		}{
			Name: "test",
		},
	}
	respBody, _ := json.Marshal(expectedResp)

	// 创建 mock 响应
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
		Header:     make(http.Header),
	}
	mockResp.Header.Set(httpLogIDKey, "test-log-id")

	// 创建测试客户端
	core := newCore(&clientOption{baseURL: "https://api.test.com", client: &mockHTTP{
		Response: mockResp,
		Error:    nil,
	}})

	// 执行请求
	var actualResp TestResponse
	actualReq := &TestReq{
		Test: "test",
		Data: "data",
	}
	err := core.Request(context.Background(), http.MethodGet, "/test", actualReq, &actualResp, withHTTPQuery("test", "data"))

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedResp.Code, actualResp.Code)
	assert.Equal(t, expectedResp.Data.Name, actualResp.Data.Name)
	assert.Equal(t, "test-log-id", actualResp.HTTPResponse.LogID())
}

func TestClient_Request_Error(t *testing.T) {
	// 测试 HTTP 错误
	t.Run("HTTP Error", func(t *testing.T) {
		core := newCore(&clientOption{baseURL: "https://api.test.com", client: &mockHTTP{
			Response: nil,
			Error:    errors.New("network error"),
		}})

		var resp TestResponse
		err := core.Request(context.Background(), http.MethodGet, "/test", nil, &resp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "network error")
	})

	// 测试业务错误
	t.Run("Business Error", func(t *testing.T) {
		errorResp := &TestResponse{}
		errorResp.Code = 1001
		errorResp.Msg = "business error"
		respBody, _ := json.Marshal(errorResp)

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respBody)),
			Header:     make(http.Header),
		}
		mockResp.Header.Set(httpLogIDKey, "test-log-id")

		core := newCore(&clientOption{baseURL: "https://api.test.com", client: &mockHTTP{
			Response: mockResp,
			Error:    nil,
		}})

		var resp TestResponse
		err := core.Request(context.Background(), http.MethodGet, "/test", nil, &resp)
		assert.Error(t, err)
		cozeErr, ok := err.(*Error)
		assert.True(t, ok)
		assert.Equal(t, 1001, cozeErr.Code)
		assert.Equal(t, "business error", cozeErr.Message)
	})

	// 测试认证错误
	t.Run("Auth Error", func(t *testing.T) {
		errorResp := &authErrorFormat{
			ErrorCode:    "invalid_token",
			ErrorMessage: "Token is invalid",
		}
		respBody, _ := json.Marshal(errorResp)

		mockResp := &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(bytes.NewReader(respBody)),
			Header:     make(http.Header),
		}
		mockResp.Header.Set(httpLogIDKey, "test-log-id")

		core := newCore(&clientOption{baseURL: "https://api.test.com", client: &mockHTTP{
			Response: mockResp,
			Error:    nil,
		}})

		var resp TestResponse
		err := core.Request(context.Background(), http.MethodGet, "/test", nil, &resp)
		assert.Error(t, err)
		authErr, ok := err.(*AuthError)
		assert.True(t, ok)
		assert.Equal(t, "invalid_token", authErr.Code.String())
		assert.Equal(t, "Token is invalid", authErr.ErrorMessage)
	})
}

func TestClient_UploadFile_Success(t *testing.T) {
	// 准备测试数据
	expectedResp := &TestResponse{
		Data: struct {
			Name string `json:"name"`
		}{
			Name: "uploaded.txt",
		},
	}
	respBody, _ := json.Marshal(expectedResp)

	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
		Header:     make(http.Header),
	}
	mockResp.Header.Set(httpLogIDKey, "test-log-id")

	core := newCore(&clientOption{baseURL: "https://api.test.com", client: &mockHTTP{
		Response: mockResp,
		Error:    nil,
	}})

	// 创建测试文件内容
	fileContent := "test file content"
	fields := map[string]string{
		"field1": "value1",
		"field2": "value2",
	}

	var actualResp TestResponse
	err := core.UploadFile(
		context.Background(),
		"/upload",
		strings.NewReader(fileContent),
		"test.txt",
		fields,
		&actualResp,
		withHTTPHeader("test", "header-value"),
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp.Code, actualResp.Code)
	assert.Equal(t, expectedResp.Data.Name, actualResp.Data.Name)
	assert.Equal(t, "test-log-id", actualResp.HTTPResponse.LogID())
}

func TestClient_UploadFile_Error(t *testing.T) {
	// 测试上传错误
	t.Run("Upload Error", func(t *testing.T) {
		core := newCore(&clientOption{baseURL: "https://api.test.com", client: &mockHTTP{
			Response: nil,
			Error:    errors.New("upload failed"),
		}})

		var resp TestResponse
		err := core.UploadFile(
			context.Background(),
			"/upload",
			strings.NewReader("test"),
			"test.txt",
			nil,
			&resp,
		)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "upload failed")
	})

	// 测试业务错误
	t.Run("Business Error", func(t *testing.T) {
		errorResp := &TestResponse{}
		errorResp.Code = 1002
		errorResp.Msg = "upload business error"
		respBody, _ := json.Marshal(errorResp)

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respBody)),
			Header:     make(http.Header),
		}
		mockResp.Header.Set(httpLogIDKey, "test-log-id")

		core := newCore(&clientOption{baseURL: "https://api.test.com", client: &mockHTTP{
			Response: mockResp,
			Error:    nil,
		}})

		var resp TestResponse
		err := core.UploadFile(
			context.Background(),
			"/upload",
			strings.NewReader("test"),
			"test.txt",
			nil,
			&resp,
		)

		assert.Error(t, err)
		cozeErr, ok := err.(*Error)
		assert.True(t, ok)
		assert.Equal(t, 1002, cozeErr.Code)
		assert.Equal(t, "upload business error", cozeErr.Message)
	})
}

func TestRequestOptions(t *testing.T) {
	// 测试请求选项
	t.Run("withHTTPHeader", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "https://api.test.com", nil)
		opt := withHTTPHeader("X-Test", "test-value")
		err := opt(req)
		assert.NoError(t, err)
		assert.Equal(t, "test-value", req.Header.Get("X-Test"))
	})

	t.Run("withHTTPQuery", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "https://api.test.com", nil)
		opt := withHTTPQuery("param", "value")
		err := opt(req)
		assert.NoError(t, err)
		assert.Equal(t, "value", req.URL.Query().Get("param"))
	})
}

func TestNewClient(t *testing.T) {
	// 测试创建客户端
	t.Run("With Custom Doer", func(t *testing.T) {
		customDoer := &mockHTTP{}
		core := newCore(&clientOption{baseURL: "https://api.test.com", client: customDoer})
		assert.Equal(t, customDoer, core.client)
	})

	t.Run("With Nil Doer", func(t *testing.T) {
		core := newCore(&clientOption{baseURL: "https://api.test.com"})
		assert.NotNil(t, core.client)
		_, ok := core.client.(*http.Client)
		assert.True(t, ok)
	})
}
