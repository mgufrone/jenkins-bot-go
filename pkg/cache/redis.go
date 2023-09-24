package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

const (
	KeySeparator = ":"
)

type redisCache struct {
	baseCache
	redisCli *redis.Client
}

func (c redisCache) Has(ctx context.Context, key string) bool {
	k := c.makeKey(key)
	v, err := c.redisCli.Exists(ctx, k).Result()
	return err == nil && v != 0
}

func (c redisCache) Get(ctx context.Context, key string, out interface{}) (err error) {
	outString, err := c.redisCli.Get(ctx, c.makeKey(key)).Result()
	if err != nil {
		return
	}
	return json.Unmarshal([]byte(outString), out)
}

func (c redisCache) Set(ctx context.Context, key string, value interface{}, duration TTL) (err error) {
	if duration.resolve() == 0 {
		return
	}
	k := c.makeKey(key)
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.redisCli.Set(ctx, k, string(val), duration.resolve()).Err()
}
func (c redisCache) Delete(ctx context.Context, key string) (err error) {
	return c.redisCli.Del(ctx, c.makeKey(key)).Err()
}

func (c redisCache) CacheOrCreate(ctx context.Context, key string, out interface{}, duration TTL, generator ICacheGenerator) (err error) {
	if c.Has(ctx, key) {
		return c.Get(ctx, key, out)
	}
	resultSet, err := generator(ctx)
	if err != nil {
		return
	}
	err = c.Set(ctx, key, resultSet, duration)
	if err != nil {
		return
	}
	return c.Get(ctx, key, out)
}

func (c redisCache) Clear(ctx context.Context) error {
	return c.ClearByPrefix(ctx, "")
}

func (c redisCache) ClearByPrefix(ctx context.Context, prefix string) error {
	kp := c.makeKey(fmt.Sprintf("%s:%s", prefix, "*"))
	if prefix == "" {
		kp = c.makeKey("*")
	}
	iter := c.redisCli.Scan(ctx, 0, kp, 0).Iterator()
	for iter.Next(ctx) {
		err := c.redisCli.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}

func Redis(client *redis.Client, prefix string) (ICache, error) {
	cache := &redisCache{redisCli: client}
	cache.prefix = []string{prefix}
	cache.separator = KeySeparator
	return cache, nil
}
