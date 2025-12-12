package cmd

import (
	"fmt"

	"github.com/pelletier/go-toml"
)


type Config struct {
	PostfixLogFile string `toml:"post_fix_log_file"`
}

func GetSystemConfig() (Config, error) {
	cfg, err := toml.LoadFile("config.toml")
	if err != nil {
		return Config{}, fmt.Errorf("can not open system config file err: %v\n", err)

	}
	var c Config
	if err := cfg.Unmarshal(&c); err != nil {
		return Config{}, fmt.Errorf("error unmarshalling config toml file")
	}
	return c, nil
}
