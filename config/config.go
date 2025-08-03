// Package config provides configuration loading and management for the application.
package config

import (
	"fmt"
	"os"
	"sync"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zenv"
	"github.com/joho/godotenv"
	"github.com/ship-labs/meet-loop-api/internal"
)

var (
	config Config
	err    error
	once   sync.Once
)

type Config struct {
	Env                string `env:"Env" zog:"Env"`
	Port               int    `env:"PORT" zog:"Port"`
	FrontendURL        string `env:"FRONTEND_URL" zog:"FrontendURL"`
	DBURL              string `env:"DB_URL" zog:"DBURL"`
	DBPassword         string `env:"DB_PASSWORD" zog:"DBPassword"`
	JwtSecret          string `env:"JWT_SECRET" zog:"JwtSecret"`
	SupabaseProjectURL string `env:"SUPABASE_PROJECT_URL" zog:"SupabaseProjectURL"`
	SupabaseAPIKey     string `env:"SUPABASE_API_KEY" zog:"SupabaseAPIKey"`
}

const (
	DefaultPort     = 8080
	DevEnvironment  = "development"
	ProdEnvironment = "production"
)

func loadConfig() (Config, error) {
	if os.Getenv("Env") != ProdEnvironment {
		if err := godotenv.Load(); err != nil {
			return Config{}, fmt.Errorf("loading environment variables: %w", err)
		}
	}

	schema := z.Struct(z.Shape{
		"Port":               z.Int().Default(DefaultPort),
		"Env":                z.String().Required().OneOf([]string{DevEnvironment, ProdEnvironment}),
		"FrontendURL":        z.String().URL().Required(),
		"DBURL":              z.String().URL().Required(),
		"DBPassword":         z.String().Required(),
		"JwtSecret":          z.String().Required(),
		"SupabaseProjectURL": z.String().URL().Required(),
		"SupabaseAPIKey":     z.String().Required(),
	})

	var c Config

	errs := schema.Parse(zenv.NewDataProvider(), &c)
	if len(errs) != 0 {
		return c, fmt.Errorf("loading config: %s", internal.ValidationError(errs))
	}

	return c, nil
}

func LoadConfig() (Config, error) {
	once.Do(func() {
		config, err = loadConfig()
	})

	return config, err
}
