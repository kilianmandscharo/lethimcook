package recipe

import (
	"bytes"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/yuin/goldmark"
)

type recipe struct {
	ID           uint
	Title        string
	Description  string
	Duration     int
	Ingredients  string
	Instructions string
}

type recipes = []recipe

func (r *recipe) renderMarkdown() errutil.RecipeError {
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
