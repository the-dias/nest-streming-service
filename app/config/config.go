package config

import (
	"nats-service/app/config/database"
	"nats-service/app/config/http"
	"nats-service/app/config/nats"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseConfig database.DatabaseConfig
	NatsConfig     nats.NatsConfig
	HttpConfig     http.HttpConfig
}

func New() *Config {
	return &Config{
		NatsConfig: nats.NatsConfig{
			ConnectWait:        time.Second * time.Duration(getEnvAsInt("NATS_CONNECT_WAIT", 0)),
			PubAckWait:         time.Second * time.Duration(getEnvAsInt("NATS_PUB_ACK_WAIT", 0)),
			Interval:           getEnvAsInt("NATS_INTERVAL", 0),
			MaxOut:             getEnvAsInt("NATS_MAX_OUT", 0),
			MaxPubAcksInflight: getEnvAsInt("NATS_MAX_PUB_ACKS_INFLIGHT", 0),
			Url:                getEnv("NATS_URL", ""),
			ClusterID:          getEnv("NATS_CLUSTER_ID", ""),
		},
		HttpConfig: http.HttpConfig{
			Port:          getEnvAsInt("HTTP_PORT_SERVER", 0),
			PatternServer: getEnv("HTTP_PATTERN_SERVER", ""),
			PatternStatic: getEnv("HTTP_PATTERN_STATIC", ""),
			StaticDir:     getEnv("HTTP_STATIC_DIR", ""),
		},
		DatabaseConfig: database.DatabaseConfig{
			Host:         getEnv("POSTGRES_HOST", ""),
			Port:         getEnvAsInt("POSTGRES_PORT", 0),
			User:         getEnv("POSTGRES_USER", ""),
			Password:     getEnv("POSTGRES_PASSWORD", ""),
			DatabaseName: getEnv("POSTGRES_DB_NAME", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
