package config

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"watchdog/database"
)

type Config struct {
	Server   ServerConfig
	Database database.Config
}

type ServerConfig struct {
	Port int
}

func Load() *Config {
	loadEnvFile()
	
	return &Config{
		Server: ServerConfig{
			Port: getIntEnv("PORT", 50051),
		},
		Database: database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getIntEnv("DB_PORT", 3306),
			Username: getEnv("DB_USERNAME", "watchdog"),
			Password: getEnv("DB_PASSWORD", "watchdog123"),
			Database: getEnv("DB_DATABASE", "watchdog_db"),
		},
	}
}

// LoadWithEntClient loads config and creates EntClient with auto-migration
func LoadWithEntClient() (*Config, *database.EntClient, error) {
	config := Load()
	
	entClient, err := database.NewEntClient(config.Database)
	if err != nil {
		return nil, nil, err
	}

	// Run auto-migration
	ctx := context.Background()
	if err := entClient.AutoMigrate(ctx); err != nil {
		entClient.Close()
		return nil, nil, err
	}

	return config, entClient, nil
}

func loadEnvFile() {
	envFiles := []string{
		".env.local",
		".env",
	}

	for _, envFile := range envFiles {
		if fileExists(envFile) {
			err := godotenv.Load(envFile)
			if err != nil {
				log.Printf("Warning: Error loading %s file: %v", envFile, err)
			} else {
				log.Printf("Loaded environment variables from %s", envFile)
				return
			}
		}
	}

	wd, err := os.Getwd()
	if err == nil {
		for _, envFile := range envFiles {
			fullPath := filepath.Join(wd, envFile)
			if fileExists(fullPath) {
				err := godotenv.Load(fullPath)
				if err != nil {
					log.Printf("Warning: Error loading %s file: %v", fullPath, err)
				} else {
					log.Printf("Loaded environment variables from %s", fullPath)
					return
				}
			}
		}
	}

	log.Println("No .env file found, using environment variables and defaults")
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}