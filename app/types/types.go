package types

import (
	"bytes"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/yuin/goldmark"
)

type Recipe struct {
	ID           uint
	Title        string
	Description  string
	Duration     int
	Ingredients  string
	Instructions string
}

func (r *Recipe) RenderMarkdown() errutil.RecipeError {
	var buf bytes.Buffer

	if err := goldmark.Convert([]byte(r.Ingredients), &buf); err != nil {
		return errutil.RecipeErrorMarkdownFailure
	}
	r.Ingredients = buf.String()

	buf.Reset()

	if err := goldmark.Convert([]byte(r.Instructions), &buf); err != nil {
		return errutil.RecipeErrorMarkdownFailure
	}
	r.Instructions = buf.String()

	return nil
}

func NewTestRecipe() Recipe {
	return Recipe{
		Title:        "Test title",
		Description:  "Test description",
		Duration:     30,
		Ingredients:  "Test ingredients",
		Instructions: "Test instructions",
	}
}
