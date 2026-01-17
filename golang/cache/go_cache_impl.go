package cache

import (
	"context"
	"github.com/bluele/gcache"
)

type goCache struct {
	cache gcache.Cache
}

func (g goCache) Get(ctx context.Context, key string) ([]byte, error) {
	panic("Not implemented")
}

func (g goCache) Set(ctx context.Context, key string, value []byte) error {
	panic("Not implemented")
}

func NewGoCache() Cache {
	return &goCache{}
}
