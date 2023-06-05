package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rromero96/roro-lib/cmd/config"
	"github.com/rromero96/roro-lib/cmd/web"

	"github.com/rromero96/stori/cmd/api/system"
)

const (
	systemGetInfo string = "/system/info/v1"
	systemGetHtml string = "/system/html/v1"
	storyLogo     string = "/static/stori_logo.jpeg"

	//this when its on docker
	connectionStringFormat string = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true"
	//connectionStringFormat string = "%s:%s@tcp/%s?charset=utf8&parseTime=true"
	mysqlDriver string = "mysql"
	storiDB     string = "stori"
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

	port := "8080"
	address := ":" + port
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	/*
	   MYSQL client
	*/
	storiDBClient, err := createDBClient(getDBConnectionStringRoutes(storiDB))
	if err != nil {
		return err
	}

	/*
		Injections
	*/
	mysqlCreateTransactions := system.MakeMySQLCreate(storiDBClient)
	readCSV := system.MakeReadCSV(mysqlCreateTransactions)
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

func createDBClient(connectionString string) (*sql.DB, error) {
	db, err := sql.Open(mysqlDriver, connectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(14 * time.Minute)

	return db, nil
}

func getDBConnectionStringRoutes(database string) string {
	dbUsername := config.String("databases", fmt.Sprintf("mysql.%s.username", database), "")
	dbPassword := config.String("databases", fmt.Sprintf("mysql.%s.password", database), "")
	dbHost := config.String("databases", fmt.Sprintf("mysql.%s.host", database), "")
	dbName := config.String("databases", fmt.Sprintf("mysql.%s.db_name", database), "")
	return fmt.Sprintf(connectionStringFormat, dbUsername, dbPassword, dbHost, dbName)
}
