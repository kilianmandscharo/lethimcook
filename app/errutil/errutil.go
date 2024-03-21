package errutil

import (
	"errors"
	"net/http"
)

type AuthError = error

var (
	AuthErrorNoAdminFound         AuthError = errors.New("Kein Admin gefunden")
	AuthErrorDatabaseFailure      AuthError = errors.New("Datenbankfehler")
	AuthErrorInvalidBody          AuthError = errors.New("Ungültiger Body")
	AuthErrorInvalidPassword      AuthError = errors.New("Falsches Passwort")
	AuthErrorPasswordTooLong      AuthError = errors.New("Maximale Passwortlänge: 72")
	AuthErrorPasswordTooShort     AuthError = errors.New("Minimale Passwortlänge: 5")
	AuthErrorAdminAlreadyExists   AuthError = errors.New("Es kann nur einen Admin geben")
	AuthErrorTokenCreationFailure AuthError = errors.New("Fehler bei der Tokengenerierung")
	AuthErrorInvalidAuthHeader    AuthError = errors.New("Ungültiger Auth Header")
	AuthErrorInvalidToken         AuthError = errors.New("Ungültiger Token")
	AuthErrorInvalidForm          AuthError = errors.New("Ungültiges Formular")
	AuthErrorNotAuthorized        AuthError = errors.New("Nicht authorisiert")
)

type RecipeError = error

var (
	RecipeErrorInvalidParam    RecipeError = errors.New("Ungültiges Pfadparameter")
	RecipeErrorInvalidFormData RecipeError = errors.New("Ungültiges Formular")
	RecipeErrorNotFound        RecipeError = errors.New("Kein Rezept gefunden")
	RecipeErrorDatabaseFailure RecipeError = errors.New("Datenbankfehler")
	RecipeErrorMarkdownFailure RecipeError = errors.New("Fehler beim Markdownparsing")
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
