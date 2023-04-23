package config

import "github.com/spf13/viper"

const (
	DB_FILE       = "MY_RECON_DB"
	DOMAIN_CONFIG = "MY_RECON_DOMAINS"
)

var c *viper.Viper

func Get() *viper.Viper {
	if c != nil {
		return c
	}
	c = viper.New()

	c.SetConfigFile(".env")
	c.ReadInConfig()
	c.AutomaticEnv()
	return c
}
