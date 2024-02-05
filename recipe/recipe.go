package recipe

type Recipe struct {
	ID          uint
	Title       string
	Description string
}

type Recipes = []Recipe

func New(title, description string) Recipe {
	return Recipe{Title: title, Description: description}
}

func (r *Recipe) Eq(other *Recipe) bool {
	return r.ID == other.ID && r.Title == other.Title && r.Description == other.Description
}
