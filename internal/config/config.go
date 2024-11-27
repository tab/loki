package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"loki/pkg/logger"
)

const (
	AppAddr   = "localhost:8080"
	ClientURL = "http://localhost:3000"
)

type Config struct {
	AppEnv    string
	AppAddr   string
	ClientURL string
}

func LoadConfig(log *logger.Logger) *Config {
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
		err := godotenv.Overload(file)
		if err == nil {
			log.Info().Msgf("Loaded %s file", file)
		}
	}

	flagAppAddr := flag.String("b", AppAddr, "server address")
	flagClientURL := flag.String("c", ClientURL, "client address")
	flag.Parse()

	return &Config{
		AppEnv:    env,
		AppAddr:   getFlagOrEnvString(*flagAppAddr, "APP_ADDRESS", AppAddr),
		ClientURL: getFlagOrEnvString(*flagClientURL, "CLIENT_URL", ClientURL),
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
