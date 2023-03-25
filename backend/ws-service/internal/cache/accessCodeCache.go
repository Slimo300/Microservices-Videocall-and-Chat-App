package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

const ACCESS_CODE_TTL = 100

type AccessCodeCache struct {
	codes map[string]User
	stop  chan struct{}
	mu    *sync.RWMutex
}

type User struct {
	ID       uuid.UUID
	Deadline time.Time
}

func NewCache(cleanupInterval time.Duration) AccessCodeCache {

	cache := AccessCodeCache{
		stop:  make(chan struct{}),
		codes: make(map[string]User),
		mu:    &sync.RWMutex{},
	}

	go func(cleanupInterval time.Duration) {
		cache.cleanupLoop(cleanupInterval * time.Second)
	}(cleanupInterval)

	return cache
}

func (c *AccessCodeCache) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-t.C:
			c.mu.Lock()
			for code, user := range c.codes {
				if user.Deadline.Unix() <= time.Now().Unix() {
					delete(c.codes, code)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *AccessCodeCache) Stop() {
	close(c.stop)
}

func (c *AccessCodeCache) Set(key string, id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.codes[key] = User{
		Deadline: time.Now().Add(ACCESS_CODE_TTL * time.Second),
		ID:       id,
	}
}

func (c *AccessCodeCache) Read(key string) (User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	user, ok := c.codes[key]
	if !ok {
		return User{}, fmt.Errorf("code %s not in cache", key)
	}

	return user, nil
}

func (c *AccessCodeCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.codes, key)
}
