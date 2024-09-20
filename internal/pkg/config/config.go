package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	port  string `mapstructure:"PORT"`
	dbUrl string `mapstructure:"DB_URL"`
}

func (c *Config) GetPort() string {
	return c.port
}

func (c *Config) GetUrl() string {
	return c.dbUrl
}

func InitConfig() (c Config, err error) {
	viper.AddConfigPath("./internal/pkg/config/envs")
	viper.SetConfigName("cfg")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	c.port = viper.Get("PORT").(string)
	c.dbUrl = viper.Get("DB_URL").(string)

	return
}
