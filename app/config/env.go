package configs

import (
	"ratelimiter/handlers"

	"github.com/spf13/viper"
	// Importe o pacote handlers ou o pacote onde RateLimitConfig está definido
)

func LoadConfig(path string) (*handlers.RateLimitConfig, error) {
	var cfg *handlers.RateLimitConfig
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}


