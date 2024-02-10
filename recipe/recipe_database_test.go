package recipe

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()

	// When
	r := newRecipe("First test recipe", "Test description")
	err := db.createRecipe(&r)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(1), r.ID)

	// When
	r = newRecipe("Second test recipe", "Test description")
	err = db.createRecipe(&r)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(2), r.ID)
}

func TestReadRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()
	r := recipe{Title: "Test recipe", Description: "Test description"}
	db.createRecipe(&r)

	// When
	retrievedRecipe, err := db.readRecipe(r.ID)

	// Then
	assert.NoError(t, err)
	assert.True(t, retrievedRecipe.eq(&r))
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
	r := newRecipe("Test recipe", "Test description")
	db.createRecipe(&r)

	// When
	recipes, err = db.readAllRecipes()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(recipes))
	assert.True(t, recipes[0].eq(&r))
}

func TestDeleteRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()
	r := newRecipe("Test recipe", "Test description")
	db.createRecipe(&r)

	// When
	err := db.deleteRecipe(r.ID)

	// Then
	assert.NoError(t, err)

	// When
	_, err = db.readRecipe(r.ID)

	// Then
	assert.Error(t, err)

	// When
	err = db.deleteRecipe(r.ID)

	// Then
	assert.Error(t, err)
}

func TestUpdateRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()
	r := newRecipe("Test recipe", "Test description")
	db.createRecipe(&r)

	// When
	r.Title = "Test recipe title modified"
	err := db.updateRecipe(&r)
	retrievedRecipe, _ := db.readRecipe(r.ID)

	// Then
	assert.NoError(t, err)
	assert.True(t, retrievedRecipe.eq(&r))
}
