package coze

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTransport 实现 http.RoundTripper 接口
type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

// mockReadCloser 实现 io.ReadCloser 接口
type mockReadCloser struct {
	*bytes.Buffer
}

func (m mockReadCloser) Close() error {
	return nil
}

// mockResponse 创建一个模拟的 HTTP 响应
func mockResponse(statusCode int, body interface{}) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(jsonBytes)
	mockResp := &http.Response{
		StatusCode: statusCode,
		Body:       mockReadCloser{buffer},
		Header:     make(http.Header),
	}
	mockResp.Header.Set(httpLogIDKey, "test_log_id")
	return mockResp, nil
}

func TestPKCEOAuthClient(t *testing.T) {
	t.Run("GenOAuthURL success", func(t *testing.T) {
		client, err := NewPKCEOAuthClient("test_client_id", WithAuthBaseURL(ComBaseURL))
		require.NoError(t, err)

		resp, err := client.GetOAuthURL(context.Background(), &GetPKCEOAuthURLReq{
			RedirectURI: "https://example.com/callback",
			State:       "test_state",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, resp.CodeVerifier)
		assert.Contains(t, resp.AuthorizationURL, "https://www.coze.com/api/permission/oauth2/authorize")
		assert.Contains(t, resp.AuthorizationURL, "client_id=test_client_id")
		assert.Contains(t, resp.AuthorizationURL, "redirect_uri=https%3A%2F%2Fexample.com%2Fcallback")
		assert.Contains(t, resp.AuthorizationURL, "state=test_state")
		assert.Contains(t, resp.AuthorizationURL, "code_challenge_method=S256")
	})

	t.Run("GenWorkspaceOAuthURL success", func(t *testing.T) {
		client, err := NewPKCEOAuthClient("test_client_id", WithAuthBaseURL(ComBaseURL))
		require.NoError(t, err)

		resp, err := client.GetOAuthURL(context.Background(), &GetPKCEOAuthURLReq{
			RedirectURI: "https://example.com/callback",
			State:       "test_state",
			Method:      CodeChallengeMethodS256.Ptr(),
			WorkspaceID: ptr("workspace_id"),
		})

		require.NoError(t, err)
		assert.NotEmpty(t, resp.CodeVerifier)
		assert.Contains(t, resp.AuthorizationURL, "https://www.coze.com/api/permission/oauth2/workspace_id/workspace_id/authorize")
		assert.Contains(t, resp.AuthorizationURL, "client_id=test_client_id")
		assert.Contains(t, resp.AuthorizationURL, "redirect_uri=https%3A%2F%2Fexample.com%2Fcallback")
		assert.Contains(t, resp.AuthorizationURL, "state=test_state")
		assert.Contains(t, resp.AuthorizationURL, "code_challenge_method=S256")
	})

	t.Run("GetAccessToken success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				expectedToken := &OAuthToken{
					AccessToken:  "test_access_token",
					ExpiresIn:    3600,
					RefreshToken: "test_refresh_token",
				}
				return mockResponse(http.StatusOK, expectedToken)
			},
		}

		client, err := NewPKCEOAuthClient("test_client_id",
			WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(&http.Client{Transport: mockTransport}))
		require.NoError(t, err)

		token, err := client.GetAccessToken(context.Background(), &GetPKCEAccessTokenReq{Code: "test_code", RedirectURI: "https://example.com/callback", CodeVerifier: "test_verifier"})
		require.NoError(t, err)
		assert.Equal(t, "test_access_token", token.AccessToken)
		assert.Equal(t, int64(3600), token.ExpiresIn)
		assert.Equal(t, "test_refresh_token", token.RefreshToken)
	})
}

func TestDeviceOAuthClient(t *testing.T) {
	t.Run("GetDeviceCode success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				expectedResp := &GetDeviceAuthResp{
					DeviceCode:      "test_device_code",
					UserCode:        "test_user_code",
					VerificationURI: "https://www.coze.com/verify",
					ExpiresIn:       1800,
					Interval:        5,
				}
				return mockResponse(http.StatusOK, expectedResp)
			},
		}

		client, err := NewDeviceOAuthClient("test_client_id",
			WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(&http.Client{Transport: mockTransport}))
		require.NoError(t, err)

		resp, err := client.GetDeviceCode(context.Background(), nil)
		require.NoError(t, err)
		assert.Equal(t, "test_device_code", resp.DeviceCode)
		assert.Equal(t, "test_user_code", resp.UserCode)
		assert.Equal(t, "https://www.coze.com/verify", resp.VerificationURI)
		assert.Equal(t, "https://www.coze.com/verify?user_code=test_user_code", resp.VerificationURL)
	})

	t.Run("GetAccessToken with polling", func(t *testing.T) {
		attempts := 0
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				attempts++
				if attempts == 1 {
					return mockResponse(http.StatusBadRequest, &authErrorFormat{
						ErrorCode:    string(AuthorizationPending),
						ErrorMessage: "Authorization pending",
					})
				}
				if attempts == 2 {
					return mockResponse(http.StatusBadRequest, &authErrorFormat{
						ErrorCode:    string(SlowDown),
						ErrorMessage: "slow down",
					})
				}
				return mockResponse(http.StatusOK, &OAuthToken{
					AccessToken:  "test_access_token",
					ExpiresIn:    3600,
					RefreshToken: "test_refresh_token",
				})
			},
		}

		httpClient, err := NewDeviceOAuthClient("test_client_id",
			WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(&http.Client{Transport: mockTransport}))
		require.NoError(t, err)

		token, err := httpClient.GetAccessToken(context.Background(), &GetDeviceOAuthAccessTokenReq{
			DeviceCode: "test_device_code",
			Poll:       true,
		})
		require.NoError(t, err)
		assert.Equal(t, "test_access_token", token.AccessToken)
		assert.Equal(t, 3, attempts)
	})
}

func TestJWTOAuthClient(t *testing.T) {
	const testPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCj1Mlf7zfg/kx4
DHogPkN7gTkAYi7FM6TktFZFHDm8Zs1KiL6WrpU+UTqBiHHhlMVB3RiaJxWH40ia
9OWJvIpM5lCaMnzGNX/4wC+4Pxc3KNoUhijP6ofS4yI5xSpUyMrjl9q95ePBNmmP
Tv+s4uTa2y0e1ZlDHwIWC8InZ5NX65RO+yIF+95gclFkANgp5l7aBHaLiSebYRJT
aluZmS4ZUH06Y9LHkS+QvuvOPaQu3Y+xdgHnzEYtNn83tTmLCBAt2ZYcJi0WIeJZ
acaLsi59N1LH+2ZFtMc6+l7qlB0i4m7Dko+9i9OGtBD4y6rMO85VKUAQTs862O3W
KIsWsKXjAgMBAAECggEAAoxg5uxK9O1WTFg3OOw7QEDoUjHLXWPKQtP8sxNxrFjo
yFcx1WQTdYRXHioasuikNn/Tc6vOyc/bXdnq/jzlXg/pbByaWEH/XwHhHgbNNJXb
JhXfrVlv+zAkGXE9czVYILF1xIcgcKI9zhsYl0IXT1gxMmwO98XX0lisPcHY7IhV
JqSGg9cpLi7agyu4E6xBnK8B7rlk34WOrQf7WElwZ+1bddqA2WLmlls5m3dcJ6IF
kJAEMmHYlkpNBC5fhocui0enfVxDncVghZFMugmY6AtxY8kB2U5Fy1hFHi0Eu9Yg
I9XDJD4S/vzpoKojeAVFr/iQkzTj/ObzeF6gaFWN0QKBgQDlM9l69oQX/p94jr9t
6U2G3BK2NJk/O2j1jcOYX7ud1erdRlfeGJwEpReYQ6Ug+cLc3n3cj8qWg2x2Yw8L
45bVuJPxfJ0KPWI03syb+llAsIY3MC70quNu4b9vDTNS6pN6F4trTvT0Woz0x4vo
i3pz3u3NPnfL1I0EoPKobDf7bwKBgQC2/FbOpXTM7a1UHVgd2y1OKzpGcuC0eOKN
/DO2P24CFCgAdySnzsfLYlIKoU8DYvEndyIVysZav6pNC5PJc0vpJ6Oqg3izXigw
viM5CJhFVxPWrtyMcN02JNUSHNWOaiuCOlZIPQEgYCTUECjE/Xl1COonVS38mO+N
FSF7Z3mSzQKBgEmX+2W7D7Dwpd284AR3m9gIg82TV/1wowPtT/d2DbThQfdopb//
YOEw7UGLvtK2v3XRztHqLZ9kdYgRyHwFyKG5EW/Bll76VLUrMMGIge3+gCnqQ7l1
wW8R9zi+IVOnVFEojDCZepeXF5llFSxG1Lutwedb/nUpO1pYH3IqxVLrAoGBAIVv
MSXzhV7CmrhRxaXP5BOydgZVUwKHfD2pgVQOoPunExxzxSkRIqRvCAB0bJe9mLj8
qMBXY5ldVqRkItqt1tcobrKyuFuj947DuA+o8tDtlKviSzWmP8lxxmY03I3DYgLO
44g95Apl0bVKK1CqvdzYKVeRR72BEH5CwG2qoP6pAoGAUpvD0LSVh171UwQFkT6K
b2mWBz1LV7EYLg4bfmi7wRBUCeEuK16+PEJ63yYUg8cSGTZqOFyRbc4tNf2Ow8BL
gpsiuY9Mn2TnbscpeK841s68IHx4l90Je4tbbjK4E/yv+vgARkiiWQbG0BZSkBjO
qI39/arl6ZhTeQMv7TrpQ6Q=
-----END PRIVATE KEY-----`

	t.Run("GetAccessToken success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusOK, &OAuthToken{
					AccessToken:  "test_access_token",
					ExpiresIn:    3600,
					RefreshToken: "test_refresh_token",
				})
			},
		}

		client, err := NewJWTOAuthClient(NewJWTOAuthClientParam{
			ClientID:      "test_client_id",
			PublicKey:     "test_public_key",
			PrivateKeyPEM: testPrivateKey,
			TTL:           nil,
		}, WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(&http.Client{Transport: mockTransport}))
		require.NoError(t, err)

		token, err := client.GetAccessToken(context.Background(), &GetJWTAccessTokenReq{
			TTL:         900,
			Scope:       BuildBotChat([]string{"bot id"}, []string{"permission id"}),
			SessionName: ptr("session"),
		})
		require.NoError(t, err)
		assert.Equal(t, "test_access_token", token.AccessToken)
	})
}

func TestWebOAuthClient(t *testing.T) {
	t.Run("GetOAuthURL success", func(t *testing.T) {
		client, err := NewWebOAuthClient("test_client_id", "test_client_secret",
			WithAuthBaseURL(ComBaseURL))
		require.NoError(t, err)

		url := client.GetOAuthURL(context.Background(), &GetWebOAuthURLReq{
			RedirectURI: "https://example.com/callback",
			State:       "test_state",
		})
		assert.Contains(t, url, "https://www.coze.com/api/permission/oauth2/authorize")
		assert.Contains(t, url, "client_id=test_client_id")
		assert.Contains(t, url, "redirect_uri=https%3A%2F%2Fexample.com%2Fcallback")
		assert.Contains(t, url, "state=test_state")
	})

	t.Run("GetAccessToken success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusOK, &OAuthToken{
					AccessToken:  "test_access_token",
					ExpiresIn:    3600,
					RefreshToken: "test_refresh_token",
				})
			},
		}

		client, err := NewWebOAuthClient("test_client_id", "test_client_secret",
			WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(&http.Client{Transport: mockTransport}))
		require.NoError(t, err)

		token, err := client.GetAccessToken(context.Background(), &GetWebOAuthAccessTokenReq{
			Code:        "test_code",
			RedirectURI: "https://example.com/callback",
		})
		require.NoError(t, err)
		assert.Equal(t, "test_access_token", token.AccessToken)
	})

	t.Run("RefreshToken success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusOK, &OAuthToken{
					AccessToken:  "new_access_token",
					ExpiresIn:    3600,
					RefreshToken: "new_refresh_token",
				})
			},
		}

		client, err := NewWebOAuthClient("test_client_id", "test_client_secret",
			WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(&http.Client{Transport: mockTransport}))
		require.NoError(t, err)

		token, err := client.RefreshToken(context.Background(), "test_refresh_token")
		require.NoError(t, err)
		assert.Equal(t, "new_access_token", token.AccessToken)
		assert.Equal(t, "new_refresh_token", token.RefreshToken)
	})
}

func TestOAuthError(t *testing.T) {
	t.Run("Handle auth error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusUnauthorized, &authErrorFormat{
					ErrorCode:    "unauthorized",
					ErrorMessage: "Invalid core credentials",
				})
			},
		}

		client, err := NewWebOAuthClient("test_client_id", "test_client_secret",
			WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(&http.Client{Transport: mockTransport}))
		require.NoError(t, err)

		_, err = client.GetAccessToken(context.Background(), &GetWebOAuthAccessTokenReq{
			Code:        "test_code",
			RedirectURI: "https://example.com/callback",
		})
		require.Error(t, err)

		authErr, ok := AsAuthError(err)
		require.True(t, ok)
		assert.Equal(t, "unauthorized", authErr.Code.String())
	})
}

func TestParsePrivateKey(t *testing.T) {
	t.Run("Invalid private key format", func(t *testing.T) {
		_, err := parsePrivateKey("invalid_key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode private key")
	})

	t.Run("Invalid PEM block", func(t *testing.T) {
		_, err := parsePrivateKey("LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCmludmFsaWQga2V5IGNvbnRlbnQKLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQo=")
		assert.Error(t, err)
	})
}
