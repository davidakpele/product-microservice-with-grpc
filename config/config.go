package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	GRPCPort   string
}

// LoadConfig loads environment variables from .env
func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
 
	// Check for missing environment variables
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("DB_HOST is required but not set in the environment")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatal("DB_PORT is required but not set in the environment")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER is required but not set in the environment")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD is required but not set in the environment")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME is required but not set in the environment")
	}

	dbsslMode := os.Getenv("DB_SSLMODE")
	if dbsslMode == "" {
		log.Fatal("DB_SSLMODE is required but not set in the environment")
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		log.Fatal("GRPC_PORT is required but not set in the environment")
	}

	// Return the config
	return &Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBSSLMode:  dbsslMode,
		GRPCPort:   grpcPort,
	}
}
