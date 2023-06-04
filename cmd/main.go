package main

import (
	"log"
	"net"
	"os"

	"github.com/rromero96/roro-lib/cmd/web"
	"github.com/rromero96/stori/cmd/api/system"
)

const (
	systemGetInfo string = "/system/info/v1"
	systemGetHtml string = "/system/html/v1"
	storyLogo     string = "/static/stori_logo.jpeg"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := ":" + port
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	/*
		Injections
	*/
	readCSV := system.MakeReadCSV()
	processTransactions := system.MakeProcessTransactions(readCSV)
	htmlProcessTransactions := system.MakeHTMLProcessTransactions(processTransactions)

	/*
		Endpoints
	*/
	app.Get(systemGetInfo, system.GetInfoV1(processTransactions))
	app.Get(systemGetHtml, system.GetHTMLInfoV1(htmlProcessTransactions))
	app.Get(storyLogo, system.GetLogoV1())

	log.Printf("server up and running in port %s", port)
	return web.Run(ln, web.DefaultTimeouts, app)
}
