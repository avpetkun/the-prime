package cache

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

func isNil(err error) bool {
	return errors.Is(err, redis.Nil)
}

type Cache struct {
	Close func() error

	c redis.Cmdable
}

func Connect(ctx context.Context, cfg Config) (Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pass,
		DB:       cfg.DB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return Cache{}, err
	}
	return Cache{Close: rdb.Close, c: rdb}, nil
}

func (c Cache) Tx(ctx context.Context, callback func(Cache) error) error {
	return c.pipelined(ctx, true, callback)
}

func (c Cache) Batch(ctx context.Context, callback func(Cache) error) error {
	return c.pipelined(ctx, false, callback)
}

func (c Cache) pipelined(ctx context.Context, tx bool, callback func(Cache) error) error {
	if c.Close == nil {
		return callback(c)
	}

	var pipe redis.Pipeliner
	if tx {
		pipe = c.c.TxPipeline()
	} else {
		pipe = c.c.Pipeline()
	}

	err := callback(Cache{c: pipe})
	if err == nil {
		_, err = pipe.Exec(ctx)
	}

	return err
}
