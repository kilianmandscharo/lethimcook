package recipe

import (
	"errors"
	"fmt"
	"os"

	"github.com/kilianmandscharo/lethimcook/errutil"
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

	db.AutoMigrate(&recipe{})

	return recipeDatabase{handler: db}
}

func (db *recipeDatabase) createRecipe(recipe *recipe) errutil.RecipeError {
	if err := db.handler.Create(recipe).Error; err != nil {
		return errutil.RecipeErrorDatabaseFailure
	}

	return nil
}

func (db *recipeDatabase) readRecipe(id uint) (recipe, errutil.RecipeError) {
	var recipe recipe

	if err := db.handler.First(&recipe, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return recipe, errutil.RecipeErrorNotFound
		}
		return recipe, errutil.RecipeErrorDatabaseFailure
	}

	return recipe, nil
}

func (db *recipeDatabase) readAllRecipes() ([]recipe, errutil.RecipeError) {
	var recipes []recipe

	if err := db.handler.Find(&recipes).Error; err != nil {
		return recipes, errutil.RecipeErrorDatabaseFailure
	}

	return recipes, nil
}

func (db *recipeDatabase) deleteRecipe(id uint) errutil.RecipeError {
	result := db.handler.Delete(&recipe{}, id)

	if result.Error != nil {
		return errutil.RecipeErrorDatabaseFailure
	}

	if result.RowsAffected == 0 {
		return errutil.RecipeErrorNotFound
	}

	return nil
}

func (db *recipeDatabase) updateRecipe(recipe *recipe) errutil.RecipeError {
	if err := db.handler.Save(recipe).Error; err != nil {
		return errutil.RecipeErrorDatabaseFailure
	}

	return nil
}
