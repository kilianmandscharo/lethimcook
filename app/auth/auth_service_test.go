package auth

import (
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/stretchr/testify/assert"
)

func newTestAuthService() AuthService {
	return AuthService{
		db:         newTestAuthDatabase(),
		privateKey: "test_private_key",
	}
}

func TestCreateAdminIfDoesNotExistCrash(t *testing.T) {
	// Given
	authService := newTestAuthService()
	if os.Getenv("BE_CRASHER") == "1" {
		authService.CreateAdminIfDoesNotExist("")
	}

	// When
	cmd := exec.Command(os.Args[0], "-test.run=TestCreateAdminIfDoesNotExistCrash")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	// Then
	e, ok := err.(*exec.ExitError)
	assert.True(t, ok)
	assert.False(t, e.Success())
}

func TestCreateAdminIfDoesNotExist(t *testing.T) {
	// Given
	authService := newTestAuthService()

	// Then
	authService.CreateAdminIfDoesNotExist("test_password")
	authService.CreateAdminIfDoesNotExist("test_password")
	authService.CreateAdminIfDoesNotExist("")
}

func TestUpdateAdminPasswordHash(t *testing.T) {
	// Given
	authService := newTestAuthService()

	// When
	err := authService.updateAdminPasswordHash("test_password")

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorNoAdminFound)

	// Given
	admin := newTestAdmin()
	assert.NoError(t, authService.db.createAdmin(&admin))

	// When
	err = authService.updateAdminPasswordHash("aaaa")

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorPasswordTooShort)

	// When
	err = authService.updateAdminPasswordHash("test_password")

	// Then
	assert.NoError(t, err)
}

func TestValidatePassword(t *testing.T) {
	// Given
	authService := newTestAuthService()

	// When
	err := authService.validatePassword("test_password")

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorNoAdminFound)

	// Given
	testHash, err := authService.hashPassword("test_password")
	assert.NoError(t, err)
	admin := admin{PasswordHash: testHash}
	assert.NoError(t, authService.db.createAdmin(&admin))

	// When
	err = authService.validatePassword("test_password")

	// Then
	assert.NoError(t, err)

	// When
	err = authService.validatePassword("invalid_password")

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorInvalidPassword)
}

func TestCreateToken(t *testing.T) {
	// Given
	authService := newTestAuthService()

	// When
	token, err := authService.createToken()

	// Then
	assert.NoError(t, err)
	assert.True(t, len(token) != 0)
}

func TestHashPassword(t *testing.T) {
	// Given
	authService := newTestAuthService()

	// When
	hash, err := authService.hashPassword("test_password")

	// Then
	assert.NoError(t, err)
	assert.True(t, len(hash) != 0)

	// When
	hash, err = authService.hashPassword("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorPasswordTooLong)
}

func TestMatchPassword(t *testing.T) {
	// Given
	authService := newTestAuthService()
	validPassword := "test_password"
	hash, err := authService.hashPassword(validPassword)
	assert.NoError(t, err)

	testCases := []struct {
		password      string
		shouldBeValid bool
	}{
		{
			password:      validPassword,
			shouldBeValid: true,
		},
		{
			password:      "invalid_password",
			shouldBeValid: false,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.shouldBeValid, authService.matchPassword(test.password, hash))
	}
}

func TestValidCookieToken(t *testing.T) {
	// Given
	authService := newTestAuthService()
	token, err := authService.createToken()
	assert.NoError(t, err)

	testCases := []struct {
		cookie    http.Cookie
		wantValid bool
	}{
		{
			cookie:    authService.newTokenCookie(token, time.Now().Add(60*time.Minute)),
			wantValid: true,
		},
		{
			cookie:    authService.newTokenCookie("invalid_token", time.Now().Add(60*time.Minute)),
			wantValid: false,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.wantValid, authService.validCookieToken(&test.cookie))
	}
}
