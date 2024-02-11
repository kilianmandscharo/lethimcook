package recipe

import (
	"testing"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/stretchr/testify/assert"
)

func TestCreateRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()

	// When
	r := newTestRecipe()
	err := db.createRecipe(&r)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(1), r.ID)

	// When
	r = newTestRecipe()
	err = db.createRecipe(&r)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(2), r.ID)
}

func TestReadRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()

	// When
	_, err := db.readRecipe(uint(1))

	// Then
	assert.ErrorIs(t, err, errutil.RecipeErrorNotFound)

	// Given
	r := newTestRecipe()
	assert.NoError(t, db.createRecipe(&r))

	// When
	retrievedRecipe, err := db.readRecipe(r.ID)

	// Then
	assert.NoError(t, err)
	assertRecipesEqual(t, r, retrievedRecipe)
}

func TestReadAllRecipes(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()

	// When
	recipes, err := db.readAllRecipes()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 0, len(recipes))

	// Given
	r := newTestRecipe()
	assert.NoError(t, db.createRecipe(&r))

	// When
	recipes, err = db.readAllRecipes()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(recipes))
	assertRecipesEqual(t, r, recipes[0])
}

func TestDeleteRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()
	r := newTestRecipe()
	assert.NoError(t, db.createRecipe(&r))

	// When
	err := db.deleteRecipe(r.ID)

	// Then
	assert.NoError(t, err)

	// When
	_, err = db.readRecipe(r.ID)

	// Then
	assert.ErrorIs(t, err, errutil.RecipeErrorNotFound)

	// When
	err = db.deleteRecipe(r.ID)

	// Then
	assert.ErrorIs(t, err, errutil.RecipeErrorNotFound)
}

func TestUpdateRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()
	r := newTestRecipe()
	assert.NoError(t, db.createRecipe(&r))

	// When
	r.Title = "Test recipe title modified"
	err := db.updateRecipe(&r)
	retrievedRecipe, _ := db.readRecipe(r.ID)

	// Then
	assert.NoError(t, err)
	assertRecipesEqual(t, r, retrievedRecipe)
}
