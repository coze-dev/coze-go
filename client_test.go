package coze

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	// Test default initialization
	t.Run("default initialization", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		api := NewCozeAPI(auth)

		assert.Equal(t, ComBaseURL, api.baseURL)
		assert.NotNil(t, api.Audio)
		assert.NotNil(t, api.Bots)
		assert.NotNil(t, api.Chat)
		assert.NotNil(t, api.Conversations)
		assert.NotNil(t, api.Workflows)
		assert.NotNil(t, api.Workspaces)
		assert.NotNil(t, api.Datasets)
		assert.NotNil(t, api.Files)
	})

	// Test with custom base URL
	t.Run("custom base URL", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		customURL := "https://custom.api.coze.com"
		api := NewCozeAPI(auth, WithBaseURL(customURL))

		assert.Equal(t, customURL, api.baseURL)
	})

	// Test with custom HTTP core
	t.Run("custom HTTP core", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		customClient := &http.Client{
			Timeout: 30,
		}
		api := NewCozeAPI(auth, WithHttpClient(customClient))

		assert.NotNil(t, api)
	})

	// Test with custom log level
	t.Run("custom log level", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		api := NewCozeAPI(auth, WithLogLevel(LogLevelDebug))

		assert.NotNil(t, api)
	})

	// Test with custom logger
	t.Run("custom logger", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		customLogger := &mockLogger{}
		api := NewCozeAPI(auth, WithLogger(customLogger))

		assert.NotNil(t, api)
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

		assert.Equal(t, customURL, api.baseURL)
		assert.NotNil(t, api)
	})

	t.Run("with logid", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		api := NewCozeAPI(auth, WithEnableLogID(true))

		assert.NotNil(t, api)
	})
}

func TestAuthTransport(t *testing.T) {
	// Test successful authentication
	t.Run("successful authentication", func(t *testing.T) {
		auth := &mockAuth{token: "test_token"}
		transport := &authTransport{
			auth: auth,
			next: &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					// Verify authorization header
					assert.Equal(t, "Bearer test_token", req.Header.Get("Authorization"))
					return &http.Response{StatusCode: http.StatusOK}, nil
				},
			},
		}

		req, _ := http.NewRequest(http.MethodGet, ComBaseURL, nil)
		resp, err := transport.RoundTrip(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Test authentication error
	t.Run("authentication error", func(t *testing.T) {
		auth := &mockAuth{
			token: "",
			err:   assert.AnError,
		}
		transport := &authTransport{
			auth: auth,
			next: http.DefaultTransport,
		}

		req, _ := http.NewRequest(http.MethodGet, ComBaseURL, nil)
		resp, err := transport.RoundTrip(req)

		require.Error(t, err)
		assert.Nil(t, resp)
	})
}

// mockLogger implements log.Logger interface for testing
type mockLogger struct{}

func (m *mockLogger) Log(ctx context.Context, level LogLevel, message string, args ...interface{}) {
}

func (m *mockLogger) Errorf(format string, args ...interface{}) {}
