package config

import (
	"os"
	"strconv"
)

// QueueConfig is queue config
type QueueConfig struct {
	MaxQueueLen int
}

// HTTPApiConfig is http api config
type HTTPApiConfig struct {
	Port int
}

// Config is scheduler config
type Config struct {
	*QueueConfig
	*HTTPApiConfig
}

var config *Config = nil

// Get config
func Get() *Config {
	port, err := strconv.Atoi(getenv("HTTP_API_PORT", "3000"))
	if err != nil {
		panic(err)
	}

	maxQueueLen, err := strconv.Atoi(getenv("MAX_QUEUE_LENGTH", "10"))
	if err != nil {
		panic(err)
	}

	if config == nil {
		config = &Config{
			QueueConfig: &QueueConfig{
				MaxQueueLen: maxQueueLen,
			},
			HTTPApiConfig: &HTTPApiConfig{
				Port: port,
			},
		}
	}

	return config
}

func getenv(envar string, dflt string) string {
	value := os.Getenv(envar)
	if value == "" {
		value = dflt
	}
	return value
}
