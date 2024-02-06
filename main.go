package main

import (
	"flag"
	"log"
	"os"

	"github.com/kilianmandscharo/lethimcook/auth"
)

func main() {
	password := flag.String("init-admin", "", "a secure password")
	flag.Parse()

	authorizer := auth.NewAuthorizer()

	if authorizer.DoesAdminExist() {
		if len(*password) != 0 {
			log.Println("Existing admin found, ignoring -init-admin")
		}
	} else {
		if len(*password) == 0 {
			log.Println("Create an admin before starting the server")
			log.Printf("Usage %s: -init-admin=<password>", os.Args[0])
			os.Exit(1)
		} else {
			authorizer.CreateAdmin(*password)
		}
	}

	// db := database.New()

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

	// server := server.New(&db)
	// server.Start()
}
