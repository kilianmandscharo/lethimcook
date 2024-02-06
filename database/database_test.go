package database

import (
	"testing"

	"github.com/kilianmandscharo/lethimcook/recipe"
	"github.com/stretchr/testify/assert"
)

func TestCreateRecipe(t *testing.T) {
	// Given
	db := NewTestDatabase()

	// When
	r := recipe.New("First test recipe", "Test description")
	err := db.CreateRecipe(&r)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(1), r.ID)

	// When
	r = recipe.New("Second test recipe", "Test description")
	err = db.CreateRecipe(&r)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(2), r.ID)
}

func TestReadRecipe(t *testing.T) {
	// Given
	db := NewTestDatabase()
	r := recipe.Recipe{Title: "Test recipe", Description: "Test description"}
	db.CreateRecipe(&r)

	// When
	retrievedRecipe, err := db.ReadRecipe(r.ID)

	// Then
	assert.NoError(t, err)
	assert.True(t, retrievedRecipe.Eq(&r))
}

func TestReadAllRecipes(t *testing.T) {
	// Given
	db := NewTestDatabase()

	// When
	recipes, err := db.ReadAllRecipes()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 0, len(recipes))

	// Given
	r := recipe.New("Test recipe", "Test description")
	db.CreateRecipe(&r)

	// When
	recipes, err = db.ReadAllRecipes()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(recipes))
	assert.True(t, recipes[0].Eq(&r))
}

func TestDeleteRecipe(t *testing.T) {
	// Given
	db := NewTestDatabase()
	r := recipe.New("Test recipe", "Test description")
	db.CreateRecipe(&r)

	// When
	err := db.DeleteRecipe(r.ID)

	// Then
	assert.NoError(t, err)

	// When
	_, err = db.ReadRecipe(r.ID)

	// Then
	assert.Error(t, err)

	// When
	err = db.DeleteRecipe(r.ID)

	// Then
	assert.Error(t, err)
}

func TestUpdateRecipe(t *testing.T) {
	// Given
	db := NewTestDatabase()
	r := recipe.New("Test recipe", "Test description")
	db.CreateRecipe(&r)

	// When
	r.Title = "Test recipe title modified"
	err := db.UpdateRecipe(&r)
	retrievedRecipe, _ := db.ReadRecipe(r.ID)

	// Then
	assert.NoError(t, err)
	assert.True(t, retrievedRecipe.Eq(&r))
}
