package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthRequestHandler struct {
	db AuthDatabase
}

func NewAuthRequestHandler() AuthRequestHandler {
	return AuthRequestHandler{
		db: NewAuthDatabase(),
	}
}

type login struct {
	Password string `json:"password"`
}

func (r *AuthRequestHandler) HandleLogin(c echo.Context) error {
	var login login

	if err := c.Bind(&login); err != nil {
		return c.String(http.StatusBadRequest, "invalid body")
	}

	match, err := r.db.validatePassword(login.Password)

	if err != nil {
		return c.String(http.StatusInternalServerError, "error validating password")
	}

	if !match {
		return c.Render(http.StatusUnauthorized, "failed-login.html", nil)
	}

	return c.String(http.StatusOK, "login successful")
}

type passwordUpdate struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (r *AuthRequestHandler) HandleUpdatePassword(c echo.Context) error {
	var passwordUpdate passwordUpdate

	if err := c.Bind(&passwordUpdate); err != nil {
		return c.String(http.StatusBadRequest, "invalid body")
	}

	match, err := r.db.validatePassword(passwordUpdate.OldPassword)

	if err != nil {
		return c.String(http.StatusInternalServerError, "error validating password")
	}

	if !match {
		return c.String(http.StatusUnauthorized, "wrong password")
	}

	err = r.db.updateAdminPasswordHash(passwordUpdate.NewPassword)

	if err != nil {
		return c.String(http.StatusInternalServerError, "error updating password")
	}

	return c.String(http.StatusOK, "password updated successfully")
}
