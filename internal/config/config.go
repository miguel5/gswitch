package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	BranchPrefix string `mapstructure:"branch_prefix"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath("$HOME/.config/gswitch")

	v.SetDefault("branch_prefix", "feature/")

	_ = v.ReadInConfig()

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
