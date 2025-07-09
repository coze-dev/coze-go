package websockets

// WebSocketsClient is the main WebSocket client that provides access to all WebSocket services
type WebSocketsClient struct {
	baseURL string
	auth    Auth
	Audio   *AudioClient
	Chat    *ChatClientBuilder
}

// AudioClient provides access to audio WebSocket services
type AudioClient struct {
	baseURL string
	auth    Auth
}

// ChatClientBuilder provides methods to create chat clients
type ChatClientBuilder struct {
	baseURL string
	auth    Auth
}

// NewWebSocketsClient creates a new WebSockets client
func NewWebSocketsClient(baseURL string, auth Auth) *WebSocketsClient {
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
