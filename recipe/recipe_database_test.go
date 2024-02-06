package recipe

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRecipe(t *testing.T) {
	// Given
	db := NewTestRecipeDatabase()

	// When
	r := New("First test recipe", "Test description")
	err := db.CreateRecipe(&r)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(1), r.ID)

	// When
	r = New("Second test recipe", "Test description")
	err = db.CreateRecipe(&r)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(2), r.ID)
}

func TestReadRecipe(t *testing.T) {
	// Given
	db := NewTestRecipeDatabase()
	r := Recipe{Title: "Test recipe", Description: "Test description"}
	db.CreateRecipe(&r)

	// When
	retrievedRecipe, err := db.ReadRecipe(r.ID)

	// Then
	assert.NoError(t, err)
	assert.True(t, retrievedRecipe.Eq(&r))
}

func TestReadAllRecipes(t *testing.T) {
	// Given
	db := NewTestRecipeDatabase()

	// When
	recipes, err := db.ReadAllRecipes()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 0, len(recipes))

	// Given
	r := New("Test recipe", "Test description")
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
	db := NewTestRecipeDatabase()
	r := New("Test recipe", "Test description")
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
	db := NewTestRecipeDatabase()
	r := New("Test recipe", "Test description")
	db.CreateRecipe(&r)

	// When
	r.Title = "Test recipe title modified"
	err := db.UpdateRecipe(&r)
	retrievedRecipe, _ := db.ReadRecipe(r.ID)

	// Then
	assert.NoError(t, err)
	assert.True(t, retrievedRecipe.Eq(&r))
}
