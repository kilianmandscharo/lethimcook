package env

import (
	"fmt"
	"os"
	"testing"

	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/stretchr/testify/assert"
)

const (
	testFileName      = ".env.test"
	testCertFilePath  = "test_cert_file_path"
	testKeyFilePath   = "test_key_file_path"
	testJWTPrivateKey = "test_jwt_private_key"
)

func TestEnvironment(t *testing.T) {
	// Given
	envString := fmt.Sprintf("%s=%s", EnvKeyCertFilePath, testCertFilePath) +
		fmt.Sprintf("\n%s=%s", EnvKeyKeyFilePath, testKeyFilePath) +
		fmt.Sprintf("\n%s=%s", EnvKeyJWTPrivateKey, testJWTPrivateKey)

	assert.NoError(t, os.WriteFile(testFileName, []byte(envString), 0644))

	// Then
	assert.Equal(t, "", Get(EnvKeyCertFilePath))
	assert.Equal(t, "", Get(EnvKeyKeyFilePath))
	assert.Equal(t, "", Get(EnvKeyJWTPrivateKey))

	// When
    logger := logging.New(logging.Debug)
	LoadEnvironment(testFileName, &logger)

	// Then
	assert.Equal(t, testCertFilePath, Get(EnvKeyCertFilePath))
	assert.Equal(t, testKeyFilePath, Get(EnvKeyKeyFilePath))
	assert.Equal(t, testJWTPrivateKey, Get(EnvKeyJWTPrivateKey))

	assert.NoError(t, os.Remove(testFileName))
}
