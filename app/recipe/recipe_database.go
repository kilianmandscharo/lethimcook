package recipe

import (
	"errors"
	"fmt"
	"os"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type recipeDatabase struct {
	handler *gorm.DB
}

func newRecipeDatabase() recipeDatabase {
	db, err := gorm.Open(sqlite.Open("./recipe.db"), &gorm.Config{})

	if err != nil {
		fmt.Println("failed to connect recipe database: ", err)
		os.Exit(1)
	}

	db.AutoMigrate(&types.Recipe{})

	return recipeDatabase{handler: db}
}

func (db *recipeDatabase) createRecipe(recipe *types.Recipe) errutil.RecipeError {
	if err := db.handler.Create(recipe).Error; err != nil {
		return errutil.RecipeErrorDatabaseFailure
	}

	return nil
}

func (db *recipeDatabase) readRecipe(id uint) (types.Recipe, errutil.RecipeError) {
	var recipe types.Recipe

	if err := db.handler.First(&recipe, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return recipe, errutil.RecipeErrorNotFound
		}
		return recipe, errutil.RecipeErrorDatabaseFailure
	}

	return recipe, nil
}

func (db *recipeDatabase) readAllRecipes() ([]types.Recipe, errutil.RecipeError) {
	var recipes []types.Recipe

	if err := db.handler.Find(&recipes).Error; err != nil {
		return recipes, errutil.RecipeErrorDatabaseFailure
	}

	return recipes, nil
}

func (db *recipeDatabase) deleteRecipe(id uint) errutil.RecipeError {
	result := db.handler.Delete(&types.Recipe{}, id)

	if result.Error != nil {
		return errutil.RecipeErrorDatabaseFailure
	}

	if result.RowsAffected == 0 {
		return errutil.RecipeErrorNotFound
	}

	return nil
}

func (db *recipeDatabase) updateRecipe(recipe *types.Recipe) errutil.RecipeError {
	if err := db.handler.Save(recipe).Error; err != nil {
		return errutil.RecipeErrorDatabaseFailure
	}

	return nil
}
