package configuration

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds information about project resources and behavior
type Config struct {
	Origin string `mapstructure:"origin"`
	Domain string `mapstructure:"domain"`

	GroupService   Service `mapstructure:"group-service"`
	UserService    Service `mapstructure:"user-service"`
	MessageService Service `mapstructure:"message-service"`
	WSService      Service `mapstructure:"ws-service"`

	SearchService SearchService `mapstructure:"search-service"`
	TokenService  TokenService  `mapstructure:"token-service"`

	S3Bucket string `mapstructure:"aws-bucket"`

	Certificate string `mapstructure:"cert"`
	PrivKeyFile string `mapstructure:"privKey"`

	BrokerType       string   `mapstructure:"brokerType"`
	BrokersAddresses []string `mapstructure:"brokerAddresses"`

	EmailFrom        string `mapstructure:"emailFrom"`
	SMTPHost         string `mapstructure:"smtpHost"`
	SMTPPort         int    `mapstructure:"smtpPort"`
	SMTPUser         string `mapstructure:"smtpUser"`
	SMTPPass         string `mapstructure:"smtpPass"`
	EmailTemplateDir string `mapstructure:"templateDir"`

	AuthAddress     string        `mapstructure:"authAddress"`
	AccessDuration  time.Duration `mapstructure:"accessDuration"`
	RefreshDuration time.Duration `mapstructure:"refreshDuration"`
}

// Service type holds regular information for REST service
type Service struct {
	DBType    string `mapstructure:"dbtype"`
	DBAddress string `mapstructure:"dbaddress"`
	HTTPPort  string `mapstructure:"httpPort"`
	HTTPSPort string `mapstructure:"httpsPort"`
}

// TokenService is a specific type for managing token-service
type TokenService struct {
	GRPCPort              string `mapstructure:"grpcPort"`
	RedisAddress          string `mapstructure:"redisAddress"`
	RedisPass             string `mapstructure:"redisPass"`
	AccessTokenPrivateKey string `mapstructure:"accessPrivKey"`
	RefreshTokenSecret    string `mapstructure:"refreshSecret"`
}

// SearchService is a specific type for managing search-service
type SearchService struct {
	HTTPPort  string   `mapstructure:"httpPort"`
	HTTPSPort string   `mapstructure:"httpsPort"`
	Addresses []string `mapstructure:"addresses"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstrucutre:"password"`
}

// LoadConfig looks up config.yaml file in path given by parameter and reads
// project configuration from it
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
