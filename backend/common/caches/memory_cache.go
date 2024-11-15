package caches

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type InMemoryCache struct {
	client *ristretto.Cache[string, any]
}

func NewInMemoryCache() *InMemoryCache {

	ristrettoCache, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 128,     // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	return &InMemoryCache{
		client: ristrettoCache,
	}
}

func (c *InMemoryCache) Set(key string, value any, expiration time.Duration) (err error) {

	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	ok := c.client.SetWithTTL(key, data, 1, expiration)
	if !ok {
		err = errors.New("error set cache")
		return
	}

	return nil
}

func (c *InMemoryCache) Get(key string) (data any, err error) {
	result, ok := c.client.Get(key)
	if !ok {
		err = os.ErrNotExist
		return
	}

	data = result

	return
}

func (c *InMemoryCache) Delete(key string) (err error) {
	c.client.Del(key)

	return
}
