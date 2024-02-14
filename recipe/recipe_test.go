package recipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testMarkdownIngredients = "- Ingredient 1\n- Ingredient 2\n- Ingredient 3"
const testMarkdownInstructions = "1. Instruction 1\n2. Instruction 2\n3. Instruction 3"

const renderedIngredients = "<ul>\n<li>Ingredient 1</li>\n<li>Ingredient 2</li>\n<li>Ingredient 3</li>\n</ul>\n"
const renderedInstructions = "<ol>\n<li>Instruction 1</li>\n<li>Instruction 2</li>\n<li>Instruction 3</li>\n</ol>\n"

func TestRenderMarkdown(t *testing.T) {
	// Given
	recipe := newTestRecipe()
	recipe.Ingredients = testMarkdownIngredients
	recipe.Instructions = testMarkdownInstructions

	// When
	err := recipe.renderMarkdown()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, renderedIngredients, recipe.Ingredients)
	assert.Equal(t, renderedInstructions, recipe.Instructions)
}
