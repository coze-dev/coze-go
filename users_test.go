package coze

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Users(t *testing.T) {
	as := assert.New(t)

	t.Run("failed", func(t *testing.T) {
		users := newUsers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("http error")
		})))
		_, err := users.Me(context.Background())
		as.NotNil(err)
		as.Contains(err.Error(), "http error")
	})

	t.Run("me", func(t *testing.T) {
		mockUser := &User{
			UserID:    randomString(10),
			UserName:  randomString(10),
			NickName:  randomString(10),
			AvatarURL: randomString(10),
		}
		users := newUsers(newCoreWithTransport(newMockTransport(func(req *http.Request) (*http.Response, error) {
			return mockResponse(http.StatusOK, &meResp{
				User: mockUser,
			})
		})))
		user, err := users.Me(context.Background())
		as.Nil(err)
		as.Equal(mockUser.UserID, user.UserID)
		as.Equal(mockUser.UserName, user.UserName)
		as.Equal(mockUser.NickName, user.NickName)
		as.Equal(mockUser.AvatarURL, user.AvatarURL)
	})
}
