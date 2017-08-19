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
func initDB() {
	db, err := sql.Open("sqlite3", config.DatabaseFile)
	if err != nil {
		log.Fatalf("Error while opening database: %+v", err)
	}
	createTableStatement := fmt.Sprintf("create table %s (id integer not null primary key, Status text, BusinessName text, Category text, Date date, Amount real);", TABLE_NAME)
	_, err = db.Exec(createTableStatement)
	if err != nil && !strings.Contains(err.Error(), "table trans already exists") {
		log.Fatalf("Error while creating table: %+v", err)
	}
	database = db
}

// checks if this transaction is already present in
// the databse
// receiveChannel	: the channel with the raw transactions
// sendChannel		: the filtered channel to which to send new
//						transactions
func filterTransactions(receiveChannel chan Transaction, sendChannel chan Transaction) {
	defer close(sendChannel)
	querystatement, err := database.Prepare(fmt.Sprintf("select id, Status, BusinessName, Category, Date, Amount from %s where Date = ?", TABLE_NAME))
	if err != nil {
		log.Printf("Error while building query statement: %+v", err)
		return
	}
	defer querystatement.Close()
	transactionFound := false
	var checkTrans Transaction
	for trans := range receiveChannel {
		rows, err := querystatement.Query(trans.Date)
		if err != nil {
			log.Printf("Error while querying: %+v.", err)
			log.Printf("Dropping transaction: %+v", trans)
			continue
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&checkTrans.Id, &checkTrans.Status, &checkTrans.BusinessName, &checkTrans.Category, &checkTrans.Date, &checkTrans.Amount)
			if err != nil {
				log.Printf("Error while scanning row: %+v", err)
			}
			if checkTrans.Compare(trans) == 0 {
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
func enterData(receiveChannel chan Transaction) error {
	transactionsAdded := 0
	transactionHandler, err := database.Begin()
	if err != nil {
		log.Printf("Error while creating transaction handler: %+v", err)
		return err
	}
	filterChannel := make(chan Transaction, config.MaxBufferedTransactions)
	go filterTransactions(receiveChannel, filterChannel)
	insertStatement, err := transactionHandler.Prepare(fmt.Sprintf("insert into %s(Status, BusinessName, Category, Date, Amount) values(?, ?, ?, ?, ?)", TABLE_NAME))
	if err != nil {
		log.Printf("Error while creating insertion statement: %+v", err)
		return err
	}
	defer insertStatement.Close()
	for trans := range filterChannel {
		_, err = insertStatement.Exec(trans.Status, trans.BusinessName, trans.Category, trans.Date, trans.Amount)
		if err != nil {
			log.Printf("Error while inserting transaction: %+v", err)
			log.Printf("Dropping transaction: %+v", trans)
			continue
		}
		transactionsAdded += 1
	}
	transactionHandler.Commit()
	fmt.Printf("Transactions added: %d\n", transactionsAdded)
	return nil
}
