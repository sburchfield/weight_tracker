package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DBConnection holds the global DB connection
type dbConnection struct {
	db *gorm.DB
}

// gloabl db connection object
var dbConn dbConnection

//InitDB Initialize db related stuff
func (dbc *dbConnection) initDB() {

	conn := "host=" + envVars.dbHost +
		" user=" + envVars.dbUsername +
		" dbname=" + envVars.dbName +
		" sslmode=" + envVars.dbSSLMode +
		" password=" + envVars.dbPassword

	var err error

	dbc.db, err = gorm.Open("postgres", conn)
	if err != nil {
		log.Fatalf("Error connecting to database, %v", err)
	}

	dbc.db.LogMode(true)

	fmt.Println("-- DB connected successfully --")

}

// Close the global resources
func (dbc *dbConnection) close() error {

	return dbc.db.Close()

}
