package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kilianmandscharo/lethimcook/components"
	"github.com/kilianmandscharo/lethimcook/env"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db         authDatabase
	privateKey string
}

func NewAuthService() AuthService {
	return AuthService{
		db:         newAuthDatabase(),
		privateKey: env.Get(env.EnvKeyJWTPrivateKey),
	}
}

func (as *AuthService) CreateAdminIfDoesNotExist(password string) {
	if as.doesAdminExist() {
		if len(password) > 0 {
			fmt.Println("Admin found, ignoring provided 'init-admin'")
		}
	} else {
		if len(password) == 0 {
			fmt.Println("No admin found")
			fmt.Printf("Usage: %s --init-admin <password>\n", os.Args[0])
			os.Exit(1)
		} else {
			as.createAdmin(password)
			fmt.Println("Initialized admin")
		}
	}
}

func (as *AuthService) updateAdminPasswordHash(newPassword string) errutil.AuthError {
	if len(newPassword) < 5 {
		return errutil.AuthErrorPasswordTooShort
	}

	newPasswordHash, err := as.hashPassword(newPassword)
	if err != nil {
		return err
	}

	return as.db.updateAdminPasswordHash(newPasswordHash)
}

func (as *AuthService) validatePassword(password string) errutil.AuthError {
	hash, err := as.db.readAdminPasswordHash()
	if err != nil {
		return err
	}

	if !as.matchPassword(password, hash) {
		return errutil.AuthErrorInvalidPassword
	}

	return nil
}

func (as *AuthService) doesAdminExist() bool {
	doesAdminExist, err := as.db.doesAdminExist()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return doesAdminExist
}

func (as *AuthService) createAdmin(password string) {
	passwordHash, err := as.hashPassword(password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = as.db.createAdmin(&admin{PasswordHash: passwordHash})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (as *AuthService) createToken() (string, errutil.AuthError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(60 * time.Minute).UnixMilli(),
	})

	tokenString, err := token.SignedString([]byte(as.privateKey))
	if err != nil {
		return "", errutil.AuthErrorTokenCreationFailure
	}

	return tokenString, nil
}

func (as *AuthService) hashPassword(password string) (string, errutil.AuthError) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return "", errutil.AuthErrorPasswordTooLong
	}

	return string(bytes), nil
}

func (as *AuthService) matchPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (as *AuthService) newTokenCookie(token string, expires time.Time) http.Cookie {
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

func (as *AuthService) validCookieToken(cookie *http.Cookie) bool {
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(as.privateKey), nil
	})

	return err == nil && token.Valid
}

func (as *AuthService) createLoginForm(disabled bool, err errutil.AuthError) []types.FormElement {
	return []types.FormElement{
		{
			Type:      types.FormElementInput,
			Name:      "password",
			Err:       err,
			InputType: "text",
			Label:     "Passwort",
			Required:  true,
			Disabled:  disabled,
		},
	}
}

func (as *AuthService) createNewPasswordForm(oldPasswordError errutil.AuthError, newPasswordError errutil.AuthError) []types.FormElement {
	return []types.FormElement{
		{
			Type:      types.FormElementInput,
			Name:      "old-password",
			Err:       oldPasswordError,
			InputType: "text",
			Label:     "Altes Passwort",
			Required:  true,
		},
		{
			Type:      types.FormElementInput,
			Name:      "new-password",
			Err:       newPasswordError,
			InputType: "text",
			Label:     "Neues Passwort",
			Required:  true,
		},
	}
}

type createAdminPageOptions struct {
	c                echo.Context
	isAuthorized     bool
	loginFormError   error
	message          string
	err              error
	oldPasswordError error
	newPasswordError error
}

func (as *AuthService) createAdminPage(options createAdminPageOptions) error {
	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context: options.c,
		Component: components.AdminPage(
			options.isAuthorized,
			as.createLoginForm(options.isAuthorized, options.loginFormError),
			as.createNewPasswordForm(options.oldPasswordError, options.newPasswordError),
		),
		Message: options.message,
		Err:     options.err,
	})
}
