package coze

import (
	"context"
	"errors"
	"time"
)

// Cache 分布式缓存接口
type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, val string, ttl time.Duration) error
}

var ErrStoreNotFound = errors.New("store not found")

func newFixedKeyMemStore() Store {
	return &fixedKeyMemStore{}
}

type fixedKeyMemStore struct {
	val         string
	expiredAtMs int64
}

func (c *fixedKeyMemStore) Get(ctx context.Context, key string) (string, error) {
	if time.Now().UnixMilli() > c.expiredAtMs {
		return "", ErrStoreNotFound
	}
	return c.val, nil
}

func (c *fixedKeyMemStore) Set(ctx context.Context, key, val string, ttl time.Duration) error {
	c.val = val
	c.expiredAtMs = time.Now().Add(ttl).UnixMilli()
	return nil
}
