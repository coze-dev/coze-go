package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenAuth(t *testing.T) {
	as := assert.New(t)
	t.Run("Token returns fixed access token", func(t *testing.T) {
		expectedToken := "test_access_token"
		auth := NewTokenAuth(expectedToken)

		token, err := auth.Token(context.Background())
		require.NoError(t, err)
		as.Equal(expectedToken, token)
	})
}

func TestJWTAuth(t *testing.T) {
	as := assert.New(t)
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

	t.Run("NewJWTAuth with default TTL", func(t *testing.T) {
		client, err := NewJWTOAuthClient(NewJWTOAuthClientParam{
			ClientID:      "test_client_id",
			PublicKey:     "test_public_key",
			PrivateKeyPEM: testPrivateKey,
		}, WithAuthBaseURL(ComBaseURL))
		require.NoError(t, err)

		auth := NewJWTAuth(client, nil)
		jwtAuth, ok := auth.(*jwtOAuthImpl)
		require.True(t, ok)
		as.Equal(900, jwtAuth.TTL)
		as.Nil(jwtAuth.SessionName)
		as.Nil(jwtAuth.Scope)
	})

	t.Run("NewJWTAuth with custom options", func(t *testing.T) {
		client, err := NewJWTOAuthClient(NewJWTOAuthClientParam{
			ClientID:      "test_client_id",
			PublicKey:     "test_public_key",
			PrivateKeyPEM: testPrivateKey,
		}, WithAuthBaseURL(ComBaseURL))
		require.NoError(t, err)

		sessionName := "test_session"
		scope := BuildBotChat([]string{"bot_id"}, []string{"permission"})
		auth := NewJWTAuth(client, &GetJWTAccessTokenReq{
			TTL:         1800,
			SessionName: &sessionName,
			Scope:       scope,
		})

		jwtAuth, ok := auth.(*jwtOAuthImpl)
		require.True(t, ok)
		as.Equal(1800, jwtAuth.TTL)
		as.Equal(&sessionName, jwtAuth.SessionName)
		as.Equal(scope, jwtAuth.Scope)
	})

	t.Run("Token returns cached token when not expired", func(t *testing.T) {
		client, err := NewJWTOAuthClient(NewJWTOAuthClientParam{
			ClientID:      "test_client_id",
			PublicKey:     "test_public_key",
			PrivateKeyPEM: testPrivateKey,
		}, WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(newHTTPClientWithTransport(func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusOK, &OAuthToken{
					AccessToken: "test_access_token",
					ExpiresIn:   3600,
				})
			})))
		as.Nil(err)
		as.NotNil(client)

		auth := NewJWTAuth(client, nil)

		// 第一次调用，获取新 token
		token1, err := auth.Token(context.Background())
		as.Nil(err)
		as.Equal("test_access_token", token1)

		// 第二次调用，使用缓存的 token
		token2, err := auth.Token(context.Background())
		as.Nil(err)
		as.Equal(token1, token2)
	})

	t.Run("Token refreshes when expired", func(t *testing.T) {
		callCount := 0
		client, err := NewJWTOAuthClient(NewJWTOAuthClientParam{
			ClientID:      "test_client_id",
			PublicKey:     "test_public_key",
			PrivateKeyPEM: testPrivateKey,
		}, WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(newHTTPClientWithTransport(func(req *http.Request) (*http.Response, error) {
				callCount++
				return mockResponse(http.StatusOK, &OAuthToken{
					AccessToken: "test_access_token_" + string(rune(callCount+'0')),
					ExpiresIn:   1, // 设置为1秒后过期
				})
			})))
		as.Nil(err)
		as.NotNil(client)

		auth := NewJWTAuth(client, nil)

		// 第一次调用，获取新 token
		token1, err := auth.Token(context.Background())
		as.Nil(err)
		as.Equal("test_access_token_1", token1)

		// 等待 token 过期
		time.Sleep(2 * time.Second)

		// 第二次调用，token 已过期，获取新 token
		token2, err := auth.Token(context.Background())
		as.Nil(err)
		as.Equal("test_access_token_2", token2)
		as.NotEqual(token1, token2)
	})

	t.Run("Token handles error", func(t *testing.T) {
		client, err := NewJWTOAuthClient(NewJWTOAuthClientParam{
			ClientID:      "test_client_id",
			PublicKey:     "test_public_key",
			PrivateKeyPEM: testPrivateKey,
		}, WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(newHTTPClientWithTransport(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("test error")
			})))
		as.Nil(err)
		as.NotNil(client)

		auth := NewJWTAuth(client, nil)

		_, err = auth.Token(context.Background())
		as.NotNil(err)
	})

	t.Run("Token with specified account_id", func(t *testing.T) {
		client, err := NewJWTOAuthClient(NewJWTOAuthClientParam{
			ClientID:      "test_client_id",
			PublicKey:     "test_public_key",
			PrivateKeyPEM: testPrivateKey,
		}, WithAuthBaseURL(ComBaseURL),
			WithAuthHttpClient(newHTTPClientWithTransport(func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusBadRequest, &OAuthToken{
					AccessToken: "test_access_token",
					ExpiresIn:   3600,
				})
			})))
		as.Nil(err)
		as.NotNil(client)

		auth := NewJWTAuth(client, &GetJWTAccessTokenReq{
			AccountID: ptr(int64(1234567890123456)),
		})
		token1, err := auth.Token(context.Background())
		as.Nil(err)
		as.Equal("test_access_token", token1)
	})

	t.Run("Test get RefreshBefore", func(t *testing.T) {
		as.Equal(int64(30), getRefreshBefore(600))
		as.Equal(int64(10), getRefreshBefore(599))
		as.Equal(int64(10), getRefreshBefore(60))
		as.Equal(int64(5), getRefreshBefore(59))
		as.Equal(int64(5), getRefreshBefore(30))
		as.Equal(int64(0), getRefreshBefore(29))
	})
}
