package redis

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/token-service/repo"
	"github.com/go-redis/redis"
)

type redisTokenRepository struct {
	*redis.Client
}

func NewRedisTokenRepository(address, password string) (*redisTokenRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &redisTokenRepository{
		Client: client,
	}, nil
}

func (rdb *redisTokenRepository) SaveToken(token string, expiration time.Duration) error {
	return rdb.Set(token, "1", expiration).Err()
}

func (rdb *redisTokenRepository) IsTokenValid(userID, tokenID string) (bool, error) {
	pattern := fmt.Sprintf("%s:*%s", userID, tokenID)
	log.Println(pattern)

	keys, err := rdb.Keys(pattern).Result()
	if err != nil {
		return false, err
	}
	if len(keys) == 0 {
		return false, repo.TokenNotFoundError
	}
	if len(keys) > 1 {
		return false, repo.TooManyTokensFoundError
	}

	res, err := rdb.Get(keys[0]).Result()
	if err != nil {
		return false, err
	}
	if repo.StringToTokenValue(res) != repo.TOKEN_VALID {
		if repo.StringToTokenValue(res) == repo.TOKEN_BLACKLISTED {
			return false, repo.TokenBlacklistedError
		}
		return false, errors.New("Unexpected token value")
	}

	return true, nil
}

func (rdb *redisTokenRepository) InvalidateTokens(userID, tokenID string) error {
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
		if len(keys) > 1 {
			return repo.TooManyTokensFoundError
		}
		if err := rdb.Do("set", keys[0], string(repo.TOKEN_BLACKLISTED), "keepttl").Err(); err != nil {
			return err
		}
		t = strings.Split(keys[0], ":")[2]
	}

	return nil
}

func (rdb *redisTokenRepository) InvalidateToken(userID, tokenID string) error {
	key := fmt.Sprintf("%s:*%s", userID, tokenID)

	keys, err := rdb.Keys(key).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return repo.TokenNotFoundError
	}
	if len(keys) > 1 {
		return repo.TooManyTokensFoundError
	}

	if err := rdb.Do("set", keys[0], string(repo.TOKEN_BLACKLISTED), "keepttl").Err(); err != nil {
		return err
	}
	return nil
}
