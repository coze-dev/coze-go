package coze

import (
	"net/http"
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
