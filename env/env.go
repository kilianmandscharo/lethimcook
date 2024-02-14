package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kilianmandscharo/lethimcook/projectpath"
)

const (
	EnvKeyCertFilePath  = "CERT_FILE_PATH"
	EnvKeyKeyFilePath   = "KEY_FILE_PATH"
	EnvKeyJWTPrivateKey = "JWT_PRIVATE_KEY"
)

func LoadEnvironment(envName string) {
	err := godotenv.Load(projectpath.Absolute(envName))
	if err != nil {
		fmt.Println("failed to load ", envName)
		os.Exit(1)
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
