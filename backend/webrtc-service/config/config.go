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
	PodName      string `mapstructure:"podName"`
	PodNamespace string `mapstructure:"podNamespace"`
	ServiceName  string `mapstructure:"serviceName"`

	DBAddress  string `mapstructure:"dbAddress"`
	DBPassword string `mapstructure:"dbPassword"`

	HTTPPort string `mapstructure:"httpPort"`

	Origin string `mapstructure:"origin"`

	BrokerAddress string `mapstructure:"brokerAddress"`
}

// LoadConfigFromEnvironment loads user service configuration from environment variables and returns an error
// if any of them is missing
func LoadConfigFromEnvironment() (conf Config, err error) {

	conf.PodName = os.Getenv("POD_NAME")
	if conf.PodName == "" {
		return Config{}, errors.New("Environment variable POD_NAME not set")
	}
	conf.PodNamespace = os.Getenv("POD_NAMESPACE")
	if conf.PodNamespace == "" {
		return Config{}, errors.New("Environment variable POD_NAMESPACE not set")
	}
	conf.ServiceName = os.Getenv("SERVICE_NAME")
	if conf.ServiceName == "" {
		return Config{}, errors.New("Environment variable SERVICE_NAME not set")
	}

	conf.DBAddress = os.Getenv("DB_ADDRESS")
	if conf.DBAddress == "" {
		return Config{}, errors.New("Environment variable DB_PASSWORD not set")
	}
	conf.DBPassword = os.Getenv("DB_PASSWORD")
	if conf.DBPassword == "" {
		return Config{}, errors.New("Environment variable DB_PASSWORD not set")
	}

	conf.HTTPPort = os.Getenv("HTTP_PORT")
	if conf.HTTPPort == "" {
		return Config{}, errors.New("Environment variable HTTP_PORT not set")
	}

	conf.Origin = os.Getenv("ORIGIN")
	if conf.Origin == "" {
		return Config{}, errors.New("Environment variable ORIGIN not set")
	}

	conf.BrokerAddress = os.Getenv("BROKER_ADDRESS")
	if conf.BrokerAddress == "" {
		return Config{}, errors.New("Environment variable BROKER_ADDRESS not set")
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
