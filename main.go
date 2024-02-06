package main

import (
	"github.com/kilianmandscharo/lethimcook/server"
)

func main() {
	server := server.New()
	server.Start()
}
