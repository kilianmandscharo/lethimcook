package errutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddMessageToAppError(t *testing.T) {
	t.Run("add to app error", func(t *testing.T) {
		// Given
		errorA := errors.New("error a")
		err := &AppError{Err: errorA}

		// When
		AddMessageToAppError(err, "error b")

		// Then
		assert.Equal(t, "error b: error a", err.Error())
		assert.True(t, errors.Is(errors.Unwrap(err), errorA))
	})

	t.Run("add to generic error", func(t *testing.T) {
		// Given
		errorA := errors.New("error a")

		// When
		err := AddMessageToAppError(errorA, "error b")

		// Then
		assert.Equal(t, "error b: error a", err.Error())
		assert.True(t, errors.Is(errors.Unwrap(err), errorA))
	})
}

func TestGetAppErrorUserMessage(t *testing.T) {
	assert.Equal(t, "test", GetAppErrorUserMessage(&AppError{UserMessage: "test"}))
	assert.Equal(t, "test", GetAppErrorUserMessage(errors.New("test")))
}

func TestGetAppErrorStatusCode(t *testing.T) {
	assert.Equal(t, 200, GetAppErrorStatusCode(&AppError{StatusCode: 200}))
	assert.Equal(t, 0, GetAppErrorStatusCode(errors.New("test")))
}
