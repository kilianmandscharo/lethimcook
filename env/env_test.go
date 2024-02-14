package env

import (
	"fmt"
	"os"
	"testing"

	"github.com/kilianmandscharo/lethimcook/projectpath"
	"github.com/stretchr/testify/assert"
)

const (
	testFileName      = ".env.test"
	testCertFilePath  = "test_cert_file_path"
	testKeyFilePath   = "test_key_file_path"
	testJWTPrivateKey = "test_jwt_private_key"
)

var testFilePath = projectpath.Absolute(testFileName)

func TestEnvironment(t *testing.T) {
	// Given
	envString := fmt.Sprintf("%s=%s", EnvKeyCertFilePath, testCertFilePath) +
		fmt.Sprintf("\n%s=%s", EnvKeyKeyFilePath, testKeyFilePath) +
		fmt.Sprintf("\n%s=%s", EnvKeyJWTPrivateKey, testJWTPrivateKey)

	err := os.WriteFile(testFilePath, []byte(envString), 0644)
	assert.NoError(t, err)

	// Then
	assert.Equal(t, "", Get(EnvKeyCertFilePath))
	assert.Equal(t, "", Get(EnvKeyKeyFilePath))
	assert.Equal(t, "", Get(EnvKeyJWTPrivateKey))

	// When
	LoadEnvironment(testFileName)

	// Then
	assert.Equal(t, testCertFilePath, Get(EnvKeyCertFilePath))
	assert.Equal(t, testKeyFilePath, Get(EnvKeyKeyFilePath))
	assert.Equal(t, testJWTPrivateKey, Get(EnvKeyJWTPrivateKey))

	err = os.Remove(testFilePath)
	assert.NoError(t, err)
}
