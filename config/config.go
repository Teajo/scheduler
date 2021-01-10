package config

import (
	"os"
	"strconv"
	"time"
)

// QueueConfig is queue config
type QueueConfig struct {
	MaxQueueLen int
	TimeChunk   time.Duration
}

// HTTPApiConfig is http api config
type HTTPApiConfig struct {
	Port int
}

// Config is scheduler config
type Config struct {
	PluginDir string
	DbName    string
	DbDriver  string
	*QueueConfig
	*HTTPApiConfig
}

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

	ms, err := strconv.Atoi(getenv("QUEUE_TIME_CHUNK", "21600000")) // default 6hours
	if err != nil {
		panic(err)
	}

	pluginDir := getenv("PLUGIN_DIR", "./plugins")
	dbDriver := getenv("DB_DRIVER", "sqlite3")
	dbName := getenv("DB_NAME", "scheduler")

	return &Config{
		DbName:    dbName,
		DbDriver:  dbDriver,
		PluginDir: pluginDir,
		QueueConfig: &QueueConfig{
			MaxQueueLen: maxQueueLen,
			TimeChunk:   time.Duration(ms) * time.Millisecond,
		},
		HTTPApiConfig: &HTTPApiConfig{
			Port: port,
		},
	}
}

func getenv(envar string, dflt string) string {
	value := os.Getenv(envar)
	if value == "" {
		value = dflt
	}
	return value
}
