package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kilianmandscharo/lethimcook/env"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db         *authDatabase
	privateKey string
	logger     *logging.Logger
}

func NewAuthService(db *authDatabase, logger *logging.Logger) *AuthService {
	return &AuthService{
		db:         db,
		privateKey: env.Get(env.EnvKeyJWTPrivateKey),
		logger:     logger,
	}
}

func (as *AuthService) CreateAdminIfDoesNotExist(password string) {
	if as.doesAdminExist() {
		if len(password) > 0 {
			as.logger.Warn("Admin found, ignoring provided 'init-admin'")
		}
	} else {
		if len(password) == 0 {
			as.logger.Fatalf("no admin found, usage: %s [--init-admin] <password> [--prod]\n", os.Args[0])
		} else {
			as.createAdmin(password)
			as.logger.Info("Initialized admin")
		}
	}
}

func (as *AuthService) updateAdminPasswordHash(newPassword string) error {
	if len(newPassword) < 5 {
		return &errutil.AppError{
			UserMessage: "Invalides Passwort",
			Err: fmt.Errorf(
				"failed at updateAdminPasswordHash(): %w",
				errutil.FormErrorPasswortTooShort,
			),
			StatusCode: http.StatusBadRequest,
		}
	}
	if len(newPassword) > 72 {
		return &errutil.AppError{
			UserMessage: "Invalides Passwort",
			Err: fmt.Errorf(
				"failed at updateAdminPasswordHash(): %w",
				errutil.FormErrorPasswortTooLong,
			),
			StatusCode: http.StatusBadRequest,
		}
	}
	newPasswordHash, err := as.hashPassword(newPassword)
	if err != nil {
		return errutil.AddMessageToAppError(
			err,
			"failed at updateAdminPasswordHash()",
		)
	}
	err = as.db.updateAdminPasswordHash(newPasswordHash)
	if err != nil {
		return errutil.AddMessageToAppError(
			err,
			"failed at updateAdminPasswordHash()",
		)
	}
	return nil
}

func (as *AuthService) validatePassword(password string) error {
	hash, err := as.db.readAdminPasswordHash()
	if err != nil {
		return errutil.AddMessageToAppError(err, "failed at validatePassword()")
	}
	if !as.matchPassword(password, hash) {
		return &errutil.AppError{
			UserMessage: "Falsches Passwort",
			Err:         errors.New("failed at validatePassword(), invalid password"),
			StatusCode:  http.StatusUnauthorized,
		}
	}
	return nil
}

func (as *AuthService) doesAdminExist() bool {
	doesAdminExist, err := as.db.doesAdminExist()
	if err != nil {
		as.logger.Fatal(err)
	}
	return doesAdminExist
}

func (as *AuthService) createAdmin(password string) {
	passwordHash, err := as.hashPassword(password)
	if err != nil {
		as.logger.Fatal(err)
	}
	err = as.db.createAdmin(&admin{PasswordHash: passwordHash})
	if err != nil {
		as.logger.Fatal(err)
	}
}

func (as *AuthService) createToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(60 * time.Minute).UnixMilli(),
	})
	tokenString, err := token.SignedString([]byte(as.privateKey))
	if err != nil {
		return "", &errutil.AppError{
			UserMessage: "Serverfehler",
			Err:         fmt.Errorf("failed at createToken(): %w", err),
			StatusCode:  http.StatusInternalServerError,
		}
	}
	return tokenString, nil
}

func (as *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return "", &errutil.AppError{
			UserMessage: "Invalides Passwort",
			Err: fmt.Errorf(
				"failed at hashPassword(), password too long: %w",
				err,
			),
			StatusCode: http.StatusBadRequest,
		}
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

func (as *AuthService) createLoginForm(disabled bool, err error) []types.FormElement {
	return []types.FormElement{
		{
			Type:      types.FormElementInput,
			Name:      "password",
			Err:       err,
			InputType: "password",
			Label:     "Passwort",
			Required:  true,
			Disabled:  disabled,
		},
	}
}

func (as *AuthService) createNewPasswordForm(oldPasswordError error, newPasswordError error) []types.FormElement {
	return []types.FormElement{
		{
			Type:      types.FormElementInput,
			Name:      "old-password",
			Err:       oldPasswordError,
			InputType: "password",
			Label:     "Altes Passwort",
			Required:  true,
		},
		{
			Type:      types.FormElementInput,
			Name:      "new-password",
			Err:       newPasswordError,
			InputType: "password",
			Label:     "Neues Passwort",
			Required:  true,
		},
	}
}
