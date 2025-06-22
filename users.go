package coze

import (
	"context"
	"net/http"
)

// Me retrieves the current user's information
func (r *users) Me(ctx context.Context) (*User, error) {
	request := &RawRequestReq{
		Method: http.MethodGet,
		URL:    "/v1/users/me",
		Body:   new(GetUserMeReq),
	}
	response := new(meResp)
	err := r.client.rawRequest(ctx, request, response)
	return response.User, err
}

type GetUserMeReq struct{}

// User represents a Coze user
type User struct {
	baseModel
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	NickName  string `json:"nick_name"`
	AvatarURL string `json:"avatar_url"`
}

type meResp struct {
	baseResponse
	User *User `json:"data"`
}

type users struct {
	client *core
}

func newUsers(core *core) *users {
	return &users{
		client: core,
	}
}
