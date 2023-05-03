package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	DbFile       = "MY_RECON_DB"
	DomainConfig = "MY_RECON_DOMAINS"
)

func New() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	v.SetConfigFile(".config.yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "reading config")
	}

	v.AutomaticEnv()
	return v, nil
}
