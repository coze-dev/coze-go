package websockets

import (
	"fmt"
	"net/http"
)

// CozeAuth adapts the coze-go Auth interface to work with WebSocket connections
type CozeAuth struct {
	authProvider interface {
		Auth(req *http.Request) error
	}
}

// NewCozeAuth creates a new CozeAuth adapter
func NewCozeAuth(authProvider interface {
	Auth(req *http.Request) error
}) *CozeAuth {
	return &CozeAuth{
		authProvider: authProvider,
	}
}

// GetAuthHeader returns the authorization header value
func (a *CozeAuth) GetAuthHeader() (string, error) {
	// Create a dummy request to get the auth header
	req, err := http.NewRequest("GET", "https://api.coze.com", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	// Apply authentication
	if err := a.authProvider.Auth(req); err != nil {
		return "", fmt.Errorf("failed to authenticate: %w", err)
	}
	
	// Get the authorization header
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header found")
	}
	
	return authHeader, nil
}