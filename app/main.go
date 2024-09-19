package main

import (
	"flag"

	"github.com/kilianmandscharo/lethimcook/auth"
	"github.com/kilianmandscharo/lethimcook/env"
	"github.com/kilianmandscharo/lethimcook/recipe"
	"github.com/kilianmandscharo/lethimcook/server"
)

func main() {
	env.LoadEnvironment(".env")

	authDatabase := auth.NewAuthDatabase()
	authService := auth.NewAuthService(authDatabase)
	authController := auth.NewAuthController(authService)

	recipeDatabase := recipe.NewRecipeDatabase()
	recipeService := recipe.NewRecipeService(recipeDatabase)
	recipeController := recipe.NewRecipeController(recipeService)

	var password = flag.String("init-admin", "", "the admin password")
	flag.Parse()
	authService.CreateAdminIfDoesNotExist(*password)

	server := server.New(authController, recipeController)
	server.Start()
}
