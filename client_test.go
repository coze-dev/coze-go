package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockAuth implements Auth interface for testing
type mockAuth struct {
	token string
	err   error
}

func (m *mockAuth) Token(ctx context.Context) (string, error) {
	return m.token, m.err
}

func TestNewCozeAPI(t *testing.T) {
	as := assert.New(t)
	t.Run("default initialization", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		api := NewCozeAPI(auth)

		as.Equal(ComBaseURL, api.baseURL)
		as.NotNil(api.Audio)
		as.NotNil(api.Bots)
		as.NotNil(api.Chat)
		as.NotNil(api.Conversations)
		as.NotNil(api.Workflows)
		as.NotNil(api.Workspaces)
		as.NotNil(api.Datasets)
		as.NotNil(api.Files)
	})

	// Test with custom base URL
	t.Run("custom base URL", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		customURL := "https://custom.api.coze.com"
		api := NewCozeAPI(auth, WithBaseURL(customURL))

		as.Equal(customURL, api.baseURL)
	})

	// Test with custom HTTP core
	t.Run("custom HTTP core", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		customClient := &http.Client{
			Timeout: 30,
		}
		api := NewCozeAPI(auth, WithHttpClient(customClient))

		as.NotNil(api)
	})

	// Test with custom log level
	t.Run("custom log level", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		api := NewCozeAPI(auth, WithLogLevel(LogLevelDebug))

		as.NotNil(api)
	})

	// Test with custom logger
	t.Run("custom logger", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		customLogger := &mockLogger{}
		api := NewCozeAPI(auth, WithLogger(customLogger))

		as.NotNil(api)
	})

	// Test with multiple options
	t.Run("multiple options", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		customURL := "https://custom.api.coze.com"
		customClient := &http.Client{
			Timeout: 30,
		}
		customLogger := &mockLogger{}

		api := NewCozeAPI(auth,
			WithBaseURL(customURL),
			WithHttpClient(customClient),
			WithLogLevel(LogLevelDebug),
			WithLogger(customLogger),
		)

		as.Equal(customURL, api.baseURL)
		as.NotNil(api)
	})

	t.Run("with logid", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		api := NewCozeAPI(auth, WithEnableLogID(true))

		as.NotNil(api)
	})
}

// mockLogger implements log.Logger interface for testing
type mockLogger struct{}

func (m *mockLogger) Log(ctx context.Context, level LogLevel, message string, args ...interface{}) {
}

func (m *mockLogger) Errorf(format string, args ...interface{}) {}
