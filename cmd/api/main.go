package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olebedev/config"

	"github.com/rromero96/stori/cmd/api/system"
)

const (
	systemGetHtml          string = "/system/html/v1"
	storyLogo              string = "/static/stori_logo.jpeg"
	connectionStringFormat string = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true"
	mysqlDriver            string = "mysql"
	storiDB                string = "stori"
)

func main() {
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
		panic(err)
	}

	/*
		Injections
	*/
	mysqlIDFinder := system.MakeMySQLFind(storiDBClient)
	mysqlCreateTransactions := system.MakeMySQLCreate(storiDBClient, mysqlIDFinder)
	readCSV := system.MakeReadCSV()
	htmlProcessTransactions := system.MakeHTMLProcessTransactions(readCSV, mysqlCreateTransactions)

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return system.GetHTMLInfoV1(ctx, request, htmlProcessTransactions)
	})
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
