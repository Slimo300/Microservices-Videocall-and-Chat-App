package configuration

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Group-service
	GS_HTTPPort string `mapstructure:"GS_HTTP_PORT"`
	GS_DBHost   string `mapstructure:"GS_DBHOST"`
	GS_DBPort   string `mapstructure:"GS_DBPORT"`
	GS_DBUser   string `mapstructure:"GS_DBUSER"`
	GS_DBPass   string `mapstructure:"GS_DBPASS"`

	// Message-service
	MS_HTTPPort string `mapstructure:"MS_HTTP_PORT"`
	MS_DBHost   string `mapstructure:"MS_DBHOST"`
	MS_DBPort   string `mapstructure:"MS_DBPORT"`
	MS_DBUser   string `mapstructure:"MS_DBUSER"`
	MS_DBPass   string `mapstructure:"MS_DBPASS"`

	// User-service
	US_HTTPPort string `mapstructure:"US_HTTP_PORT"`
	US_DBHost   string `mapstructure:"US_DBHOST"`
	US_DBPort   string `mapstructure:"US_DBPORT"`
	US_DBUser   string `mapstructure:"US_DBUSER"`
	US_DBPass   string `mapstructure:"US_DBPASS"`

	//Token-service
	REDISHost             string        `mapstructure:"REDIS_HOST"`
	REDISPort             string        `mapstructure:"REDIS_PORT"`
	REDISPass             string        `mapstructure:"REDIS_PASS"`
	AccessTokenPrivateKey string        `mapstructure:"PRIV_KEY_FILE"`
	RefreshTokenSecret    string        `mapstructure:"REFRESH_SECRET"`
	AccessDuration        time.Duration `mapstructure:"ACCESS_DURATION"`
	RefreshDuration       time.Duration `mapstructure:"REFRESH_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	vp := viper.New()
	vp.AddConfigPath(path)
	vp.SetConfigType("env")
	vp.SetConfigName("")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
