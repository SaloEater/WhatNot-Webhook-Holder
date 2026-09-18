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
	// On a miss go-cache hands back a nil interface, and asserting a nil interface to any
	// concrete T panics ("interface conversion: interface {} is nil") — so a bare Get on a cold
	// key used to take the whole HTTP handler down. The older callers only survived by checking
	// Has() first; a Get that is safe to call on a miss removes that trap for everyone.
	if !found {
		var zero T
		return zero, false
	}
	return entity.(T), true
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

func (c *GoCache[T]) Clear() {
	c.Cache.Flush()
}

func CreateCache[T any](duration time.Duration) GoCache[T] {
	c := cache.New(duration-5*time.Minute, duration)
	return GoCache[T]{
		Cache: c,
	}
}
