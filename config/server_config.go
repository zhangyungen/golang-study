package config

import (
	"fmt"
	"time"
)

type ServerConfig struct {
	Name    string        `mapstructure:"name"`
	Version string        `mapstructure:"version"`
	Port    int           `mapstructure:"port"`
	Debug   bool          `mapstructure:"debug"`
	Timeout time.Duration `mapstructure:"timeout"`
}

func (a *ServerConfig) GetAddr() string {
	return fmt.Sprintf(":%d", a.Port)
}
