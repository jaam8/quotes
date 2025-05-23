package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port uint16 `env:"PORT" yaml:"port"`
}

func New() (Config, error) {
	var config Config
	val, exists := os.LookupEnv("PORT")
	if !exists {
		val = "8080"
	}

	port, err := strconv.Atoi(val)
	if err != nil || port < 0 || port > 65535 {
		return config, fmt.Errorf("invalid PORT value: %v", err)
	}
	config.Port = uint16(port)
	return config, nil
}
