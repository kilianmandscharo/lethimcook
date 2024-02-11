package recipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestRecipe() recipe {
	return recipe{
		Title:        "Test title",
		Description:  "Test description",
		Duration:     30,
		Ingredients:  "Test ingredients",
		Instructions: "Test instructions",
	}
}

func assertRecipesEqual(t *testing.T, first, second recipe) {
	assert.True(
		t,
		first.ID == second.ID &&
			first.Description == second.Description &&
			first.Duration == second.Duration &&
			first.Ingredients == second.Ingredients &&
			first.Instructions == second.Instructions,
	)
}
