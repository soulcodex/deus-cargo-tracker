package configs

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppServiceName string `env:"SERVICE_NAME"`
	AppEnv         string `env:"ENV"`
	AppVersion     string `env:"VERSION"`
}
type RedisConfig struct {
	RedisURL string `env:"URL" envDefault:"redis://localhost:6379"`
}

type PostgresConfig struct {
	PostgresUser   string `env:"USER" envDefault:"cargo_tracker_role"`
	PostgresPass   string `env:"PASSWORD" envDefault:"cargo_tracker_pass"`
	PostgresHost   string `env:"HOST" envDefault:"localhost"`
	PostgresPort   uint16 `env:"PORT" envDefault:"5432"`
	PostgresDB     string `env:"DB" envDefault:"cargo_tracker"`
	PostgresSchema string `env:"SCHEMA" envDefault:"cargo_tracker"`
	PostgresSSL    string `env:"SSL_MODE" envDefault:"disable"`
}

type DBMigrationsConfig struct {
	MigrationsPath  string `env:"PATH" envDefault:"./migrations"`
	MigrationsTable string `env:"TABLE_NAME" envDefault:"migrations"`
}

type HTTPConfig struct {
	HTTPHost         string `env:"HOST" envDefault:"0.0.0.0"`
	HTTPPort         int    `env:"PORT" envDefault:"8080"`
	HTTPReadTimeout  int    `env:"READ_TIMEOUT" envDefault:"30"`
	HTTPWriteTimeout int    `env:"WRITE_TIMEOUT" envDefault:"30"`
}

type UncategorizedConfig struct {
	JSONSchemaPath string `env:"JSON_SCHEMA_PATH" envDefault:"./schemas"`
	LogLevel       string `env:"LOG_LEVEL" envDefault:"debug"`
}
type Config struct {
	AppConfig           `envPrefix:"APP_"`
	HTTPConfig          `envPrefix:"HTTP_"`
	RedisConfig         `envPrefix:"REDIS_"`
	PostgresConfig      `envPrefix:"POSTGRES_"`
	DBMigrationsConfig  `envPrefix:"MIGRATIONS_"`
	UncategorizedConfig `envPrefix:""`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	config, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
