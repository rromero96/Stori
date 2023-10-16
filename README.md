# Stori
This application was made by Rodrigo Romero using: 
- Golang Gin Gonic Framework
- Package oriented design architecture
- Unit Tests
- functional programming
- Amazon RDS database
- MySQL
- Postman

## Instructions for local usage
- Download the github repository from https://github.com/rromero96/Stori
- Open the terminal and write "go mod tidy"
- In the terminal place yourself in cmd/api and write "go run main.go"
- Open the browser and write this URL "http://localhost:8080/system/html/v1"

## Information
The application is conected to the RDS DB so you don't have to configure MYSQL in your local machine, in case you want to do that the SQL folder has the information about the db and you can set the credentials in production_test.yml. After that in main.go you have to make sure that the yaml that has to be used is the test one.

The folder automatedtests has the collection in it if you want to try in on Postman.

## Code Interface
- The application reads the csv file with the transactions information. The csv file is located in system/data
- Transforms this into a data type transactions and saves it into the database
- Then it prepares the information to make the email with the templated placed th system/html
- Renders the information email styled to the client