package go_cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type GoCache[T any] struct {
	Cache *cache.Cache
}

func (c *GoCache[T]) Get(key string) (T, bool) {
	entity, found := c.Cache.Get(key)
	return entity.(T), found
}

func (c *GoCache[T]) Set(key string, entity T) {
	c.Cache.Set(key, entity, cache.DefaultExpiration)
}

func (c *GoCache[T]) Delete(key string) {
	c.Cache.Delete(key)
}

func (c *GoCache[T]) Has(key string) bool {
	_, found := c.Cache.Get(key)
	return found
}

func CreateCache[T any](duration time.Duration) GoCache[T] {
	c := cache.New(duration-5*time.Minute, duration)
	return GoCache[T]{
		Cache: c,
	}
}
