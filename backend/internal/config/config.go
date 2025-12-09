package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	JWTSecret         string
	JWTExpiration     time.Duration
	ServerPort        string
	ServerHost        string
	MaxFileSize       int64
	AllowedImageTypes []string
	AllowedAudioTypes []string
	EnableTLS         bool
	TLSCertFile       string
	TLSKeyFile        string
}

var AppConfig *Config

func LoadConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "10485760"), 10, 64)

	jwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	jwtExpiration := time.Duration(jwtExpHours) * time.Hour

	enableTLS, _ := strconv.ParseBool(getEnv("ENABLE_TLS", "false"))

	AppConfig = &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "messaging_db"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		JWTSecret:         getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		JWTExpiration:     jwtExpiration,
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		ServerHost:        getEnv("SERVER_HOST", "localhost"),
		MaxFileSize:       maxFileSize,
		AllowedImageTypes: parseCommaSeparated(getEnv("ALLOWED_IMAGE_TYPES", "jpg,jpeg,png,gif")),
		AllowedAudioTypes: parseCommaSeparated(getEnv("ALLOWED_AUDIO_TYPES", "mp3,wav,ogg,m4a")),
		EnableTLS:         enableTLS,
		TLSCertFile:       getEnv("TLS_CERT_FILE", "cert.pem"),
		TLSKeyFile:        getEnv("TLS_KEY_FILE", "key.pem"),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseCommaSeparated(s string) []string {
	result := []string{}
	for _, item := range splitComma(s) {
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func splitComma(s string) []string {
	result := []string{}
	start := 0
	for i, char := range s {
		if char == ',' {
			if i > start {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}
