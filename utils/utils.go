package utils

import (
	"crypto/rand"
	"log/slog"
	"math/big"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Map map[string]any

func IsEnvProduction() bool {
	slog.Debug("Checking environment", "env", os.Getenv("ENV"))
	return os.Getenv("ENV") == "production"
}

func AppEnv() string {
	return os.Getenv("ENV")
}

func GeneratePassphrase(wordCount ...int) string {
	// Default to 4 words if not specified
	count := 4
	if len(wordCount) > 0 && wordCount[0] > 0 {
		count = wordCount[0]
	}

	// List of common words to generate readable passphrases
	words := []string{
		"apple", "banana", "orange", "grape", "melon", "cherry", "peach", "lemon", "plum", "kiwi",
		"castle", "mountain", "river", "forest", "ocean", "desert", "valley", "island", "beach", "hill",
		"happy", "clever", "brave", "calm", "strong", "gentle", "wise", "quick", "bright", "careful",
		"jump", "dance", "sing", "write", "read", "sleep", "walk", "run", "swim", "climb",
		"blue", "green", "red", "yellow", "purple", "white", "black", "orange", "brown", "gray",
		"time", "world", "year", "day", "night", "hour", "minute", "second", "moment", "decade",
		"table", "chair", "window", "door", "floor", "wall", "ceiling", "lamp", "desk", "shelf",
	}

	var result []string
	for i := 0; i < count; i++ {
		// Use crypto/rand for better randomness
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
		if err != nil {
			// Fallback to first word on error
			slog.Warn("Error generating passphrase", "error", err)
			result = append(result, words[0])
		} else {
			slog.Debug("Generated passphrase", "word", words[nBig.Int64()])
			result = append(result, words[nBig.Int64()])
		}
	}

	return strings.Join(result, "-")
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Debug("Password hashing failed", "error", err)
		return "", err
	}
	slog.Debug("Password hashed successfully")
	return string(hashedPassword), nil
}

func VerifyPasswordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		slog.Debug("Password verification failed", "error", err)
		return false
	}
	slog.Debug("Password verified successfully")
	return true
}
