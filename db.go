package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

var database *sql.DB
const TABLE_NAME = "trans"

// initializes the database and its connection
func initDB(){
	db, err := sql.Open("sqlite3", config.DatabaseFile)
	if err != nil {
		log.Fatalf("Error while opening database: %+v", err)
	}
	createTableStatement := fmt.Sprintf("create table %s (id integer not null primary key, Status text, BusinessName text, Category text, Date text, Amount real);", TABLE_NAME)
	_, err = db.Exec(createTableStatement)
	if err != nil && !strings.Contains(err.Error(), "table trans already exists"){
		log.Fatalf("Error while creating table: %+v", err)
	}
	database = db
}

}

// enters data into the database
// via the provided channel
func enterData(receiveChannel chan Transaction){
	for trans := range receiveChannel{
		fmt.Printf("Transaction: %+v\n", trans)
	}
}

