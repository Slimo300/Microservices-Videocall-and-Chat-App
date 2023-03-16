package main

import (
	"os"
	"time"
)

type Config struct {
	GRPCPort string

	RedisAddress  string
	RedisPassword string

	AccessDuration     time.Duration
	RefreshDuration    time.Duration
	RefreshTokenSecret string
}

func loadConfig() (conf Config, err error) {

	conf.GRPCPort = os.Getenv("GRPC_PORT")
	conf.RedisAddress = os.Getenv("REDIS_ADDRESS")
	conf.RedisPassword = os.Getenv("REDIS_PASSWORD")
	conf.RefreshTokenSecret = os.Getenv("REFRESH_SECRET")

	conf.RefreshDuration, err = time.ParseDuration(os.Getenv("REFRESH_DURATION"))
	if err != nil {
		return Config{}, err
	}
	conf.AccessDuration, err = time.ParseDuration(os.Getenv("ACCESS_DURATION"))
	if err != nil {
		return Config{}, err
	}

	return
}
