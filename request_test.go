package coze

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newHTTPClientWithTransport(fn func(req *http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{Transport: newMockTransport(fn)}
}

func newCoreWithTransport(transport http.RoundTripper) *core {
	return newCore(&clientOption{
		baseURL:     CnBaseURL,
		client:      &http.Client{Transport: transport},
		logLevel:    LogLevelDebug,
		logger:      newStdLogger(),
		auth:        NewTokenAuth("token"),
		enableLogID: true,
	})
}

func randomString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

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
	as := assert.New(t)
	// 测试创建客户端
	t.Run("With Custom Doer", func(t *testing.T) {
		customDoer := &mockHTTP{}
		core := newCore(&clientOption{baseURL: "https://api.test.com", client: customDoer})
		as.Equal(customDoer, core.client)
	})

	t.Run("With Nil Doer", func(t *testing.T) {
		core := newCore(&clientOption{baseURL: "https://api.test.com"})
		as.NotNil(core.client)
		_, ok := core.client.(*http.Client)
		as.True(ok)
	})
}
