package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	AppAddr    = "0.0.0.0:8080"
	ClientURL  = "http://localhost:3000"
	DebugLevel = "debug"
)

type SmartId struct {
	BaseURL string

	RelyingPartyName string
	RelyingPartyUUID string

	Text string
}

type MobileId struct {
	BaseURL string

	RelyingPartyName string
	RelyingPartyUUID string

	Text       string
	TextFormat string

	Language string
}

type Config struct {
	AppEnv       string
	AppName      string
	AppAddr      string
	ClientURL    string
	SecretKey    string
	DatabaseDSN  string
	RedisURI     string
	TelemetryURI string
	SmartId      SmartId
	MobileId     MobileId
	LogLevel     string
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
	flagTelemetryURI := flag.String("t", "", "OpenTelemetry collector URI")
	flag.Parse()

	return &Config{
		AppEnv:    env,
		AppName:   getEnvString("APP_NAME"),
		AppAddr:   getFlagOrEnvString(*flagAppAddr, "APP_ADDRESS", AppAddr),
		ClientURL: getFlagOrEnvString(*flagClientURL, "CLIENT_URL", ClientURL),

		SecretKey:    getFlagOrEnvString(*flagSecretKey, "SECRET_KEY", ""),
		DatabaseDSN:  getFlagOrEnvString(*flagDatabaseDSN, "DATABASE_DSN", ""),
		RedisURI:     getFlagOrEnvString(*flagRedisURI, "REDIS_URI", ""),
		TelemetryURI: getFlagOrEnvString(*flagTelemetryURI, "TELEMETRY_URI", ""),

		SmartId: SmartId{
			BaseURL:          getEnvString("SMART_ID_API_URL"),
			RelyingPartyName: getEnvString("RELYING_PARTY_NAME"),
			RelyingPartyUUID: getEnvString("RELYING_PARTY_UUID"),
			Text:             getEnvString("SMART_ID_DISPLAY_TEXT"),
		},
		MobileId: MobileId{
			BaseURL:          getEnvString("MOBILE_ID_API_URL"),
			RelyingPartyName: getEnvString("RELYING_PARTY_NAME"),
			RelyingPartyUUID: getEnvString("RELYING_PARTY_UUID"),
			Text:             getEnvString("MOBILE_ID_DISPLAY_TEXT"),
			TextFormat:       getEnvString("MOBILE_ID_TEXT_FORMAT"),
			Language:         getEnvString("MOBILE_ID_LANGUAGE"),
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
