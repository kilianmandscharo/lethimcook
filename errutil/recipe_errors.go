package errutil

import (
	"errors"
	"net/http"
)

type RecipeError = error

var (
	RecipeErrorInvalidParam    RecipeError = errors.New("invalid path parameter")
	RecipeErrorInvalidFormData RecipeError = errors.New("invalid form data")
	RecipeErrorNotFound        RecipeError = errors.New("no recipe found")
	RecipeErrorDatabaseFailure RecipeError = errors.New("database error")
)

var RecipeErrorHttpCodes = map[RecipeError]int{
	RecipeErrorInvalidParam:    http.StatusBadRequest,
	RecipeErrorInvalidFormData: http.StatusBadRequest,
	RecipeErrorNotFound:        http.StatusNotFound,
	RecipeErrorDatabaseFailure: http.StatusInternalServerError,
}
