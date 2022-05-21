package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/go-redis/redis/v8"
)

var ErrNotFound = errors.New("key not found")

//go:generate .bin/mockery --name RedisWrapper
type RedisWrapper interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Del(key string) error
	io.Closer
}

type redisWrap struct {
	client *redis.Client
}

func newRedisWrap(conf *config) *redisWrap {
	return &redisWrap{
		client: redis.NewClient(&redis.Options{
			Addr: conf.RedisURL,
		}),
	}
}

func (r *redisWrap) Get(key string) (string, error) {
	value, err := r.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}

	if err != nil {
		return "", fmt.Errorf("failed to get value for key %s: %w", key, err)
	}

	return value, nil
}

func (r *redisWrap) Set(key, value string) error {
	if _, err := r.client.Set(context.Background(), key, value, 0).Result(); err != nil {
		return fmt.Errorf("failed to set value for key %s: %w", key, err)
	}

	return nil
}

func (r *redisWrap) Del(key string) error {
	if _, err := r.client.Del(context.Background(), key).Result(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return nil
}

func (r *redisWrap) Close() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close redis client: %w", err)
	}

	return nil
}
