package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds token service configuration
type Config struct {
	GRPCPort string

	RedisAddress  string
	RedisPassword string

	AccessDuration     time.Duration
	RefreshDuration    time.Duration
	RefreshTokenSecret string
}

// LoadConfigFromEnvironment loads token service from environment variables and returns an error
// if any of them is missing
func LoadConfigFromEnvironment() (conf Config, err error) {

	conf.GRPCPort = os.Getenv("GRPC_PORT")
	if conf.GRPCPort == "" {
		return Config{}, errors.New("Environment variable GRPC_PORT not set")
	}

	conf.RedisAddress = os.Getenv("REDIS_ADDRESS")
	if conf.RedisAddress == "" {
		return Config{}, errors.New("Environment variable GRPC_PORT not set")
	}

	conf.RedisPassword = os.Getenv("REDIS_PASSWORD")
	if conf.RedisPassword == "" {
		return Config{}, errors.New("Environment variable GRPC_PORT not set")
	}

	conf.RefreshTokenSecret = os.Getenv("REFRESH_SECRET")
	if conf.RefreshTokenSecret == "" {
		return Config{}, errors.New("Environment variable GRPC_PORT not set")
	}

	refreshDuration := os.Getenv("REFRESH_DURATION")
	if refreshDuration == "" {
		return Config{}, errors.New("Environment variable REFRESH_DURATION not set")
	}
	conf.RefreshDuration, err = time.ParseDuration(refreshDuration)
	if err != nil {
		return Config{}, err
	}

	accessDuration := os.Getenv("ACCESS_DURATION")
	if accessDuration == "" {
		return Config{}, errors.New("Environment variable ACCESS_DURATION not set")
	}
	conf.AccessDuration, err = time.ParseDuration(accessDuration)
	if err != nil {
		return Config{}, err
	}

	return
}

// LoadConfigFromFile loads config from specified path
func LoadConfigFromFile(path string) (config Config, err error) {
	vp := viper.New()

	vp.AddConfigPath(filepath.Dir(path))

	filename := strings.Split(filepath.Base(path), ".")
	vp.SetConfigName(filename[0])
	vp.SetConfigType(filename[1])

	if err = vp.ReadInConfig(); err != nil {
		return
	}

	err = vp.Unmarshal(&config)
	return
}
