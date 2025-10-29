package coze

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"testing"
)

func TestStores_Plugin_List(t *testing.T) {
	// Mock response data
	mockResponse := &ListPluginsResp{
		baseResponse: baseResponse{
			HTTPResponse: &httpResponse{
				Status: 200,
				Header: http.Header{
					"X-Log-Id": []string{"test-log-id"},
				},
				ContentLength: 0,
			},
		},
		Data: &ListPluginsData{
			Plugins: []PluginInfo{
				{
					PluginID:     "plugin-123",
					Name:         "Weather Plugin",
					Description:  "Get weather information",
					AvatarURL:    "https://example.com/avatar.jpg",
					Author:       "Coze Team",
					PluginType:   "tool",
					IsOfficial:   true,
					InstallCount: 1000,
					Rating:       4.5,
					Tags:         []string{"weather", "tool"},
					CreateTime:   1234567890,
					UpdateTime:   1234567890,
				},
				{
					PluginID:     "plugin-456",
					Name:         "Translation Plugin",
					Description:  "Translate text between languages",
					AvatarURL:    "https://example.com/avatar2.jpg",
					Author:       "Community",
					PluginType:   "tool",
					IsOfficial:   false,
					InstallCount: 500,
					Rating:       4.2,
					Tags:         []string{"translation", "language"},
					CreateTime:   1234567890,
					UpdateTime:   1234567890,
				},
			},
			TotalCount: 2,
			Page:       1,
			PageSize:   20,
			HasMore:    false,
		},
	}

	// Create test server
	server := createTestServer(t, mockResponse)
	defer server.Close()

	// Create client with test server URL
	auth := NewTokenAuth("test-token")
	client := NewCozeAPI(auth, WithBaseURL("http://"+server.Addr))

	// Test basic list
	t.Run("BasicList", func(t *testing.T) {
		req := &ListPluginsReq{
			Page:     1,
			PageSize: 20,
		}

		resp, err := client.Stores.Plugin.List(context.Background(), req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if resp.Data.TotalCount != 2 {
			t.Errorf("Expected TotalCount to be 2, got %d", resp.Data.TotalCount)
		}

		if len(resp.Data.Plugins) != 2 {
			t.Errorf("Expected 2 plugins, got %d", len(resp.Data.Plugins))
		}

		// Verify first plugin
		plugin1 := resp.Data.Plugins[0]
		if plugin1.PluginID != "plugin-123" {
			t.Errorf("Expected PluginID to be 'plugin-123', got %s", plugin1.PluginID)
		}
		if plugin1.Name != "Weather Plugin" {
			t.Errorf("Expected Name to be 'Weather Plugin', got %s", plugin1.Name)
		}
		if !plugin1.IsOfficial {
			t.Error("Expected first plugin to be official")
		}

		// Verify second plugin
		plugin2 := resp.Data.Plugins[1]
		if plugin2.PluginID != "plugin-456" {
			t.Errorf("Expected PluginID to be 'plugin-456', got %s", plugin2.PluginID)
		}
		if plugin2.Name != "Translation Plugin" {
			t.Errorf("Expected Name to be 'Translation Plugin', got %s", plugin2.Name)
		}
		if plugin2.IsOfficial {
			t.Error("Expected second plugin to not be official")
		}
	})

	// Test with search keyword
	t.Run("SearchPlugins", func(t *testing.T) {
		req := &ListPluginsReq{
			Page:          1,
			PageSize:      10,
			SearchKeyword: "weather",
		}

		resp, err := client.Stores.Plugin.List(context.Background(), req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if resp.Data.TotalCount != 2 {
			t.Errorf("Expected TotalCount to be 2, got %d", resp.Data.TotalCount)
		}
	})

	// Test with official filter
	t.Run("OfficialPlugins", func(t *testing.T) {
		isOfficial := true
		req := &ListPluginsReq{
			Page:       1,
			PageSize:   10,
			IsOfficial: &isOfficial,
		}

		resp, err := client.Stores.Plugin.List(context.Background(), req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(resp.Data.Plugins) != 2 {
			t.Errorf("Expected 2 plugins, got %d", len(resp.Data.Plugins))
		}
	})

	// Test with sorting
	t.Run("SortedPlugins", func(t *testing.T) {
		req := &ListPluginsReq{
			Page:      1,
			PageSize:  10,
			SortBy:    "install_count",
			SortOrder: "desc",
		}

		resp, err := client.Stores.Plugin.List(context.Background(), req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(resp.Data.Plugins) != 2 {
			t.Errorf("Expected 2 plugins, got %d", len(resp.Data.Plugins))
		}

		// Verify sorting (first plugin has more installs)
		if resp.Data.Plugins[0].InstallCount < resp.Data.Plugins[1].InstallCount {
			t.Error("Expected plugins to be sorted by install count in descending order")
		}
	})

	// Test nil request (should use defaults)
	t.Run("NilRequest", func(t *testing.T) {
		resp, err := client.Stores.Plugin.List(context.Background(), nil)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if resp.Data.Page != 1 {
			t.Errorf("Expected default Page to be 1, got %d", resp.Data.Page)
		}
		if resp.Data.PageSize != 20 {
			t.Errorf("Expected default PageSize to be 20, got %d", resp.Data.PageSize)
		}
	})
}

// Helper function to create test server
func createTestServer(t *testing.T, response interface{}) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/store/plugins", func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Verify query parameters
		if r.URL.Query().Get("page") == "" {
			t.Error("Expected 'page' query parameter")
		}
		if r.URL.Query().Get("page_size") == "" {
			t.Error("Expected 'page_size' query parameter")
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Log-Id", "test-log-id")
		w.WriteHeader(http.StatusOK)

		// Write response
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})

	server := &http.Server{
		Addr:    ":0",
		Handler: mux,
	}

	// Start server on a random port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	go server.Serve(listener)

	// Update server URL to include the actual port
	server.Addr = listener.Addr().String()

	return server
}