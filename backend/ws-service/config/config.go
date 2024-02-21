package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds user service configuration
type Config struct {
	PodName string `mapstructure:"podName"`

	DBAddress  string `mapstructure:"dbAddress"`
	DBPassword string `mapstructure:"dbPassword"`

	HTTPPort string `mapstructure:"httpPort"`

	Origin string `mapstructure:"origin"`

	BrokerType    string `mapstructure:"brokerType"`
	BrokerAddress string `mapstructure:"brokerAddress"`
}

// LoadConfigFromEnvironment loads user service configuration from environment variables and returns an error
// if any of them is missing
func LoadConfigFromEnvironment() (conf Config, err error) {

	conf.PodName = os.Getenv("POD_NAME")
	if len(conf.PodName) == 0 {
		return Config{}, errors.New("environment variable POD_NAME not set")
	}

	conf.DBAddress = os.Getenv("REDIS_ADDRESS")
	if len(conf.DBAddress) == 0 {
		return Config{}, errors.New("environment variable REDIS_PASSWORD not set")
	}
	conf.DBPassword = os.Getenv("REDIS_PASSWORD")
	if len(conf.DBPassword) == 0 {
		return Config{}, errors.New("environment variable REDIS_PASSWORD not set")
	}

	conf.HTTPPort = os.Getenv("HTTP_PORT")
	if len(conf.HTTPPort) == 0 {
		return Config{}, errors.New("environment variable HTTP_PORT not set")
	}

	conf.Origin = os.Getenv("ORIGIN")
	if len(conf.Origin) == 0 {
		return Config{}, errors.New("environment variable ORIGIN not set")
	}

	conf.BrokerType = os.Getenv("BROKER_TYPE")
	if len(conf.BrokerType) == 0 {
		return Config{}, errors.New("environment variable BROKER_TYPE not set")
	}

	conf.BrokerAddress = os.Getenv("BROKER_ADDRESS")
	if len(conf.BrokerAddress) == 0 {
		return Config{}, errors.New("environment variable BROKER_ADDRESS not set")
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
