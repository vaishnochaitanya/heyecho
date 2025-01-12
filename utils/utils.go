package utils

import (
	"crypto/sha256"
	"os"
)

type Map map[string]any

func IsEnvProduction() bool {
	return os.Getenv("ENV") == "production"
}

func AppEnv() string {
	return os.Getenv("ENV")
}

func HashPassword(password string) (string, error) {
	bytes := sha256.Sum256([]byte(password))
	return string(bytes[:]), nil
}
