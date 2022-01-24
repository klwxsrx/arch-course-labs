package main

import (
	"fmt"
	"os"
)

type config struct {
	DBName     string
	DBAddress  string
	DBUser     string
	DBPassword string
}

func parseEnvString(key string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	str, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("undefined environment variable %s", key)
	}
	return str, nil
}

func parseConfig() (*config, error) {
	var err error
	dbName, err := parseEnvString("DATABASE_NAME", err)
	dbAddress, err := parseEnvString("DATABASE_ADDRESS", err)
	dbUser, err := parseEnvString("DATABASE_USER", err)
	dbPassword, err := parseEnvString("DATABASE_PASSWORD", err)

	if err != nil {
		return nil, err
	}

	return &config{
		dbName,
		dbAddress,
		dbUser,
		dbPassword,
	}, nil
}
