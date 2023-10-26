package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds user service configuration
type Config struct {
	DBAddress string `mapstructure:"dbAddress"`
	HTTPPort  string `mapstructure:"httpPort"`

	Origin string `mapstructure:"origin"`

	BrokerAddress string `mapstructure:"brokerAddress"`

	StorageKeyID     string `mapstructure:"storageKeyID"`
	StorageKeySecret string `mapstructure:"storageKeySecret"`
	Bucket           string `mapstructure:"bucketname"`
}

// LoadConfigFromEnvironment loads user service configuration from environment variables and returns an error
// if any of them is missing
func LoadConfigFromEnvironment() (conf Config, err error) {

	mySQLAddress := os.Getenv("MYSQL_ADDRESS")
	if len(mySQLAddress) == 0 {
		return Config{}, errors.New("Environment variable MYSQL_ADDRESS not set")
	}
	mySQLDatabase := os.Getenv("MYSQL_DATABASE")
	if len(mySQLDatabase) == 0 {
		return Config{}, errors.New("Environment variable MYSQL_DATABASE not set")
	}
	mySQLUser := os.Getenv("MYSQL_USER")
	if len(mySQLUser) == 0 {
		return Config{}, errors.New("Environment variable MYSQL_USER not set")
	}
	mySQLPassword := os.Getenv("MYSQL_PASSWORD")
	if len(mySQLPassword) == 0 {
		return Config{}, errors.New("Environment variable MYSQL_PASSWORD not set")
	}

	conf.DBAddress = fmt.Sprintf("%s:%s@tcp(%s)/%s", mySQLUser, mySQLPassword, mySQLAddress, mySQLDatabase)

	conf.HTTPPort = os.Getenv("HTTP_PORT")
	if len(conf.HTTPPort) == 0 {
		return Config{}, errors.New("Environment variable HTTP_PORT not set")
	}

	conf.Origin = os.Getenv("ORIGIN")
	if len(conf.Origin) == 0 {
		return Config{}, errors.New("Environment variable ORIGIN not set")
	}

	conf.BrokerAddress = os.Getenv("BROKER_ADDRESS")
	if len(conf.BrokerAddress) == 0 {
		return Config{}, errors.New("Environment variable BROKER_ADDRESS not set")
	}

	conf.StorageKeyID = os.Getenv("STORAGE_KEY_ID")
	if len(conf.StorageKeyID) == 0 {
		return Config{}, errors.New("Environment variable STORAGE_KEY_ID not set")
	}
	conf.StorageKeySecret = os.Getenv("STORAGE_KEY_SECRET")
	if len(conf.StorageKeySecret) == 0 {
		return Config{}, errors.New("Environment variable STORAGE_KEY_SECRET not set")
	}
	conf.Bucket = os.Getenv("STORAGE_BUCKET")
	if len(conf.Bucket) == 0 {
		return Config{}, errors.New("Environment variable STORAGE_BUCKET not set")
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
