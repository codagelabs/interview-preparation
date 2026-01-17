package cache

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type cache struct {
	bigCache *bigcache.BigCache
}

func NewCache(ctx context.Context, evictionTime time.Duration) Cache {
	cacheConfig := bigcache.DefaultConfig(evictionTime)
	bigCache, initErr := bigcache.New(context.Background(), cacheConfig)
	if initErr != nil {
		panic(initErr)
	}
	go func() {
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC,
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		select {
		case <-sigC:
			bigCache.Close()
		case <-ctx.Done():
			bigCache.Close()
		}
	}()
	return &cache{
		bigCache: bigCache,
	}

}

func (c cache) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := c.bigCache.Get(key)
	if err != nil {
		return nil, err
	}
	c.Set(ctx, key, value)
	return value, nil

}

func (c cache) Set(ctx context.Context, key string, value []byte) error {
	return c.bigCache.Set(key, value)
}
