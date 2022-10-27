package configuration

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GroupService   Service `mapstructure:"group-service"`
	UserService    Service `mapstructure:"user-service"`
	MessageService Service `mapstructure:"message-service"`
	WSService      Service `mapstructure:"ws-service"`

	TokenService TokenService `mapstructure:"token-service"`

	S3Bucket string `mapstructure:"aws-bucket"`

	Certificate string `mapstructure:"cert"`
	PrivKeyFile string `mapstructure:"privKey"`
}

type Service struct {
	DBType    string `mapstructure:"dbtype"`
	DBAddress string `mapstructure:"dbaddress"`
	HTTPPort  string `mapstructure:"httpPort"`
	HTTPSPort string `mapstructure:"httpsPort"`
}

type TokenService struct {
	GRPCPort              string        `mapstructure:"grpcPort"`
	RedisAddress          string        `mapstructure:"redisAddress"`
	RedisPass             string        `mapstructure:"redisPass"`
	AccessTokenPrivateKey string        `mapstructure:"accessPrivKey"`
	RefreshTokenSecret    string        `mapstructure:"refreshSecret"`
	AccessDuration        time.Duration `mapstructure:"accessDuration"`
	RefreshDuration       time.Duration `mapstructure:"refreshDuration"`
}

func LoadConfig(path string) (config Config, err error) {
	vp := viper.New()
	vp.AddConfigPath(path)
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")

	err = vp.ReadInConfig()
	if err != nil {
		return
	}

	err = vp.Unmarshal(&config)
	return
}
