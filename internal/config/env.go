package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	ServiceName      string         `envconfig:"service_name" required:"true" default:"ciam"`
	BindAddress      string         `envconfig:"bind_address" required:"true" default:":13702"`
	Issuer           string         `envconfig:"issuer" required:"true" default:"http://localhost:13702"`
	CryptoKey        string         `envconfig:"crypto_key" required:"true"`
	EnableTracing    bool           `envconfig:"enable_tracing" default:"false"`
	EnableMetrics    bool           `envconfig:"enable_metrics" default:"false"`
	Postgres         PostgresConfig `envconfig:"postgres"`
	DevLogin         DevLoginConfig `envconfig:"dev_login"`
}

// DevLoginConfig is the stand-in for a real account store. The built-in login
// page is disabled unless both values are set.
type DevLoginConfig struct {
	Email    string `envconfig:"email"`
	Password string `envconfig:"password"`
}

type PostgresConfig struct {
	DSN          string `envconfig:"dsn" required:"true"`
	MaxOpenConns int    `envconfig:"max_open_conns" default:"10"`
	MaxIdleConns int    `envconfig:"max_idle_conns" default:"2"`
}

func LoadConfig(envFiles ...string) *EnvConfig {
	if err := godotenv.Load(envFiles...); err != nil {
		log.Fatalf("no .env file found or error loading .env file: %v", err)
	}

	var cfg EnvConfig
	if err := envconfig.Process("CIAM", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	return &cfg
}
