package config

import (
	"log"
	"os"
)

type Config struct {
	Port            string
	AuthToken       string
	DBURL           string
	BoxofficeURL    string
	BoxofficeAPIKey string
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return v
}

func Load() Config {
	return Config{
		Port:            getEnv("PORT", "8080"),
		AuthToken:       mustEnv("AUTH_TOKEN"),
		DBURL:           mustEnv("DB_URL"),
		BoxofficeURL:    mustEnv("BOXOFFICE_URL"),
		BoxofficeAPIKey: mustEnv("BOXOFFICE_API_KEY"),
	}
}
