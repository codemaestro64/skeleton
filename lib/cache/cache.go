package cache

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"

	"github.com/codemaestro64/skeleton/config"
)

type Store interface {
	Put(key string, data interface{}, duration time.Duration) error
	// Check if item exists in cache
	Has(key string) (bool, error)
	// Get an item from the cache
	Get(key string) (interface{}, error)
	// Remove item from cache
	Remove(key string) error
	// Remove all items from cache
	Flush()
	// Close
	Close()
}

type Cache struct {
	store Store
}

var (
	driverFuncMap = map[string]func(cfg *config.Config) (Store, error){
		"redis": NewRedis,
	}
)

func New(cfg *config.Config) (*Cache, error) {
	c := &Cache{}

	driverFunc, ok := driverFuncMap[cfg.Cache.Driver]
	if ok {
		var err error
		if c.store, err = driverFunc(cfg); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Cache) Put(key string, data interface{}, duration int64) error {
	return c.store.Put(key, data, time.Duration(duration)*time.Second)
}

func (c *Cache) Has(key string) (bool, error) {
	return c.store.Has(key)
}

func (c *Cache) Add(key string, data interface{}, duration int64) error {
	has, err := c.Has(key)
	if err != nil {
		return err
	}

	if has {
		return nil
	}

	return c.Put(key, data, duration)
}

func (c *Cache) Get(key string) (interface{}, error) {
	res, err := c.store.Get(key)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Cache) GetOrDefault(key string, def interface{}) (interface{}, error) {
	has, err := c.Has(key)
	if err != nil {
		return nil, err
	}

	if !has {
		return def, nil
	}

	return c.Get(key)
}

func (c *Cache) Remember(key string, duration int64, cb func() interface{}) error {
	has, err := c.Has(key)
	if err != nil {
		return err
	}

	var data interface{}
	if !has {
		data = cb()
	}

	return c.Put(key, data, duration)
}

func (c *Cache) RememberForever(key string, cb func() interface{}) error {
	has, err := c.Has(key)
	if err != nil {
		return err
	}

	var data interface{}
	if !has {
		data = cb()
	}

	return c.Put(key, data, 86400*365)
}

func (c *Cache) Remove(key string) error {
	return c.store.Remove(key)
}

func (c *Cache) Flush() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		c.store.Flush()
		wg.Done()
	}()

	wg.Wait()
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
