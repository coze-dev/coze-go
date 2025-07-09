package websockets

// WebSocketsClient is the main WebSocket client that provides access to all WebSocket services
type WebSocketsClient struct {
	baseURL string
	auth    Auth
	Audio   *AudioClient
	Chat    *ChatClientBuilder
}

// NewWebSockets creates a new WebSockets client
func NewWebSockets(baseURL string, auth Auth) *WebSocketsClient {
	// Convert HTTP URL to WebSocket URL
	wsURL := baseURL
	if baseURL[:7] == "http://" {
		wsURL = "ws://" + baseURL[7:]
	} else if baseURL[:8] == "https://" {
		wsURL = "wss://" + baseURL[8:]
	}

	return &WebSocketsClient{
		baseURL: wsURL,
		auth:    auth,
		Audio: &AudioClient{
			baseURL: wsURL,
			auth:    auth,
		},
		Chat: &ChatClientBuilder{
			baseURL: wsURL,
			auth:    auth,
		},
	}
}
