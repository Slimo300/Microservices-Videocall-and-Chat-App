package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// DB is a representation of redis database
type DB struct {
	*redis.Client
	AccessCodeTTL time.Duration
}

// DBOption is a function whose only goal is modifying DB instance
type DBOption func(db *DB) *DB

// WithTTL takes DB and sets its TTL to given value
func WithAccessCodeTTL(ttl time.Duration) DBOption {
	return func(db *DB) *DB {
		db.AccessCodeTTL = ttl
		return db
	}
}

// Setup creates a new redis database client with given parameters
func Setup(addr, password string, options ...DBOption) (*DB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	db := &DB{
		Client:        client,
		AccessCodeTTL: 10 * time.Second,
	}

	for _, option := range options {
		db = option(db)
	}

	return db, nil
}
