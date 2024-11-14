package config

import (
	"errors"
	"fmt"
	"os"
)

// Config holds user service configuration
type Config struct {
	DBAddress string `mapstructure:"dbAddress"`
	HTTPPort  string `mapstructure:"httpPort"`

	Origin string `mapstructure:"origin"`

	BrokerType    string `mapstructure:"brokerType"`
	BrokerAddress string `mapstructure:"brokerAddress"`

	StorageKeyID     string `mapstructure:"storageKeyID"`
	StorageKeySecret string `mapstructure:"storageKeySecret"`
	StorageRegion    string `mapstructure:"storageRegion"`
	Bucket           string `mapstructure:"bucketname"`
}

// LoadConfigFromEnvironment loads user service configuration from environment variables and returns an error
// if any of them is missing
func LoadConfigFromEnvironment() (conf Config, err error) {

	mySQLAddress := os.Getenv("MYSQL_ADDRESS")
	if len(mySQLAddress) == 0 {
		return Config{}, errors.New("environment variable MYSQL_ADDRESS not set")
	}
	mySQLDatabase := os.Getenv("MYSQL_DATABASE")
	if len(mySQLDatabase) == 0 {
		return Config{}, errors.New("environment variable MYSQL_DATABASE not set")
	}
	mySQLUser := os.Getenv("MYSQL_USER")
	if len(mySQLUser) == 0 {
		return Config{}, errors.New("environment variable MYSQL_USER not set")
	}
	mySQLPassword := os.Getenv("MYSQL_PASSWORD")
	if len(mySQLPassword) == 0 {
		return Config{}, errors.New("environment variable MYSQL_PASSWORD not set")
	}

	conf.DBAddress = fmt.Sprintf("%s:%s@tcp(%s)/%s", mySQLUser, mySQLPassword, mySQLAddress, mySQLDatabase)

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

	conf.StorageKeyID = os.Getenv("STORAGE_KEY_ID")
	if len(conf.StorageKeyID) == 0 {
		return Config{}, errors.New("environment variable STORAGE_KEY_ID not set")
	}
	conf.StorageKeySecret = os.Getenv("STORAGE_KEY_SECRET")
	if len(conf.StorageKeySecret) == 0 {
		return Config{}, errors.New("environment variable STORAGE_KEY_SECRET not set")
	}
	conf.Bucket = os.Getenv("STORAGE_BUCKET")
	if len(conf.Bucket) == 0 {
		return Config{}, errors.New("environment variable STORAGE_BUCKET not set")
	}

	return
}
