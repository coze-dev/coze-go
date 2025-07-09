package websockets

import (
	"context"
	"fmt"
)

// CozeAuth adapts the coze-go Auth interface to work with WebSocket connections
type CozeAuth struct {
	authProvider interface {
		Token(ctx context.Context) (string, error)
	}
}

// NewCozeAuth creates a new CozeAuth adapter
func NewCozeAuth(authProvider interface {
	Token(ctx context.Context) (string, error)
},
) *CozeAuth {
	return &CozeAuth{
		authProvider: authProvider,
	}
}

// GetAuthHeader returns the authorization header value
func (a *CozeAuth) GetAuthHeader() (string, error) {
	// Get token from auth provider
	token, err := a.authProvider.Token(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	return "Bearer " + token, nil
}
