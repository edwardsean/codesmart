package config

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	DBHost         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBPort         string
	JWTSecret      string
	GithubClientID string
	GithubSecret   string
	GolangAPIURL   string
	FrontendOrigin string
	EncryptionKey  string
}

var Envs = initConfig()

func initConfig() *Config {
	godotenv.Load("./backend/auth-service/.env")
	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", "mypassword"),
		DBName:         getEnv("DB_NAME", "codesmart"),
		DBPort:         getEnv("DB_PORT", "5432"),
		JWTSecret:      getEnv("JWT_SECRET", "jwt-secret-notfound"),
		GithubClientID: getEnv("GITHUB_ID", "not-found"),
		GithubSecret:   getEnv("GITHUB_SECRET", "not-found"),
		GolangAPIURL:   getEnv("GOLANG_API_URL", "not-found"),
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),
		EncryptionKey:  getEnv("ENCRYPTION_KEY", "not-found"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		log.Println("Value: ", value)
		return value
	}

	log.Println("Returning feedback")
	return fallback
}

// func getEnvAsInt(key string, fallback int64) int64 {
// 	if value, ok := os.LookupEnv(key); ok {
// 		i, err := strconv.ParseInt(value, 10, 64)
// 		if err != nil {
// 			return fallback
// 		}
// 		return i
// 	}

// 	return fallback
// }
