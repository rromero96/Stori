package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olebedev/config"

	"github.com/rromero96/stori/cmd/api/system"
)

const (
	systemGetHtml string = "/system/html/v1"

	connectionStringFormat string = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true"
	mysqlDriver            string = "mysql"
	storiDB                string = "stori"
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
	app := gin.Default()

	port := "8080"
	address := ":" + port

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
	app.GET(systemGetHtml, system.GetHTMLInfoV1(htmlProcessTransactions))

	log.Printf("server up and running in port %s", port)
	app.Run(address)
	return nil
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
