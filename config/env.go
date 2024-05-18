package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	// PORT returns the server listening port
	PORT = "4242"
	// DB returns the database URL
	DB = ""
	// TOKENKEY returns the jwt token secret
	TOKENKEY = ""
	// TOKENEXP returns the jwt token expiration duration.
	// Should be time.ParseDuration string. Source: https://golang.org/pkg/time/#ParseDuration
	TOKENEXP = "24h"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("❌ Error loading environment variables")
	}

	// PORT returns the server listening port
	PORT = getEnv("PORT", "4242")
	// DB returns the database URL
	DB = getEnv("DB", "")
	// TOKENKEY returns the jwt token secret
	TOKENKEY = getEnv("TOKEN_KEY", "")
	// TOKENEXP returns the jwt token expiration duration.
	// Should be time.ParseDuration string. Source: https://golang.org/pkg/time/#ParseDuration
	TOKENEXP = getEnv("TOKEN_EXP", "24h")

	fmt.Println("🚀 Loaded environment variables")
}

func getEnv(name string, fallback string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}

	if fallback != "" {
		return fallback
	}

	panic(fmt.Sprintf(`Environment variable not found :: %v`, name))
}
