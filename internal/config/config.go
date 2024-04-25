package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Name        string
	URI         string
	Collections DatabaseCollectionConfig
}

type DatabaseCollectionConfig struct {
	Users string
	Forms string
}

type JWTConfig struct {
	Secret     string
	Expiration string
}

var cachedConfig *Config

func loadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", ""),
		},
		Database: DatabaseConfig{
			Name: getEnv("DATABASE_NAME", ""),
			URI:  getEnv("DATABASE_URI", ""),
			Collections: DatabaseCollectionConfig{
				Users: getEnv("DATABASE_COLLECTION_USERS", ""),
				Forms: getEnv("DATABASE_COLLECTION_FORMS", ""),
			},
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			Expiration: getEnv("JWT_EXPIRATION", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func GetConfig() *Config {
	if cachedConfig == nil {
		cachedConfig = loadConfig()
	}

	return cachedConfig
}

func (c *Config) GetServer() ServerConfig {
	return c.Server
}

func (c *Config) GetDatabaseConfig() DatabaseConfig {
	return c.Database
}

func (c *Config) GetJWTConfig() JWTConfig {
	return c.JWT
}
