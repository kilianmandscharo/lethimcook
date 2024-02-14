package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kilianmandscharo/lethimcook/env"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db         authDatabase
	privateKey string
}

func newTestAuthService() AuthService {
	return AuthService{
		db:         newTestAuthDatabase(),
		privateKey: "test_private_key",
	}
}

func NewAuthService() AuthService {
	return AuthService{
		db:         newAuthDatabase(),
		privateKey: env.Get(env.EnvKeyJWTPrivateKey),
	}
}

func (a *AuthService) CreateAdminIfDoesNotExist(password string) {
	if a.doesAdminExist() {
		if len(password) > 0 {
			fmt.Println("Admin found, ignoring provided 'init-admin'")
		}
	} else {
		if len(password) == 0 {
			fmt.Println("No admin found")
			fmt.Printf("Usage: %s --init-admin <password>\n", os.Args[0])
			os.Exit(1)
		} else {
			a.createAdmin(password)
			fmt.Println("Initialized admin")
		}
	}
}

func (a *AuthService) updateAdminPasswordHash(newPassword string) errutil.AuthError {
	if len(newPassword) < 5 {
		return errutil.AuthErrorPasswordTooShort
	}

	newPasswordHash, err := a.hashPassword(newPassword)
	if err != nil {
		return err
	}

	return a.db.updateAdminPasswordHash(newPasswordHash)
}

func (a *AuthService) validatePassword(password string) errutil.AuthError {
	hash, err := a.db.readAdminPasswordHash()
	if err != nil {
		return err
	}

	if !a.matchPassword(password, hash) {
		return errutil.AuthErrorInvalidPassword
	}

	return nil
}

func (a *AuthService) doesAdminExist() bool {
	doesAdminExist, err := a.db.doesAdminExist()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return doesAdminExist
}

func (a *AuthService) createAdmin(password string) {
	passwordHash, err := a.hashPassword(password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = a.db.createAdmin(&admin{PasswordHash: passwordHash})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (a *AuthService) createToken() (string, errutil.AuthError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(60 * time.Minute).UnixMilli(),
	})

	tokenString, err := token.SignedString([]byte(a.privateKey))
	if err != nil {
		return "", errutil.AuthErrorTokenCreationFailure
	}

	return tokenString, nil
}

func (a *AuthService) hashPassword(password string) (string, errutil.AuthError) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return "", errutil.AuthErrorPasswordTooLong
	}

	return string(bytes), nil
}

func (a *AuthService) matchPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (a *AuthService) newTokenCookie(token string, expires time.Time) http.Cookie {
	return http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expires,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
}

func (r *AuthService) validCookieToken(cookie *http.Cookie) bool {
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.privateKey), nil
	})

	return err == nil && token.Valid
}
