package recipe

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestRecipeDatabase() *recipeDatabase {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("failed to connect test database: ", err)
		os.Exit(1)
	}
	db.Migrator().DropTable(&types.Recipe{})
	db.AutoMigrate(&types.Recipe{})
	return &recipeDatabase{handler: db, logger: logging.New(logging.Debug, false)}
}

func TestCreateRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()

	// When
	r := types.NewTestRecipe()
	err := db.createRecipe(&r)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(1), r.ID)

	// When
	r = types.NewTestRecipe()
	err = db.createRecipe(&r)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(2), r.ID)
}

func TestReadRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()

	// When
	_, err := db.readRecipe(uint(1))

	// Then
	assert.Error(t, err)
	appError, ok := err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, appError.StatusCode)
	assert.Equal(t, "Rezept nicht gefunden", appError.UserMessage)

	// Given
	r := types.NewTestRecipe()
	assert.NoError(t, db.createRecipe(&r))

	// When
	retrievedRecipe, err := db.readRecipe(r.ID)

	// Then
	assert.NoError(t, err)
	assert.True(t, r == retrievedRecipe)
}

func TestReadAllRecipes(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()

	// When
	recipes, err := db.readAllRecipesWithoutPending()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 0, len(recipes))

	// Given
	r := types.NewTestRecipe()
	assert.NoError(t, db.createRecipe(&r))

	// When
	recipes, err = db.readAllRecipesWithoutPending()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(recipes))
	assert.True(t, r == recipes[0])

	// Given
	r = types.NewTestRecipe()
	r.Pending = true
	assert.NoError(t, db.createRecipe(&r))

	// When
	recipes, err = db.readAllRecipesWithoutPending()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(recipes))
	assert.False(t, recipes[0].Pending)
}

func TestReadAllRecipesWithPending(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()

	// When
	recipes, err := db.readAllRecipesWithPending()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 0, len(recipes))

	// Given
	r := types.NewTestRecipe()
	assert.NoError(t, db.createRecipe(&r))

	// When
	recipes, err = db.readAllRecipesWithPending()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(recipes))
	assert.True(t, r == recipes[0])

	// Given
	r = types.NewTestRecipe()
	r.Pending = true
	assert.NoError(t, db.createRecipe(&r))

	// When
	recipes, err = db.readAllRecipesWithPending()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 2, len(recipes))
}

func TestDeleteRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()
	r := types.NewTestRecipe()
	assert.NoError(t, db.createRecipe(&r))

	// When
	err := db.deleteRecipe(r.ID)

	// Then
	assert.NoError(t, err)

	// When
	_, err = db.readRecipe(r.ID)

	// Then
	assert.Error(t, err)
	appError, ok := err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, appError.StatusCode)
	assert.Equal(t, "Rezept nicht gefunden", appError.UserMessage)

	// When
	err = db.deleteRecipe(r.ID)

	// Then
	assert.Error(t, err)
	appError, ok = err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, appError.StatusCode)
	assert.Equal(t, "Rezept nicht gefunden", appError.UserMessage)
}

func TestUpdateRecipe(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()
	r := types.NewTestRecipe()
	assert.NoError(t, db.createRecipe(&r))

	// When
	r.Title = "Test recipe title modified"
	err := db.updateRecipe(&r)
	retrievedRecipe, _ := db.readRecipe(r.ID)

	// Then
	assert.NoError(t, err)
	assert.True(t, r == retrievedRecipe)
}

func TestUpdatePending(t *testing.T) {
	// Given
	db := newTestRecipeDatabase()

	// When
	err := db.updatePending(1, true)

	// Then
	assert.Error(t, err)
	appError, ok := err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, appError.StatusCode)
	assert.Equal(t, "Rezept nicht gefunden", appError.UserMessage)

	// Given
	r := types.NewTestRecipe()
	assert.NoError(t, db.createRecipe(&r))

	// When
	err = db.updatePending(1, true)

	// Then
	assert.NoError(t, err)
}
