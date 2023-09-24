package cache

import (
	"context"
	"strings"
	"time"
)

type ICacheGenerator func(ctx context.Context) (interface{}, error)

type ICache interface {
	Has(ctx context.Context, key string) bool
	Get(ctx context.Context, key string, out interface{}) error
	Set(ctx context.Context, key string, value interface{}, duration TTL) error
	Delete(ctx context.Context, key string) error
	CacheOrCreate(ctx context.Context, key string, out interface{}, duration TTL, generator ICacheGenerator) error
	Clear(ctx context.Context) error
	ClearByPrefix(ctx context.Context, prefix string) error
}

type TTL int

const (
	Short     TTL = iota
	Forever       = -1
	Medium        = 1
	Long          = 2
	ExtraLong     = 3
	SuperLong     = 4
)

func (dur TTL) resolve() time.Duration {
	switch dur {
	case Short:
		return time.Second * 60
	case Medium:
		return time.Minute * 10
	case Long:
		return time.Hour
	case ExtraLong:
		return time.Hour * 6
	case SuperLong:
		return time.Hour * 18
	case Forever:
		return -1
	}
	return 0
}

type baseCache struct {
	prefix    []string
	separator string
}

func (b *baseCache) makeKey(str string) string {
	return strings.Join(append(b.prefix, str), b.separator)
}
