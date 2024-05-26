package types

import (
	"bytes"
	"strings"

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
	Tags         string
}

func (r *Recipe) ParseTags() []string {
	if len(strings.TrimSpace(r.Tags)) == 0 {
		return []string{}
	}

	tags := []string{}

	for _, tag := range strings.Split(r.Tags, ",") {
		trimmedTag := strings.TrimSpace(tag)
		if len(trimmedTag) > 0 {
			tags = append(tags, trimmedTag)
		}
	}

	return tags
}

func (r *Recipe) ContainsQuery(query string) bool {
	query = strings.ToLower(query)

	return strings.Contains(strings.ToLower(r.Title), query) ||
		strings.Contains(strings.ToLower(r.Description), query) ||
		r.containsQueryInTags(query)
}

func (r *Recipe) containsQueryInTags(query string) bool {
	for _, tag := range r.ParseTags() {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}

	return false
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

type FormElementType int

const (
	FormElementInput    FormElementType = 1
	FormElementTextArea FormElementType = 2
)

type FormElement struct {
	Type      FormElementType
	Name      string
	Err       error
	Value     string
	InputType string
	Label     string
	Required  bool
}
