package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Server   ServerConfig
	Auth     AuthConfig
	Security SecurityConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type AuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	CallbackURL        string
	Scopes             []string
}

type SecurityConfig struct {
	SessionSecret string
	IsProd        bool
	CookieMaxAge  int
}

// Load carga y valida la configuración desde variables de entorno
func Load() (*Config, error) {
	// if err := godotenv.Load(); err != nil {
	// 	return nil, fmt.Errorf("error loading .env file: %w", err)
	// }

	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", "8080"),
			Host: getEnvOrDefault("HOST", "localhost"),
		},
		Auth: AuthConfig{
			GoogleClientID:     mustGetEnv("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: mustGetEnv("GOOGLE_CLIENT_SECRET"),
			CallbackURL:        getEnvOrDefault("AUTH_CALLBACK_URL", "https://go-auth-google.onrender.com/auth/google/callback"),
			Scopes:             []string{"email", "profile"},
		},
		Security: SecurityConfig{
			SessionSecret: mustGetEnv("SESSION_SECRET"),
			IsProd:        getEnvOrDefault("ENV", "development") == "production",
			CookieMaxAge:  86400 * 30, // 30 días por defecto
		},
	}

	return config, nil
}

// mustGetEnv obtiene una variable de entorno o paniquea si no existe
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("environment variable %s is required", key))
	}
	return value
}

// getEnvOrDefault obtiene una variable de entorno o retorna un valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
