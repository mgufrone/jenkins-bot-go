package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/ristretto"
)

type inMemory struct {
	baseCache
	cache *ristretto.Cache
}

func (i *inMemory) Has(ctx context.Context, key string) bool {
	_, ok := i.cache.Get(i.makeKey(key))
	return ok
}

const (
	maxKeys = 10 << 20
	maxCost = 1 << 32
	buffers = 64
)

func InMemory(prefix string) (ICache, error) {
	dest, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: maxKeys, // 10 mb of keys
		MaxCost:     maxCost, // 1 gb max cost
		BufferItems: buffers,
	})
	if err != nil {
		return nil, err
	}
	cache := &inMemory{cache: dest}
	cache.prefix = []string{prefix}
	cache.separator = KeySeparator
	return cache, nil
}

func (i *inMemory) Get(ctx context.Context, key string, out interface{}) error {
	val, found := i.cache.Get(i.makeKey(key))
	if !found {
		return nil
	}
	return json.Unmarshal(val.([]byte), out)
}

func (i *inMemory) Set(ctx context.Context, key string, value interface{}, duration TTL) error {
	by, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ok := i.cache.SetWithTTL(i.makeKey(key), by, 1, duration.resolve())
	if !ok {
		return fmt.Errorf("failed to set cache %s", key)
	}
	return nil
}

func (i *inMemory) Delete(ctx context.Context, key string) error {
	k := i.makeKey(key)
	if _, ok := i.cache.Get(k); !ok {
		return nil
	}
	i.cache.Del(k)
	return nil
}

func (i *inMemory) CacheOrCreate(ctx context.Context, key string, out interface{}, duration TTL, generator ICacheGenerator) error {
	k := i.makeKey(key)
	if v, ok := i.cache.Get(k); ok {
		return json.Unmarshal(v.([]byte), out)
	}
	v, err := generator(ctx)
	if err != nil {
		return err
	}
	err = i.Set(ctx, key, v, duration)
	if err != nil {
		return err
	}
	i.cache.Wait()
	return i.Get(ctx, key, out)
}

func (i *inMemory) Clear(ctx context.Context) error {
	i.cache.Clear()
	return nil
}

func (i *inMemory) ClearByPrefix(ctx context.Context, prefix string) error {
	// risretto doesn't support clearing cache by prefix
	return i.Clear(ctx)
}
