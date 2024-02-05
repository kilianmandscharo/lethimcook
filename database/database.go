package database

import (
	"fmt"
	"log"

	"github.com/kilianmandscharo/lethimcook/recipe"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	handler *gorm.DB
}

func NewTestDatabase() Database {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("failed to connect test database")
	}

	db.Migrator().DropTable(&recipe.Recipe{})

	db.AutoMigrate(&recipe.Recipe{})

	return Database{handler: db}
}

func New() Database {
	db, err := gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")
	}

	db.AutoMigrate(&recipe.Recipe{})

	return Database{handler: db}
}

func (db *Database) CreateRecipe(recipe *recipe.Recipe) error {
	if err := db.handler.Create(recipe).Error; err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteRecipe(id uint) error {
	result := db.handler.Delete(&recipe.Recipe{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("row with id=%d does not exist", id)
	}

	return nil
}

func (db *Database) UpdateRecipe(recipe *recipe.Recipe) error {
	if err := db.handler.Save(recipe).Error; err != nil {
		return err
	}

	return nil
}

func (db *Database) ReadRecipe(id uint) (recipe.Recipe, error) {
	var recipe recipe.Recipe

	if err := db.handler.First(&recipe, id).Error; err != nil {
		return recipe, err
	}

	return recipe, nil
}

func (db *Database) ReadAllRecipes() ([]recipe.Recipe, error) {
	var recipes []recipe.Recipe

	if err := db.handler.Find(&recipes).Error; err != nil {
		return recipes, err
	}

	return recipes, nil
}
