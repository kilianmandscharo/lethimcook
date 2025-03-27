package types

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/yuin/goldmark"
)

type Recipe struct {
	ID             uint   `json:"id"`
	Author         string `json:"author"`
	Source         string `json:"source"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Duration       int    `json:"duration"`
	Ingredients    string `json:"ingredients"`
	Instructions   string `json:"instructions"`
	Tags           string `json:"tags"`
	Pending        bool   `json:"-"`
	CreatedAt      string `json:"createdAt"`
	LastModifiedAt string `json:"-"`
	LastModifiedBy string `json:"-"`
}

func (r *Recipe) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("\n    id: ")
	buf.WriteString(strconv.Itoa(int(r.ID)))
	buf.WriteString("\n    author: ")
	buf.WriteString(r.Author)
	buf.WriteString("\n    source: ")
	buf.WriteString(r.Source)
	buf.WriteString("\n    title: ")
	buf.WriteString(r.Title)
	buf.WriteString("\n    description: ")
	buf.WriteString(r.Description)
	buf.WriteString("\n    duration: ")
	buf.WriteString(strconv.Itoa(r.Duration))
	buf.WriteString("\n    ingredients: ")
	buf.WriteString(r.Ingredients)
	buf.WriteString("\n    instructions: ")
	buf.WriteString(r.Instructions)
	buf.WriteString("\n    tags: ")
	buf.WriteString(r.Tags)
	buf.WriteString("\n    pending: ")
	buf.WriteString(strconv.FormatBool(r.Pending))
	buf.WriteString("\n    created at: ")
	buf.WriteString(r.CreatedAt)
	buf.WriteString("\n    last modified at: ")
	buf.WriteString(r.LastModifiedAt)
	return buf.String()
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
		strings.Contains(strings.ToLower(r.Author), query) ||
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

func (r *Recipe) RenderMarkdown() error {
	var buf bytes.Buffer

	if err := goldmark.Convert([]byte(r.Ingredients), &buf); err != nil {
		return &errutil.AppError{
			UserMessage: "Fehler beim Markdownparsing",
			StatusCode:  http.StatusInternalServerError,
			Err:         fmt.Errorf("failed at RenderMarkdown() with ingredients %s", r.Ingredients),
		}
	}
	r.Ingredients = buf.String()

	buf.Reset()

	if err := goldmark.Convert([]byte(r.Instructions), &buf); err != nil {
		return &errutil.AppError{
			UserMessage: "Fehler beim Markdownparsing",
			StatusCode:  http.StatusInternalServerError,
			Err:         fmt.Errorf("failed at RenderMarkdown() with instructions %s", r.Instructions),
		}
	}
	r.Instructions = buf.String()

	return nil
}

func (r *Recipe) ToVersion() RecipeVersion {
	return RecipeVersion{
		RecipeID:       r.ID,
		Author:         r.Author,
		Source:         r.Source,
		Title:          r.Title,
		Description:    r.Description,
		Duration:       r.Duration,
		Ingredients:    r.Ingredients,
		Instructions:   r.Instructions,
		Tags:           r.Tags,
		Pending:        r.Pending,
		CreatedAt:      r.CreatedAt,
		LastModifiedAt: r.LastModifiedAt,
		LastModifiedBy: r.LastModifiedBy,
	}
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
	FormElementInput FormElementType = iota + 1
	FormElementTextArea
)

type FormElement struct {
	Type           FormElementType
	Name           string
	Err            error
	Value          string
	InputType      string
	Label          string
	Required       bool
	Disabled       bool
	Placeholder    string
	LabelComponent templ.Component
}

func (f *FormElement) GetPlaceholder() string {
	if len(f.Placeholder) > 0 {
		return f.Placeholder
	}
	return f.Label
}

func (f *FormElement) GetLabel() string {
	if f.Required {
		return f.Label + "*"
	}
	return f.Label
}

type PaginationInfo struct {
	TotalRecipes int
	CurrentPage  int
	TotalPages   int
}

func (p *PaginationInfo) GetPageLabels() int {
	if p.TotalPages > 5 {
		return 5
	}
	return p.TotalPages
}

type RecipeLinkData struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type RecipeVersion struct {
	ID             uint   `gorm:"primaryKey"`
	RecipeID       uint   `json:"recipeId"`
	Author         string `json:"author"`
	Source         string `json:"source"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Duration       int    `json:"duration"`
	Ingredients    string `json:"ingredients"`
	Instructions   string `json:"instructions"`
	Tags           string `json:"tags"`
	Pending        bool   `json:"-"`
	CreatedAt      string `json:"createdAt"`
	LastModifiedAt string `json:"-"`
	LastModifiedBy string `json:"-"`
}
