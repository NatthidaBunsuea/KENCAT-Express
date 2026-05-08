// หมวย
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv          string
	Port            string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	JWTSecret       string
	JWTTTL          time.Duration
	PasswordSalt    string
	AutoMigrate     bool
	AutoSeed        bool
	CORSAllowOrigin string
	FrontendDir     string
}

func Load() Config {
	return Config{
		AppEnv:          env("APP_ENV", "development"),
		Port:            env("PORT", "8080"),
		DBHost:          env("DB_HOST", "127.0.0.1"),
		DBPort:          env("DB_PORT", "3306"),
		DBUser:          env("DB_USER", "root"),
		DBPassword:      env("DB_PASSWORD", "Muaymin_ly12548"),
		DBName:          env("DB_NAME", "kencat"),
		JWTSecret:       env("JWT_SECRET", "kencat-super-secret-key"),
		JWTTTL:          envDuration("JWT_TTL", 24*time.Hour),
		PasswordSalt:    env("PASSWORD_SALT", "kencat-express-salt"),
		AutoMigrate:     envBool("AUTO_MIGRATE", true),
		AutoSeed:        envBool("AUTO_SEED", true),
		CORSAllowOrigin: env("CORS_ALLOW_ORIGIN", "http://localhost:8080"),
		FrontendDir:     env("FRONTEND_DIR", "../Frontend"),
	}
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		parsed, err := time.ParseDuration(v)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
