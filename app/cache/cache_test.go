package cache

import (
	"testing"

	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/stretchr/testify/assert"
)

func newTestRecipeCache() *RecipeCache {
	return NewRecipeCache(logging.New(logging.Debug, false))
}

func TestSet(t *testing.T) {
	c := newTestRecipeCache()
	c.Set([]types.Recipe{
		{ID: 1},
		{ID: 2},
	})
	assert.Equal(t, []types.Recipe{
		{ID: 1},
		{ID: 2},
	}, c.recipes)
}

func TestInvalidate(t *testing.T) {
	c := newTestRecipeCache()
	c.Set([]types.Recipe{
		{ID: 1},
		{ID: 2},
	})
	c.Invalidate()
	assert.Nil(t, c.recipes)
}

func TestIsValid(t *testing.T) {
	c := newTestRecipeCache()
	assert.False(t, c.IsValid())
	c.Set([]types.Recipe{
		{ID: 1},
		{ID: 2},
	})
	assert.True(t, c.IsValid())
	c.Invalidate()
	assert.False(t, c.IsValid())
}

func TestGet(t *testing.T) {
	c := newTestRecipeCache()
	assert.Equal(t, []types.Recipe{}, *c.Get(true))
	assert.Equal(t, []types.Recipe{}, *c.Get(false))
	c.Set([]types.Recipe{
		{ID: 1},
		{ID: 2, Pending: true},
	})
	assert.Equal(t, []types.Recipe{
		{ID: 1},
		{ID: 2, Pending: true},
	}, *c.Get(true))
	assert.Equal(t, []types.Recipe{
		{ID: 1},
	}, *c.Get(false))
}
