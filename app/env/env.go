package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kilianmandscharo/lethimcook/logging"
)

const (
	EnvKeyCertFilePath  = "CERT_FILE_PATH"
	EnvKeyKeyFilePath   = "KEY_FILE_PATH"
	EnvKeyJWTPrivateKey = "JWT_PRIVATE_KEY"
)

func LoadEnvironment(envName string, logger *logging.Logger) {
	err := godotenv.Load(envName)
	if err != nil {
		logger.Fatal("failed to load", envName)
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
