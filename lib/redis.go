package lib

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

func ConnectRedis(cfg config.RedisConfig) (*Redis, error) {
	r := &Redis{
		config: cfg,
	}

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: cfg.Addresses,
	})

	if err := rdb.Ping(r.getTimeoutContext()).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis cluster: %s", err.Error())
	}
	r.client = rdb

	return r, nil
}

func (r *Redis) Get(key string) (interface{}, error) {
	val, err := r.client.Get(r.getTimeoutContext(), key).Result()

	return val, err
}

func (r *Redis) Put(key string, data interface{}) {
	//r.client.Set(key)
}

func (r *Redis) Disconnect() {
	r.client.Close()
}

func (r *Redis) getTimeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.config.Timeout)*time.Second)
	return ctx
}
