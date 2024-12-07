package config

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

const (
	AppAddr    = "localhost:8080"
	ClientURL  = "http://localhost:3000"
	DebugLevel = "debug"
)

type SmartId struct {
	BaseURL string

	RelyingPartyName string
	RelyingPartyUUID string

	Text string
}

type Config struct {
	AppEnv      string
	AppAddr     string
	ClientURL   string
	SecretKey   string
	DatabaseDSN string
	RedisURI    string
	SmartId     SmartId
	LogLevel    string
}

func LoadConfig() *Config {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	envFiles := []string{
		".env",
		fmt.Sprintf(".env.%s", env),
		fmt.Sprintf(".env.%s.local", env),
	}
	for _, file := range envFiles {
		_ = godotenv.Overload(file)
	}

	flagAppAddr := flag.String("b", AppAddr, "server address")
	flagClientURL := flag.String("c", ClientURL, "client address")
	flagSecretKey := flag.String("s", "", "JWT secret key")
	flagDatabaseDSN := flag.String("d", "", "database DSN")
	flagRedisURI := flag.String("r", "", "Redis URI")
	flag.Parse()

	return &Config{
		AppEnv:    env,
		AppAddr:   getFlagOrEnvString(*flagAppAddr, "APP_ADDRESS", AppAddr),
		ClientURL: getFlagOrEnvString(*flagClientURL, "CLIENT_URL", ClientURL),

		SecretKey:   getFlagOrEnvString(*flagSecretKey, "SECRET_KEY", ""),
		DatabaseDSN: getFlagOrEnvString(*flagDatabaseDSN, "DATABASE_DSN", ""),
		RedisURI:    getFlagOrEnvString(*flagRedisURI, "REDIS_URI", ""),

		SmartId: SmartId{
			BaseURL:          getEnvString("SMART_ID_API_URL"),
			RelyingPartyName: getEnvString("RELYING_PARTY_NAME"),
			RelyingPartyUUID: getEnvString("RELYING_PARTY_UUID"),
			Text:             getEnvString("SMART_ID_DISPLAY_TEXT"),
		},

		LogLevel: getEnvString("LOG_LEVEL"),
	}
}

func getFlagOrEnvString(flagValue, envVar, defaultValue string) string {
	if flagValue != "" {
		return flagValue
	}

	if envValue, ok := os.LookupEnv(envVar); ok && envValue != "" {
		return envValue
	}

	return defaultValue
}

func getEnvString(envVar string) string {
	if envValue, ok := os.LookupEnv(envVar); ok && envValue != "" {
		return envValue
	}

	return ""
}
