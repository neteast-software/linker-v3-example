package config

import (
	"os"
	"strconv"
	"time"

	postgresql "github.com/neteast-software/go-module/db/postgresql"
	http "github.com/neteast-software/go-module/http/gin"
	rpcgrpc "github.com/neteast-software/go-module/rpc/grpc"
)

const ExampleLoginPassword = "linfunlinfun"
const ExampleUserPhone = "18558755877"

type Config struct {
	HTTP            http.Config          `json:"http" yaml:"http"`
	GRPC            rpcgrpc.ServerConfig `json:"grpc" yaml:"grpc"`
	TTSClient       rpcgrpc.ClientConfig `json:"ttsClient" yaml:"ttsClient"`
	PostgreSQL      postgresql.Config    `json:"postgresql" yaml:"postgresql"`
	ShutdownTimeout time.Duration        `json:"shutdownTimeout" yaml:"shutdownTimeout"`
}

func Default() Config {
	return Config{
		HTTP: http.Config{
			Addr:     "127.0.0.1:8800",
			Recovery: true,
			Health:   http.HealthConfig{Enabled: true, Path: "ping"},
		},
		GRPC: rpcgrpc.ServerConfig{
			Addr: "127.0.0.1:9900",
		},
		TTSClient: rpcgrpc.ClientConfig{
			Discovery: rpcgrpc.ClientDiscoveryConfig{
				Scheme:  "example",
				Service: "tts",
				Group:   "DEFAULT_GROUP",
				Metadata: map[string]string{
					"version": "v1",
				},
			},
			Timeout: time.Second,
			Metadata: map[string]string{
				"scope": "front",
			},
		},
		PostgreSQL: postgresql.Config{
			Host:    "192.168.3.13",
			Port:    5432,
			User:    "neteast",
			DBName:  "linker_v3_example",
			SSLMode: "disable",
		},
		ShutdownTimeout: 3 * time.Second,
	}
}

func FromEnv() Config {
	config := Default()
	ApplyEnv(&config)
	return config
}

func ApplyEnv(config *Config) {
	if config == nil {
		return
	}
	config.HTTP.Addr = env("LINKER_V3_EXAMPLE_HTTP_ADDR", config.HTTP.Addr)
	config.GRPC.Addr = env("LINKER_V3_EXAMPLE_GRPC_ADDR", config.GRPC.Addr)
	config.PostgreSQL.Host = env("LINKER_V3_EXAMPLE_PG_HOST", config.PostgreSQL.Host)
	config.PostgreSQL.Port = envInt("LINKER_V3_EXAMPLE_PG_PORT", config.PostgreSQL.Port)
	config.PostgreSQL.User = env("LINKER_V3_EXAMPLE_PG_USER", config.PostgreSQL.User)
	config.PostgreSQL.Password = env("LINKER_V3_EXAMPLE_PG_PASSWORD", config.PostgreSQL.Password)
	config.PostgreSQL.DBName = env("LINKER_V3_EXAMPLE_PG_DB", config.PostgreSQL.DBName)
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
