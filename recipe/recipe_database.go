package recipe

import (
	"errors"
	"log"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type recipeDatabase struct {
	handler *gorm.DB
}

func newTestRecipeDatabase() recipeDatabase {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("failed to connect test database")
	}

	db.Migrator().DropTable(&recipe{})

	db.AutoMigrate(&recipe{})

	return recipeDatabase{handler: db}
}

func newRecipeDatabase() recipeDatabase {
	db, err := gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect recipe database")
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
