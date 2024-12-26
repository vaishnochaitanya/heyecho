package util

import (
	"os"
)

type Map map[string]any

func IsEnvProduction() bool {
	return os.Getenv("ENV") == "production"
}

func AppEnv() string {
	return os.Getenv("ENV")
}
