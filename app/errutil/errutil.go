package errutil

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	UserMessage string
	Err         error
	StatusCode  int
}

func (a *AppError) Error() string {
	return a.Err.Error()
}

func (a *AppError) Unwrap() error {
	return errors.Unwrap(a.Err)
}

func (a *AppError) AddMessage(message string) {
	a.Err = fmt.Errorf("%s: %w", message, a.Err)
}

func NewAppErrorNotAuthorized(functionName string) error {
	return &AppError{
		UserMessage: "Nicht authorisiert",
		Err:         errors.New(fmt.Sprintf("failed at %s, not authorized", functionName)),
		StatusCode:  http.StatusUnauthorized,
	}
}

func AddMessageToAppError(err error, message string) error {
	if appError, ok := err.(*AppError); ok {
		appError.AddMessage(message)
		return appError
	}
	return fmt.Errorf("%s: %w", message, err)
}

func GetAppErrorUserMessage(err error) string {
	if appError, ok := err.(*AppError); ok {
		return appError.UserMessage
	}
	return err.Error()
}

func GetAppErrorStatusCode(err error) int {
	if appError, ok := err.(*AppError); ok {
		return appError.StatusCode
	}
	return 0
}

var (
	FormErrorPasswortTooShort = errors.New("Minimale Passwortlänge: 5")
	FormErrorPasswortTooLong  = errors.New("Maximale Passwortlänge: 72")

	FormErrorNoTitle         = errors.New("Bitte trage einen Rezepttitel ein")
	FormErrorNoDescription   = errors.New("Bitte trage eine Rezeptbeschreibung ein")
	FormErrorNoDuration      = errors.New("Bitte trage die Zubereitungszeit ein")
	FormErrorNoIngredients   = errors.New("Bitte trage die Rezeptzutaten ein")
	FormErrorNoInstructions  = errors.New("Bitte trage die Rezeptanleitung ein")
	FormErrorInvalidPassword = errors.New("Falsches Passwort")
)
