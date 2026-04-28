package config

import (
	"log"
	"os"
	"strconv"
)

/*
	To use in other part of app
	cfg := config.LoadConfig()
	dsn := cfg.Postgres.DSN()
*/

type Config struct {
	App      AppConfig
	Auth     AuthConfig
	Postgres PostgresConfig
	Minio    MinioConfig
	AI		 AIConfig
}

type AppConfig struct {
	Environment string
	IsDev       bool
	Port		string
	Proxies		string
}

type AuthConfig struct {
	JWTSecret string
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

type AIConfig struct {
	IP string
}

// LoadConfig loads and returns application config
func LoadConfig() *Config {
	cfg := &Config{}

	// App env
	cfg.App.Environment = getEnv("APP_ENV", "development")
	cfg.App.Port = getEnv("PORT", "8080")
	cfg.App.Proxies = getEnv("PROXIES", "*")
	cfg.App.IsDev = cfg.App.Environment == "development"

	//auth
	cfg.Auth.JWTSecret = getEnv("JWT_SECRET", "dev-secret")

	// Postgres
	cfg.Postgres.Host = getEnv("POSTGRES_HOST", "localhost")
	cfg.Postgres.User = getEnv("POSTGRES_USER", "postgres")
	cfg.Postgres.Password = getEnv("POSTGRES_PASSWORD", "")
	cfg.Postgres.DBName = getEnv("POSTGRES_DB", "postgres")
	cfg.Postgres.SSLMode = getEnv("POSTGRES_SSLMODE", "disable")

	port, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid POSTGRES_PORT: %v", err)
	}
	cfg.Postgres.Port = port

	//minio
	cfg.Minio.Endpoint = getEnv("MINIO_ENDPOINT", "localhost:9000")
	cfg.Minio.AccessKey = getEnv("MINIO_ACCESS_KEY", "minioadmin")
	cfg.Minio.SecretKey = getEnv("MINIO_SECRET_KEY", "minioadmin")
	cfg.Minio.Bucket = getEnv("MINIO_BUCKET", "uploads")

	cfg.Minio.UseSSL = getEnv("MINIO_USE_SSL", "0") == "1"

	//ai
	cfg.AI.IP = getEnv("AI_IP", "localhost:8080")

	return cfg
}

// DSN builder for GORM/Postgres
func (p PostgresConfig) DSN() string {

	return "host=" + p.Host +
		" user=" + p.User +
		" password=" + p.Password +
		" dbname=" + p.DBName +
		" port=" + strconv.Itoa(p.Port) +
		" sslmode=" + p.SSLMode +
		" connect_timeout=5"
}

// Helper env loader
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
