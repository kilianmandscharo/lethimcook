package recipe

type recipe struct {
	ID           uint
	Title        string
	Description  string
	Duration     int
	Ingredients  string
	Instructions string
}

type recipes = []recipe

func newRecipe(title, description string) recipe {
	return recipe{Title: title, Description: description}
}

func (r *recipe) eq(other *recipe) bool {
	return r.ID == other.ID && r.Title == other.Title && r.Description == other.Description
}
