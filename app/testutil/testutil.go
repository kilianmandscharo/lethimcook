package testutil

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type RequestOptions struct {
	HandlerFunc    func(c echo.Context) error
	Method         string
	Route          string
	StatusWant     int
	Authorized     bool
	WithFormData   bool
	FormData       string
	WithCookie     bool
	Cookie         http.Cookie
	WithPathParam  bool
	PathParamName  string
	PathParamValue string
	WithQueryParam bool
	QueryParam     string
}

func AssertRequest(t *testing.T, options RequestOptions) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	rr := httptest.NewRecorder()

	var body io.Reader

	if options.WithFormData {
		body = bytes.NewBufferString(options.FormData)
	} else {
		body = nil
	}

	if options.WithQueryParam {
		options.Route += options.QueryParam
	}

	req, err := http.NewRequest(options.Method, options.Route, body)
	assert.NoError(t, err)

	if options.WithFormData {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	}

	c := e.NewContext(req, rr)

	if options.Authorized {
		c.Set("authorized", true)
	} else {
		c.Set("authorized", false)
	}

	if options.WithPathParam {
		c.SetParamNames(options.PathParamName)
		c.SetParamValues(options.PathParamValue)
	}

	assert.NoError(t, options.HandlerFunc(c))

	assert.Equal(t, options.StatusWant, rr.Code)

	return rr, c
}
