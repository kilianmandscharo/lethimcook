package types

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
	recipe := NewTestRecipe()
	recipe.Ingredients = testMarkdownIngredients
	recipe.Instructions = testMarkdownInstructions

	// When
	err := recipe.RenderMarkdown()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, renderedIngredients, recipe.Ingredients)
	assert.Equal(t, renderedInstructions, recipe.Instructions)
}

func TestContainsQuery(t *testing.T) {
	testCases := []struct {
		recipe   Recipe
		query    string
		contains bool
	}{
		{
			recipe:   Recipe{Title: "Naan"},
			query:    "naan",
			contains: true,
		},
		{
			recipe:   Recipe{Description: "Indisches Fladenbrot – zubereitet in der Pfanne"},
			query:    "indisch",
			contains: true,
		},
		{
			recipe:   Recipe{Tags: "indisch, Beilage"},
			query:    "indisch",
			contains: true,
		},
		{
			recipe:   Recipe{Tags: "indisch, Beilage"},
			query:    "Beilage",
			contains: true,
		},
		{
			recipe: Recipe{
				Title:       "Naan",
				Description: "Indisches Fladenbrot – zubereitet in der Pfanne",
				Tags:        "indisch, Beilage",
			},
			query:    "test",
			contains: false,
		},
	}

	for _, test := range testCases {
		contains := test.recipe.ContainsQuery(test.query)
		assert.Equal(t, test.contains, contains)
	}
}

func TestParseTags(t *testing.T) {
	testCases := []struct {
		tags   string
		parsed []string
	}{
		{
			tags:   "",
			parsed: []string{},
		},
		{
			tags:   "Fleisch",
			parsed: []string{"Fleisch"},
		},
		{
			tags:   "Fleisch, ,",
			parsed: []string{"Fleisch"},
		},
		{
			tags:   "vegetarisch, asiatisch, scharf",
			parsed: []string{"vegetarisch", "asiatisch", "scharf"},
		},
		{
			tags:   "  vegetarisch  , asiatisch    , scharf  ",
			parsed: []string{"vegetarisch", "asiatisch", "scharf"},
		},
	}

	for _, test := range testCases {
		recipe := NewTestRecipe()
		recipe.Tags = test.tags
		parsed := recipe.ParseTags()
		assert.Equal(t, test.parsed, parsed)
	}
}
