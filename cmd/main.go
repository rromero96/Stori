package main

import (
	"log"
	"net"

	"github.com/rromero96/stori/internal/web"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	/*
		Server Configuration
	*/
	app := web.New()

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	return web.Run(ln, web.DefaultTimeouts, app)
}
