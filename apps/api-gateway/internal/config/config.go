package config

import (
	"os"
)

type Config struct {
	Server           ServerConfig
	Database         DatabaseConfig
	MinIO            MinIOConfig
	UniDocLicenseKey string
	Ollama           OllamaConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type MinIOConfig struct {
	HostAPI     string
	PortAPI     string
	HostConsole string
	PortConsole string
	User        string
	Password    string
}

type OllamaConfig struct {
	Host  string
	Port  string
	Model string
}

func getEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		value = defaultValue[0]
	}
	return value
}

func Load() Config {
	cfg := Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "nexus"),
		},
		MinIO: MinIOConfig{
			HostAPI:     getEnv("MINIO_HOST_API", "localhost"),
			PortAPI:     getEnv("MINIO_PORT_API", "9000"),
			HostConsole: getEnv("MINIO_HOST_CONSOLE", "localhost"),
			PortConsole: getEnv("MINIO_PORT_CONSOLE", "9001"),
			User:        getEnv("MINIO_USER", "admin"),
			Password:    getEnv("MINIO_PASSWORD", "password123"),
		},
		Ollama: OllamaConfig{
			Host:  getEnv("OLLAMA_HOST", "localhost"),
			Port:  getEnv("OLLAMA_PORT", "11434"),
			Model: getEnv("OLLAMA_MODEL", "phi3:mini"),
		},
	}

	return cfg
}
