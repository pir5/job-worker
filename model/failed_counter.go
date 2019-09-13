package model

import (
	"strconv"

	"github.com/go-redis/redis"
)

func NewFailedCounter(r *redis.Client) *FailedCounter {
	return &FailedCounter{
		redis: r,
	}
}

type FailedCounterModel interface {
	Increment(string) (int, error)
	Reset(string) error
	Get(string) (int, error)
}

type FailedCounter struct {
	redis *redis.Client
}

func (f *FailedCounter) Increment(key string) (int, error) {
	i, err := f.redis.Incr(key).Result()
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func (f *FailedCounter) Reset(key string) error {
	return f.redis.Set(key, 0, 0).Err()
}

func (f *FailedCounter) Get(key string) (int, error) {
	v, err := f.redis.Get(key).Result()
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return i, nil
}
