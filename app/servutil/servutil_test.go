package servutil

import (
	"testing"

	"github.com/kilianmandscharo/lethimcook/testutil"
	"github.com/stretchr/testify/assert"
)

func TestIsHxRequest(t *testing.T) {
	c := testutil.NewEmptyTestContext(t)
	assert.False(t, IsHxRequest(c))

	c.Request().Header.Add("Hx-Request", "false")
	assert.False(t, IsHxRequest(c))

	c.Request().Header.Add("Hx-Request", "true")
	assert.False(t, IsHxRequest(c))

	c.Request().Header.Set("Hx-Request", "true")
	assert.True(t, IsHxRequest(c))
}

func TestIsAuthorized(t *testing.T) {
	c := testutil.NewEmptyTestContext(t)
	assert.False(t, IsAuthorized(c))

	c.Set("authorized", false)
	assert.False(t, IsAuthorized(c))

	c.Set("authorized", "true")
	assert.False(t, IsAuthorized(c))

	c.Set("authorized", true)
	assert.True(t, IsAuthorized(c))
}
