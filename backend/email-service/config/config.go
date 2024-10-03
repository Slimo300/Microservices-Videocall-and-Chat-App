package config

import (
	"errors"
	"os"
	"strconv"
)

// Config holds email service configuration
type Config struct {
	BrokerAddress string

	Origin string

	EmailFrom string
	SMTPHost  string
	SMTPPort  int
	SMTPUser  string
	SMTPPass  string
}

// LoadConfigFromEnvironment populates Config fields with environment variables
// throws an error if any environment variable is not found
func LoadConfigFromEnvironment() (conf Config, err error) {
	conf.BrokerAddress = os.Getenv("BROKER_ADDRESS")
	if len(conf.BrokerAddress) == 0 {
		return Config{}, errors.New("environment variable BROKER_ADDRESS not found")
	}

	conf.Origin = os.Getenv("ORIGIN")
	if len(conf.Origin) == 0 {
		return Config{}, errors.New("environment variable ORIGIN not found")
	}

	conf.EmailFrom = os.Getenv("EMAIL_FROM")
	if len(conf.EmailFrom) == 0 {
		return Config{}, errors.New("environment variable EMAIL_FROM not found")
	}

	conf.SMTPHost = os.Getenv("SMTP_HOST")
	if len(conf.SMTPHost) == 0 {
		return Config{}, errors.New("environment variable SMTP_HOST not found")
	}

	port := os.Getenv("SMTP_PORT")
	if len(port) == 0 {
		return Config{}, errors.New("environment variable SMTP_PORT not found")
	}
	conf.SMTPPort, err = strconv.Atoi(port)
	if err != nil {
		return Config{}, err
	}

	conf.SMTPUser = os.Getenv("SMTP_USER")
	if len(conf.SMTPUser) == 0 {
		return Config{}, errors.New("environment variable SMTP_USER not found")
	}

	conf.SMTPPass = os.Getenv("SMTP_PASS")
	if len(conf.SMTPPass) == 0 {
		return Config{}, errors.New("environment variable SMTP_PASS not found")
	}

	return conf, nil
}
