package redis

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/token-service/database"
	"github.com/go-redis/redis"
)

type redisTokenDB struct {
	*redis.Client
}

func NewRedisTokenDB(address, password string) (database.TokenDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &redisTokenDB{
		Client: client,
	}, nil
}

func (rdb *redisTokenDB) SaveToken(token string, expiration time.Duration) error {
	return rdb.Set(token, "1", expiration).Err()
}

func (rdb *redisTokenDB) IsTokenValid(userID, tokenID string) (bool, error) {
	pattern := fmt.Sprintf("%s:*%s", userID, tokenID)

	keys, err := rdb.Keys(pattern).Result()
	if err != nil {
		return false, err
	}
	if len(keys) == 0 {
		return false, database.TokenNotFoundError
	}

	res, err := rdb.Get(keys[0]).Result()
	if err != nil {
		return false, err
	}
	if database.StringToTokenValue(res) != database.TOKEN_VALID {
		if database.StringToTokenValue(res) == database.TOKEN_BLACKLISTED {
			return false, database.TokenBlacklistedError
		}
		return false, errors.New("Unexpected token value")
	}

	return true, nil
}

func (rdb *redisTokenDB) InvalidateTokens(userID, tokenID string) error {
	t := tokenID
	for {
		key := fmt.Sprintf("%s:%s:*", userID, t)

		keys, err := rdb.Keys(key).Result()
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			break
		}
		if err := rdb.Do("set", keys[0], string(database.TOKEN_BLACKLISTED), "keepttl").Err(); err != nil {
			return err
		}
		t = strings.Split(keys[0], ":")[2]
	}

	return nil
}

func (rdb *redisTokenDB) InvalidateToken(userID, tokenID string) error {
	key := fmt.Sprintf("%s:*%s", userID, tokenID)

	keys, err := rdb.Keys(key).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return database.TokenNotFoundError
	}

	if err := rdb.Do("set", keys[0], string(database.TOKEN_BLACKLISTED), "keepttl").Err(); err != nil {
		return err
	}
	return nil
}
