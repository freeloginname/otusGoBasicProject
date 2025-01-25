package config

import "github.com/spf13/viper"

type Config struct {
	Port  string `mapstructure:"APP_HTTP_PORT"`
	DBUrl string `mapstructure:"DATABASE_URL"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath(".env")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
