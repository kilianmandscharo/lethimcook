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

func TestCreateAdminIfDoesNotExistCrash(t *testing.T) {
	// Given
	a := newTestAuthService()
	if os.Getenv("BE_CRASHER") == "1" {
		a.CreateAdminIfDoesNotExist("")
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
	a := newTestAuthService()

	// Then
	a.CreateAdminIfDoesNotExist("test_password")
	a.CreateAdminIfDoesNotExist("test_password")
	a.CreateAdminIfDoesNotExist("")
}

func TestUpdateAdminPasswordHash(t *testing.T) {
	// Given
	a := newTestAuthService()

	// When
	err := a.updateAdminPasswordHash("test_password")

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorNoAdminFound)

	// Given
	admin := newTestAdmin()
	assert.NoError(t, a.db.createAdmin(&admin))

	// When
	err = a.updateAdminPasswordHash("aaaa")

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorPasswordTooShort)

	// When
	err = a.updateAdminPasswordHash("test_password")

	// Then
	assert.NoError(t, err)
}

func TestValidatePassword(t *testing.T) {
	// Given
	a := newTestAuthService()

	// When
	err := a.validatePassword("test_password")

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorNoAdminFound)

	// Given
	testHash, err := a.hashPassword("test_password")
	assert.NoError(t, err)
	admin := admin{PasswordHash: testHash}
	assert.NoError(t, a.db.createAdmin(&admin))

	// When
	err = a.validatePassword("test_password")

	// Then
	assert.NoError(t, err)

	// When
	err = a.validatePassword("invalid_password")

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorInvalidPassword)
}

func TestCreateToken(t *testing.T) {
	// Given
	a := newTestAuthService()

	// When
	token, err := a.createToken()

	// Then
	assert.NoError(t, err)
	assert.True(t, len(token) != 0)
}

func TestHashPassword(t *testing.T) {
	// Given
	a := newTestAuthService()

	// When
	hash, err := a.hashPassword("test_password")

	// Then
	assert.NoError(t, err)
	assert.True(t, len(hash) != 0)

	// When
	hash, err = a.hashPassword("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	// Then
	assert.ErrorIs(t, err, errutil.AuthErrorPasswordTooLong)
}

func TestMatchPassword(t *testing.T) {
	// Given
	a := newTestAuthService()
	validPassword := "test_password"
	hash, err := a.hashPassword(validPassword)
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
		assert.Equal(t, test.shouldBeValid, a.matchPassword(test.password, hash))
	}
}

func TestValidCookieToken(t *testing.T) {
	// Given
	a := newTestAuthService()
	token, err := a.createToken()
	assert.NoError(t, err)

	testCases := []struct {
		cookie    http.Cookie
		wantValid bool
	}{
		{
			cookie:    a.newTokenCookie(token, time.Now().Add(60*time.Minute)),
			wantValid: true,
		},
		{
			cookie:    a.newTokenCookie("invalid_token", time.Now().Add(60*time.Minute)),
			wantValid: false,
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.wantValid, a.validCookieToken(&test.cookie))
	}
}
