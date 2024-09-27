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
	var password = flag.String("init-admin", "", "the admin password")
	var isProd = flag.Bool("prod", false, "shoudl the app run in production mode")
	flag.Parse()

    var logLevel logging.LoggerLevel
    if *isProd {
        logLevel = logging.Info
    } else {
        logLevel = logging.Debug
    }

	logger := logging.New(logLevel)
	renderer := render.New(&logger)

	env.LoadEnvironment(".env", &logger)

	authDatabase := auth.NewAuthDatabase(&logger)
	authService := auth.NewAuthService(authDatabase, &logger)
	authController := auth.NewAuthController(authService, &logger, &renderer)

	recipeDatabase := recipe.NewRecipeDatabase(&logger)
	recipeService := recipe.NewRecipeService(recipeDatabase, &logger)
	recipeController := recipe.NewRecipeController(recipeService, &logger, &renderer)

	authService.CreateAdminIfDoesNotExist(*password)
	server := server.New(authController, recipeController, &logger, &renderer, *isProd)
	server.Start()
}
