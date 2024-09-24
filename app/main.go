package main

import (
	"flag"

	"github.com/kilianmandscharo/lethimcook/auth"
	"github.com/kilianmandscharo/lethimcook/env"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/recipe"
	"github.com/kilianmandscharo/lethimcook/render"
	"github.com/kilianmandscharo/lethimcook/server"
)

func main() {
	env.LoadEnvironment(".env")

	logger := logging.New()
	renderer := render.New(&logger)

	authDatabase := auth.NewAuthDatabase(&logger)
	authService := auth.NewAuthService(authDatabase, &logger)
	authController := auth.NewAuthController(authService, &logger, &renderer)

	recipeDatabase := recipe.NewRecipeDatabase(&logger)
	recipeService := recipe.NewRecipeService(recipeDatabase, &logger)
	recipeController := recipe.NewRecipeController(recipeService, &logger, &renderer)

	var password = flag.String("init-admin", "", "the admin password")
	flag.Parse()
	authService.CreateAdminIfDoesNotExist(*password)

	server := server.New(authController, recipeController, &logger, &renderer)
	server.Start()
}
