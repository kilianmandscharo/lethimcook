package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvKeyCertFilePath  = "CERT_FILE_PATH"
	EnvKeyKeyFilePath   = "KEY_FILE_PATH"
	EnvKeyJWTPrivateKey = "JWT_PRIVATE_KEY"
)

func LoadEnvironment() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("failed to load .env")
		os.Exit(1)
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
