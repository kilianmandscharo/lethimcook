package errutil

import (
	"errors"
	"net/http"
)

type AuthError = error

var (
	AuthErrorNoAdminFound         AuthError = errors.New("no admin found")
	AuthErrorDatabaseFailure      AuthError = errors.New("database error")
	AuthErrorInvalidBody          AuthError = errors.New("invalid body")
	AuthErrorInvalidPassword      AuthError = errors.New("wrong password provided")
	AuthErrorPasswordTooLong      AuthError = errors.New("the maximum password length is 72 bytes")
	AuthErrorPasswordTooShort     AuthError = errors.New("the minimum password length is 5 bytes")
	AuthErrorAdminAlreadyExists   AuthError = errors.New("there can only be one admin")
	AuthErrorTokenCreationFailure AuthError = errors.New("failed to create token")
	AuthErrorInvalidAuthHeader    AuthError = errors.New("invalid auth header provided")
	AuthErrorInvalidToken         AuthError = errors.New("invalid token provided")
	AuthErrorInvalidForm          AuthError = errors.New("invalid form")
	AuthErrorNotAuthorized        AuthError = errors.New("not authorized")
)

type RecipeError = error

var (
	RecipeErrorInvalidParam    RecipeError = errors.New("invalid path parameter")
	RecipeErrorInvalidFormData RecipeError = errors.New("invalid form data")
	RecipeErrorNotFound        RecipeError = errors.New("no recipe found")
	RecipeErrorDatabaseFailure RecipeError = errors.New("database error")
	RecipeErrorMarkdownFailure RecipeError = errors.New("error parsing markdown")
)

var ErrorHttpCodes = map[error]int{
	AuthErrorNoAdminFound:         http.StatusNotFound,
	AuthErrorDatabaseFailure:      http.StatusInternalServerError,
	AuthErrorInvalidBody:          http.StatusBadRequest,
	AuthErrorInvalidPassword:      http.StatusUnauthorized,
	AuthErrorPasswordTooLong:      http.StatusBadRequest,
	AuthErrorPasswordTooShort:     http.StatusBadRequest,
	AuthErrorAdminAlreadyExists:   http.StatusConflict,
	AuthErrorTokenCreationFailure: http.StatusInternalServerError,
	AuthErrorInvalidAuthHeader:    http.StatusBadRequest,
	AuthErrorInvalidToken:         http.StatusUnauthorized,
	AuthErrorInvalidForm:          http.StatusBadRequest,
	AuthErrorNotAuthorized:        http.StatusUnauthorized,

	RecipeErrorInvalidParam:    http.StatusBadRequest,
	RecipeErrorInvalidFormData: http.StatusBadRequest,
	RecipeErrorNotFound:        http.StatusNotFound,
	RecipeErrorDatabaseFailure: http.StatusInternalServerError,
}
