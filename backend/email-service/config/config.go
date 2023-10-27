package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config holds email service configuration
type Config struct {
	GRPCPort string `mapstructure:"grpcPort"`

	Origin string `mapstructure:"origin"`

	EmailFrom string `mapstructure:"emailFrom"`
	SMTPHost  string `mapstructure:"smtpHost"`
	SMTPPort  int    `mapstructure:"smtpPort"`
	SMTPUser  string `mapstructure:"smtpUser"`
	SMTPPass  string `mapstructure:"smtpPass"`
}

// LoadConfigFromEnvironment populates Config fields with environment variables
// throws an error if any environment variable is not found
func LoadConfigFromEnvironment() (conf Config, err error) {
	conf.GRPCPort = os.Getenv("GRPC_PORT")
	if len(conf.GRPCPort) == 0 {
		return Config{}, errors.New("Environment variable GRPC_PORT not found")
	}

	conf.Origin = os.Getenv("ORIGIN")
	if len(conf.Origin) == 0 {
		return Config{}, errors.New("Environment variable ORIGIN not found")
	}

	conf.EmailFrom = os.Getenv("EMAIL_FROM")
	if len(conf.EmailFrom) == 0 {
		return Config{}, errors.New("Environment variable EMAIL_FROM not found")
	}

	conf.SMTPHost = os.Getenv("SMTP_HOST")
	if len(conf.SMTPHost) == 0 {
		return Config{}, errors.New("Environment variable SMTP_HOST not found")
	}

	port := os.Getenv("SMTP_PORT")
	if len(port) == 0 {
		return Config{}, errors.New("Environment variable SMTP_PORT not found")
	}
	conf.SMTPPort, err = strconv.Atoi(port)
	if err != nil {
		return Config{}, err
	}

	conf.SMTPUser = os.Getenv("SMTP_USER")
	if len(conf.SMTPUser) == 0 {
		return Config{}, errors.New("Environment variable SMTP_USER not found")
	}

	conf.SMTPPass = os.Getenv("SMTP_PASS")
	if len(conf.SMTPPass) == 0 {
		return Config{}, errors.New("Environment variable SMTP_PASS not found")
	}

	return conf, nil
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
