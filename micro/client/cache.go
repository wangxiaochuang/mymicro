package client

import (
	"context"
	"encoding/json"
	"time"

	cache "github.com/patrickmn/go-cache"
)

func NewCache() *Cache {
	return &Cache{
		cache: cache.New(cache.NoExpiration, 30*time.Second),
	}
}

type Cache struct {
	cache *cache.Cache
}

func (c *Cache) Get(ctx context.Context, req *Request) (interface{}, bool) {
	return c.cache.Get(key(ctx, req))
}

func (c *Cache) Set(ctx context.Context, req *Request, rsp interface{}, expiry time.Duration) {
	c.cache.Set(key(ctx, req), rsp, expiry)
}

func (c *Cache) List() map[string]string {
	items := c.cache.Items()

	rsp := make(map[string]string, len(items))

	for k, v := range items {
		bytes, _ := json.Marshal(v.Object)
		rsp[k] = string(bytes)
	}

	return rsp
}

func key(ctx context.Context, req *Request) string {
	panic("key")
}
