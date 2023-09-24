package cache

import (
	"context"
	"encoding/json"
	"github.com/mgufrone/go-utils/try"
	"github.com/mgufrone/jenkins-slackbot/pkg/env"
	nts "github.com/nutsdb/nutsdb"
)

type nuts struct {
	db            *nts.DB
	defaultBucket string
}

func OpenNuts() (*nts.DB, error) {
	opt := nts.DefaultOptions
	opt.Dir = env.GetOr("NUTS_CACHE_DIR", "/tmp/nutsdb")
	return nts.Open(opt)
}

func (b *nuts) Has(ctx context.Context, key string) (res bool) {
	_ = b.db.View(func(tx *nts.Tx) error {
		bucket, err := tx.Get(b.defaultBucket, []byte(key))
		if err == nil {
			res = !bucket.IsZero()
		}
		return err
	})
	return
}

func Nuts(db *nts.DB, defaultBucket string) ICache {
	return &nuts{db: db, defaultBucket: defaultBucket}
}

func (b *nuts) Get(ctx context.Context, key string, out interface{}) error {
	return b.db.View(func(tx *nts.Tx) error {
		var (
			bucket *nts.Entry
		)
		return try.RunOrError(func() (err error) {
			bucket, err = tx.Get(b.defaultBucket, []byte(key))
			return
		}, func() error {
			return json.Unmarshal(bucket.Value, out)
		})
	})
}

func (b *nuts) Set(ctx context.Context, key string, value interface{}, duration TTL) error {
	return b.db.Update(func(tx *nts.Tx) error {
		by, _ := json.Marshal(value)
		return tx.Put(b.defaultBucket, []byte(key), by, uint32(duration.resolve().Seconds()))
	})
}

func (b *nuts) Delete(ctx context.Context, key string) error {
	return b.db.Update(func(tx *nts.Tx) error {
		return tx.Delete(b.defaultBucket, []byte(key))
	})
}

func (b *nuts) CacheOrCreate(ctx context.Context, key string, out interface{}, duration TTL, generator ICacheGenerator) error {
	return b.db.Update(func(tx *nts.Tx) error {
		ent, err := tx.Get(b.defaultBucket, []byte(key))
		if err == nil && ent.Value != nil {
			return json.Unmarshal(ent.Value, out)
		}
		v, err := generator(ctx)
		if err != nil {
			return err
		}
		by, _ := json.Marshal(v)
		err = tx.Put(b.defaultBucket, []byte(key), by, uint32(duration.resolve().Seconds()))
		if err != nil {
			return err
		}
		return json.Unmarshal(by, out)
	})
}

func (b *nuts) Clear(ctx context.Context) error {
	return b.ClearByPrefix(ctx, "")
}

func (b *nuts) ClearByPrefix(ctx context.Context, prefix string) error {
	return b.db.Update(func(tx *nts.Tx) error {
		if entries, err := tx.PrefixScan(b.defaultBucket, []byte(prefix), 0, 1000); err != nil {
			return err
		} else {
			for _, v := range entries {
				err = tx.Delete(b.defaultBucket, v.Key)
				if err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}
		return nil
	})
}
