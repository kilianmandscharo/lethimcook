package recipe

import (
	"os"
	"path"
)

type recipeMarkdownLoader struct {
	dir string
}

func newRecipeMarkdownLoader() recipeMarkdownLoader {
	return recipeMarkdownLoader{
		dir: "markdown",
	}
}

func (l *recipeMarkdownLoader) writeRecipe(fileName, markdown string) error {
	return os.WriteFile(l.getFilePath(fileName), []byte(markdown), 0644)
}

func (l *recipeMarkdownLoader) readRecipe(name string) (string, error) {
	markdown, err := os.ReadFile(l.getFilePath(name))

	if err != nil {
		return "", err
	}

	return string(markdown), nil
}

func (l *recipeMarkdownLoader) getFilePath(fileName string) string {
	return path.Join(l.dir, fileName)
}
