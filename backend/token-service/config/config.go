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
	GRPCPort string `mapstructure:"grpcPort"`

	RedisAddress  string `mapstructure:"redisAddress"`
	RedisPassword string `mapstructure:"redisPassword"`

	AccessDuration     time.Duration `mapstructure:"accessDuration"`
	RefreshDuration    time.Duration `mapstructure:"refreshDuration"`
	RefreshTokenSecret string        `mapstructure:"secret"`
}

// LoadConfigFromEnvironment loads token service from environment variables and returns an error
// if any of them is missing
func LoadConfigFromEnvironment() (conf Config, err error) {

	conf.GRPCPort = os.Getenv("GRPC_PORT")
	if len(conf.GRPCPort) == 0 {
		return Config{}, errors.New("Environment variable GRPC_PORT not set")
	}

	conf.RedisAddress = os.Getenv("REDIS_ADDRESS")
	if len(conf.RedisAddress) == 0 {
		return Config{}, errors.New("Environment variable GRPC_PORT not set")
	}

	conf.RedisPassword = os.Getenv("REDIS_PASSWORD")
	if len(conf.RedisPassword) == 0 {
		return Config{}, errors.New("Environment variable GRPC_PORT not set")
	}

	conf.RefreshTokenSecret = os.Getenv("REFRESH_SECRET")
	if len(conf.RefreshTokenSecret) == 0 {
		return Config{}, errors.New("Environment variable GRPC_PORT not set")
	}

	refreshDuration := os.Getenv("REFRESH_DURATION")
	if len(refreshDuration) == 0 {
		return Config{}, errors.New("Environment variable REFRESH_DURATION not set")
	}
	conf.RefreshDuration, err = time.ParseDuration(refreshDuration)
	if err != nil {
		return Config{}, err
	}

	accessDuration := os.Getenv("ACCESS_DURATION")
	if len(accessDuration) == 0 {
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
