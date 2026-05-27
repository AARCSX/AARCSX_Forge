package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Environment string

const (
	EnvLocal   Environment = "local"
	EnvStaging Environment = "staging"
	EnvProd    Environment = "production"
)

type Config struct {
	Env           Environment
	API           APIConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	Storage       StorageConfig
	JWT           JWTConfig
	Queue         QueueConfig
	Observability ObservabilityConfig
}

type APIConfig struct {
	HTTPAddress string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type StorageConfig struct {
	Provider string
	MinIO    MinIOConfig
	S3       S3Config
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type S3Config struct {
	Endpoint string
	Region   string
	Bucket   string
}

type JWTConfig struct {
	Issuer         string
	Audience       string
	SigningKey     string
	KeyID          string
	AccessTTLMin   int
	RefreshTTLHour int
}

type QueueConfig struct {
	Concurrency int
}

type ObservabilityConfig struct {
	ServiceName string
	LogLevel    string
}

func Load() (Config, error) {
	cfg := Config{
		Env: Environment(getEnv("APP_ENV", string(EnvLocal))),
		API: APIConfig{
			HTTPAddress: getEnv("HTTP_ADDR", ":8080"),
		},
		Database: DatabaseConfig{
			URL: strings.TrimSpace(getEnv("POSTGRES_URL", "")),
		},
		Redis: RedisConfig{
			URL: strings.TrimSpace(getEnv("REDIS_URL", "")),
		},
		Storage: StorageConfig{
			Provider: strings.ToLower(strings.TrimSpace(getEnv("STORAGE_PROVIDER", "minio"))),
			MinIO: MinIOConfig{
				Endpoint:  getEnv("MINIO_ENDPOINT", ""),
				AccessKey: getEnv("MINIO_ACCESS_KEY", ""),
				SecretKey: getEnv("MINIO_SECRET_KEY", ""),
				Bucket:    getEnv("MINIO_BUCKET", ""),
				UseSSL:    strings.EqualFold(getEnv("MINIO_USE_SSL", "false"), "true"),
			},
			S3: S3Config{
				Endpoint: getEnv("S3_ENDPOINT", ""),
				Region:   getEnv("S3_REGION", ""),
				Bucket:   getEnv("S3_BUCKET", ""),
			},
		},
		JWT: JWTConfig{
			Issuer:         getEnv("JWT_ISSUER", "aarcsx-forge"),
			Audience:       getEnv("JWT_AUDIENCE", "aarcsx-internal"),
			SigningKey:     getEnv("JWT_SIGNING_KEY", ""),
			KeyID:          getEnv("JWT_KID", "default"),
			AccessTTLMin:   getIntEnv("JWT_ACCESS_TTL_MIN", 15),
			RefreshTTLHour: getIntEnv("JWT_REFRESH_TTL_HOUR", 720),
		},
		Queue: QueueConfig{
			Concurrency: getIntEnv("QUEUE_CONCURRENCY", 10),
		},
		Observability: ObservabilityConfig{
			ServiceName: getEnv("OTEL_SERVICE_NAME", "aarcsx-forge"),
			LogLevel:    strings.ToLower(getEnv("LOG_LEVEL", "info")),
		},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	var errs []error

	if c.Database.URL == "" {
		errs = append(errs, errors.New("POSTGRES_URL is required"))
	}
	if c.Redis.URL == "" {
		errs = append(errs, errors.New("REDIS_URL is required"))
	}
	if c.JWT.SigningKey == "" {
		errs = append(errs, errors.New("JWT_SIGNING_KEY is required"))
	}
	if c.Storage.Provider == "" {
		errs = append(errs, errors.New("STORAGE_PROVIDER is required"))
	}
	switch c.Storage.Provider {
	case "minio":
		if c.Storage.MinIO.Endpoint == "" || c.Storage.MinIO.Bucket == "" {
			errs = append(errs, errors.New("MINIO_ENDPOINT and MINIO_BUCKET are required for minio provider"))
		}
	case "s3":
		if c.Storage.S3.Region == "" || c.Storage.S3.Bucket == "" {
			errs = append(errs, errors.New("S3_REGION and S3_BUCKET are required for s3 provider"))
		}
	default:
		errs = append(errs, fmt.Errorf("unsupported STORAGE_PROVIDER %q", c.Storage.Provider))
	}
	return errors.Join(errs...)
}

func getEnv(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	raw := strings.TrimSpace(getEnv(key, ""))
	if raw == "" {
		return defaultValue
	}
	var out int
	_, err := fmt.Sscanf(raw, "%d", &out)
	if err != nil {
		return defaultValue
	}
	return out
}
