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
	var isProd = flag.Bool("prod", false, "should the app run in production mode")
	var shouldLogToFile = flag.Bool("log-file", false, "should the app log to a file")

	flag.Parse()

	var logLevel logging.LoggerLevel
	var logToFile bool
	if *isProd {
		logLevel = logging.Info
		logToFile = true
	} else {
		logLevel = logging.Debug
		logToFile = *shouldLogToFile
	}

	logger := logging.New(logLevel, logToFile)
	renderer := render.New(logger)

	env.LoadEnvironment(".env", logger)

	authDatabase := auth.NewAuthDatabase(logger)
	authService := auth.NewAuthService(authDatabase, logger)
	authController := auth.NewAuthController(authService, logger, renderer)

	recipeDatabase := recipe.NewRecipeDatabase(logger)
	recipeService := recipe.NewRecipeService(recipeDatabase, logger)
	recipeController := recipe.NewRecipeController(recipeService, logger, renderer)

	authService.CreateAdminIfDoesNotExist(*password)
	server := server.New(authController, recipeController, logger, renderer, *isProd)
	server.Start()
}
