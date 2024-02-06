package recipe

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RecipeDatabase struct {
	handler *gorm.DB
}

func NewTestRecipeDatabase() RecipeDatabase {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("failed to connect test database")
	}

	db.Migrator().DropTable(&Recipe{})

	db.AutoMigrate(&Recipe{})

	return RecipeDatabase{handler: db}
}

func NewRecipeDatabase() RecipeDatabase {
	db, err := gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect recipe database")
	}

	db.AutoMigrate(&Recipe{})

	return RecipeDatabase{handler: db}
}

func (db *RecipeDatabase) CreateRecipe(recipe *Recipe) error {
	if err := db.handler.Create(recipe).Error; err != nil {
		return err
	}

	return nil
}

func (db *RecipeDatabase) DeleteRecipe(id uint) error {
	result := db.handler.Delete(&Recipe{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("row with id=%d does not exist", id)
	}

	return nil
}

func (db *RecipeDatabase) UpdateRecipe(recipe *Recipe) error {
	if err := db.handler.Save(recipe).Error; err != nil {
		return err
	}

	return nil
}

func (db *RecipeDatabase) ReadRecipe(id uint) (Recipe, error) {
	var recipe Recipe

	if err := db.handler.First(&recipe, id).Error; err != nil {
		return recipe, err
	}

	return recipe, nil
}

func (db *RecipeDatabase) ReadAllRecipes() ([]Recipe, error) {
	var recipes []Recipe

	if err := db.handler.Find(&recipes).Error; err != nil {
		return recipes, err
	}

	return recipes, nil
}
