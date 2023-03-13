package cfg

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Grpc string `mapstructure:"GRPC_PORT"`
	Http string `mapstructure:"HTTP_PORT"`

	TokenSecret  string `mapstructure:"TOKEN_SECRET"`
	TokenBuilder string `mapstructure:"TOKEN_BUILDER"`

	DefaultTokenDuration time.Duration `mapstructure:"DEFAULT_TOKEN_DURATION"`
	MaximumTokenDuration time.Duration `mapstructure:"MAXIMUM_TOKEN_DURATION"`

	ApiHeader string `mapstructure:"API_HEADER"`
	ApiSecret string `mapstructure:"API_SECRET"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("example")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
