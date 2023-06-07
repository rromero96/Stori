package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/olebedev/config"
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
	   YML Configuration
	*/
	file, err := ioutil.ReadFile(system.GetFileName("../conf", "production.yml"))
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, _ := config.ParseYaml(yamlString)

	/*
	   MYSQL client
	*/
	storiDBClient, err := createDBClient(getDBConnectionStringRoutes(storiDB, cfg))
	if err != nil {
		return err
	}

	/*
		Injections
	*/
	mysqlIDFinder := system.MakeMySQLFind(storiDBClient)
	mysqlCreateTransactions := system.MakeMySQLCreate(storiDBClient, mysqlIDFinder)
	readCSV := system.MakeReadCSV()
	htmlProcessTransactions := system.MakeHTMLProcessTransactions(readCSV, mysqlCreateTransactions)

	/*
		Endpoints
	*/
	app.Get(systemGetHtml, system.GetHTMLInfoV1(htmlProcessTransactions))

	/*
		Static Files serve
	*/
	app.Get(storyLogo, func(w http.ResponseWriter, r *http.Request) error {
		http.ServeFile(w, r, system.GetFileName(system.HtmlFolder, system.StoriLogo))
		return nil
	})

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

func getDBConnectionStringRoutes(database string, yml *config.Config) string {
	dbUserName, _ := yml.String(fmt.Sprintf("databases.mysql.%s.username", database))
	dbPassword, _ := yml.String(fmt.Sprintf("databases.mysql.%s.password", database))
	dbHost, _ := yml.String(fmt.Sprintf("databases.mysql.%s.db_host", database))
	dbName, _ := yml.String(fmt.Sprintf("databases.mysql.%s.db_name", database))
	return fmt.Sprintf(connectionStringFormat, dbUserName, dbPassword, dbHost, dbName)
}
