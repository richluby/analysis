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

// checks if this transaction is already present in
// the databse
// receiveChannel	: the channel with the raw transactions
// sendChannel		: the filtered channel to which to send new 
//						transactions
func filterTransactions(receiveChannel chan Transaction, sendChannel chan Transaction){
	defer close(sendChannel)
	querystatement, err := database.Prepare(fmt.Sprintf("select * from %s where date = ?", TABLE_NAME))
	if err != nil{
		log.Printf("Error while building query statement: %+v", err)
		return
	}
	defer querystatement.Close()
	transactionFound := false
	var checkTrans Transaction
	for trans := range receiveChannel{
		rows, err := querystatement.Query(trans.Date)
		if err != nil{
			log.Printf("Error while querying: %+v.", err)
			log.Printf("Dropping transaction: %+v", trans)
			continue
		}
		for rows.Next() {
			rows.Scan(&checkTrans)
			if checkTrans == trans {
				log.Printf("Duplicate transaction: %+v", checkTrans)
				transactionFound = true
				break
			}
		}
		if !transactionFound {
			sendChannel <- trans
		} else {
			transactionFound = false
		}
	}
}

// enters data into the database
// via the provided channel
func enterData(receiveChannel chan Transaction){
	for trans := range receiveChannel{
		fmt.Printf("Transaction: %+v\n", trans)
	}
}

