package main

import (
	// "log"

	"github.com/kilianmandscharo/lethimcook/database"
	// "github.com/kilianmandscharo/lethimcook/recipe"
	"github.com/kilianmandscharo/lethimcook/server"
)

func main() {
	db := database.New()

	// recipes := []recipe.Recipe{
	// 	recipe.New("Spaghetti aglio e olio", "empty"),
	// 	recipe.New("Masala Chai", "empty"),
	// 	recipe.New("Butter Chicken", "empty"),
	// }
	//
	// for _, recipe := range recipes {
	// 	err := db.CreateRecipe(&recipe)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	server := server.New(&db)
	server.Start()
}
