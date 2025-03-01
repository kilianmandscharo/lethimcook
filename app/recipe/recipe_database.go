package recipe

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type recipeDatabase struct {
	handler *gorm.DB
	logger  *logging.Logger
}

func NewRecipeDatabase(logger *logging.Logger) *recipeDatabase {
	db, err := gorm.Open(sqlite.Open("./recipe.db"), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect recipe database: ", err)
	}
	db.AutoMigrate(&types.Recipe{}, &types.RecipeVersion{})
	return &recipeDatabase{handler: db, logger: logger}
}

func (db *recipeDatabase) createRecipe(recipe *types.Recipe) error {
	if err := db.handler.Create(recipe).Error; err != nil {
		return &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at createRecipe() with recipe %v, database failure: %w",
				*recipe,
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (db *recipeDatabase) readRecipe(id uint) (types.Recipe, error) {
	var recipe types.Recipe
	if err := db.handler.First(&recipe, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return recipe, &errutil.AppError{
				UserMessage: "Rezept nicht gefunden",
				Err: fmt.Errorf(
					"failed at readRecipe(), recipe with id %d not found",
					id,
				),
				StatusCode: http.StatusNotFound,
			}
		}
		return recipe, &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at readRecipe(), database failure: %w",
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return recipe, nil
}

func (db *recipeDatabase) readAllRecipesWithoutPending() ([]types.Recipe, error) {
	var recipes []types.Recipe
	if err := db.handler.Order("id desc").Where("pending IS NULL OR pending = 0").Find(&recipes).Error; err != nil {
		return recipes, &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at readAllRecipes(), database failure: %w",
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return recipes, nil
}

func (db *recipeDatabase) readAllRecipesWithPending() ([]types.Recipe, error) {
	var recipes []types.Recipe
	if err := db.handler.Order("id desc").Find(&recipes).Error; err != nil {
		return recipes, &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at readAllRecipes(), database failure: %w",
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return recipes, nil
}

func (db *recipeDatabase) deleteRecipe(id uint) error {
	result := db.handler.Delete(&types.Recipe{}, id)
	if err := result.Error; err != nil {
		return &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at deleteRecipe() with id %d, database failure: %w",
				id,
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	if result.RowsAffected == 0 {
		return &errutil.AppError{
			UserMessage: "Rezept nicht gefunden",
			Err: fmt.Errorf(
				"failed at deleteRecipe(), recipe with id %d not found",
				id,
			),
			StatusCode: http.StatusNotFound,
		}
	}
	return nil
}

func (db *recipeDatabase) updateRecipe(recipe *types.Recipe) error {
	if err := db.handler.Save(recipe).Error; err != nil {
		return &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at updateRecipe(), database failure: %w",
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (db *recipeDatabase) updatePending(id uint, pending bool) error {
	createError := func(err error) error {
		return errutil.AddMessageToAppError(
			err,
			fmt.Sprintf("failed at updatePending() with id %d", id),
		)
	}
	recipe, err := db.readRecipe(id)
	if err != nil {
		return createError(err)
	}
	recipe.Pending = pending
	err = db.updateRecipe(&recipe)
	if err != nil {
		return createError(err)
	}
	return nil
}

func (db *recipeDatabase) createRecipeVersion(recipeVersion *types.RecipeVersion) error {
	if err := db.handler.Create(recipeVersion).Error; err != nil {
		return &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at createRecipeVersion() with recipe %v, database failure: %w",
				*recipeVersion,
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}

func (db *recipeDatabase) readRecipeVersionsForRecipe(id uint) ([]types.RecipeVersion, error) {
	var recipeVersions []types.RecipeVersion
	if err := db.handler.Order("version_id desc").Where("recipe_id = ?", id).Find(&recipeVersions).Error; err != nil {
		return recipeVersions, &errutil.AppError{
			UserMessage: "Datenbankfehler",
			Err: fmt.Errorf(
				"failed at readRecipeVersionsForRecipe(), database failure: %w",
				err,
			),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return recipeVersions, nil
}
