package redis

import (
	"github.com/go-redis/redis"
)

type DB struct {
	*redis.Client
}

// Setup creates a new redis database client with given parameters
func Setup(addr, password string) (*DB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	db := &DB{
		Client: client,
	}

	return db, nil
}
