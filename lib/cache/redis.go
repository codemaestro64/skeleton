package cache

import (
	"fmt"
	"time"
	"context"

	"github.com/codemaestro64/skeleton/config"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client redis.UniversalClient
	config config.RedisConfig 
}

const redis_prefix = "REDIS"

func NewRedis(cfg *config.Config) (Store, error) {
	r := &Redis{
		config: cfg.Redis,
	}

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: cfg.Redis.Addresses,
	})

	if err := rdb.Ping(r.getTimeoutContext()).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis cluster: %s", err.Error())
	}
	r.client = rdb

	return r, nil
}

func (r *Redis) getTimeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.config.Timeout)*time.Second)
	return ctx
}


func (r *Redis) Put(key string, data interface{}, duration time.Duration) error {
	_, err := r.client.Set(r.getTimeoutContext(), key, data, duration).Result()
	if err != nil {
		return fmt.Errorf("%s: error putting item: %s", redis_prefix, err.Error())
	} 

	return nil
}

/**func (r *Redis) Add(key string, data interface{}) bool {
	return false
}**/

func (r *Redis) Has(key string) (bool, error) {
	exists, err := r.client.Exists(r.getTimeoutContext(), key).Result()
	if err != nil {
		return false, fmt.Errorf("%s: error checking if item exists: %s", redis_prefix, err.Error())
	}

	return exists > 0, nil
}

func (r *Redis) Get(key string) (interface{}, error) {
	val, err := r.client.Get(r.getTimeoutContext(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: error getting item: %s", redis_prefix, err.Error())
	}

	return val, nil
}

/**func (r *Redis) GetOrDefault(key string, data interface{}) interface{} {
	return false
}**/

/**func (r *Redis) Remember(key string, duration int64, cb func() interface{}) {
	
}**/

/**func (r *Redis) RememberForever(key string, cb func() interface{}) {
	
}**/

/**func (r *Redis) Pull(key string) interface{} {
	return false
}**/

func (r *Redis) Remove(key string) error {
	_, err := r.client.Del(r.getTimeoutContext(), key).Result()
	if err != nil {
		return fmt.Errorf("%s: error deleting item: %s", redis_prefix, err.Error())
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