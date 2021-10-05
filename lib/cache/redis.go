package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/codemaestro64/skeleton/config"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client redis.UniversalClient
	config config.RedisConfig
}

const redisPrefix = "REDIS"

func NewRedis(cfg *config.Config) (Store, error) {
	r := &Redis{
		config: cfg.Redis,
	}

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: cfg.Redis.Addresses,
	})

	ctx, cancel := r.getTimeoutContext()
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis cluster: %s", err.Error())
	}
	r.client = rdb

	return r, nil
}

func (r *Redis) getTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), r.config.Timeout*time.Second)
}

func (r *Redis) Put(key string, data interface{}, duration time.Duration) error {
	ctx, cancel := r.getTimeoutContext()
	defer cancel()

	_, err := r.client.Set(ctx, key, data, duration).Result()
	if err != nil {
		return fmt.Errorf("%s: error putting item: %s", redisPrefix, err.Error())
	}

	return nil
}

func (r *Redis) Has(key string) (bool, error) {
	ctx, cancel := r.getTimeoutContext()
	defer cancel()

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("%s: error checking if item exists: %s", redisPrefix, err.Error())
	}

	return exists > 0, nil
}

func (r *Redis) Get(key string) (interface{}, error) {
	ctx, cancel := r.getTimeoutContext()
	defer cancel()

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: error getting item: %s", redisPrefix, err.Error())
	}

	return val, nil
}

func (r *Redis) Remove(key string) error {
	ctx, cancel := r.getTimeoutContext()
	defer cancel()

	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("%s: error deleting item: %s", redisPrefix, err.Error())
	}

	return nil
}

func (r *Redis) Flush() {
	r.client.FlushAll(context.Background())
}

func (r *Redis) Close() {
	if r.client != nil {
		r.client.Close()
	}
}
