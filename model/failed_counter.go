package model

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewFailedCounter(r *redis.Client) *FailedCounter {
	return &FailedCounter{
		redis: r,
	}
}

type FailedCounterModel interface {
	Increment(context.Context, string) (int, error)
	Reset(context.Context, string) error
	Get(context.Context, string) (int, error)
}

type FailedCounter struct {
	redis *redis.Client
}

func (f *FailedCounter) Increment(ctx context.Context, key string) (int, error) {
	i, err := f.redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func (f *FailedCounter) Reset(ctx context.Context, key string) error {
	return f.redis.Set(ctx, key, 0, 0).Err()
}

func (f *FailedCounter) Get(ctx context.Context, key string) (int, error) {
	v, err := f.redis.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return i, nil
}
