package coze

import (
	"context"
	"time"
)

type Auth interface {
	Token(ctx context.Context) (string, error)
}

var (
	_ Auth = &tokenAuthImpl{}
	_ Auth = &jwtOAuthImpl{}
)

// tokenAuthImpl implements the Auth interface with fixed access token.
type tokenAuthImpl struct {
	accessToken string
}

// NewTokenAuth creates a new token authentication instance.
func NewTokenAuth(accessToken string) Auth {
	return &tokenAuthImpl{
		accessToken: accessToken,
	}
}

func NewJWTAuth(client *JWTOAuthClient, opt *GetJWTAccessTokenReq) Auth {
	if opt == nil {
		opt = &GetJWTAccessTokenReq{}
	}
	if opt.TTL <= 0 {
		opt.TTL = 900
	}
	if opt.Store == nil {
		opt.Store = newFixedKeyMemStore()
	}

	return &jwtOAuthImpl{
		TTL:           opt.TTL,
		Scope:         opt.Scope,
		SessionName:   opt.SessionName,
		refreshBefore: getRefreshBefore(opt.TTL),
		client:        client,
		accountID:     opt.AccountID,
		store:         opt.Store,
	}
}

// Token returns the access token.
func (r *tokenAuthImpl) Token(ctx context.Context) (string, error) {
	return r.accessToken, nil
}

func getRefreshBefore(ttl int) int64 {
	if ttl >= 600 {
		return 30
	} else if ttl >= 60 {
		return 10
	} else if ttl >= 30 {
		return 5
	}
	return 0
}

type jwtOAuthImpl struct {
	TTL           int
	SessionName   *string
	Scope         *Scope
	client        *JWTOAuthClient
	refreshBefore int64 // 在到期前多少秒刷新
	accountID     *int64
	store         Store
	storeKey      string
}

func (r *jwtOAuthImpl) Token(ctx context.Context) (string, error) {
	token, _ := r.store.Get(ctx, r.storeKey)
	if token != "" {
		return token, nil
	}

	resp, err := r.client.GetAccessToken(ctx, &GetJWTAccessTokenReq{
		TTL:         r.TTL,
		SessionName: r.SessionName,
		Scope:       r.Scope,
		AccountID:   r.accountID,
	})
	if err != nil {
		return "", err
	}

	// resp.ExpiresIn 是到期的时间戳, 减去 r.refreshBefore buffer 时间, 再减去当前时间, 得到缓存 ttl
	ttl := time.Second * time.Duration(resp.ExpiresIn-r.refreshBefore-time.Now().Unix())
	_ = r.store.Set(ctx, r.storeKey, resp.AccessToken, ttl)

	return resp.AccessToken, nil
}
