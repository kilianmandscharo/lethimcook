package main

import (
	"flag"

	"github.com/kilianmandscharo/lethimcook/auth"
	"github.com/kilianmandscharo/lethimcook/server"
)

func main() {
	var password = flag.String("init-admin", "", "the admin password")
	flag.Parse()

	authService := auth.NewAuthService()
	authService.CreateAdminIfDoesNotExist(*password)

	server := server.New(authService)
	server.Start()
}
