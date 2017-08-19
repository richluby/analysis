package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// builds a Transaction from a record line
func buildTransaction(record []string) (Transaction, error) {
	var err error
	trans := Transaction{}
	trans.Status = record[0]
	trans.Date, err = time.Parse("01/02/2006", record[2])
	if err != nil {
		return Transaction{}, err
	}
	trans.BusinessName = record[4]
	trans.Category = record[5]
	amount := strings.Replace(record[6], "--", "", 1)
	trans.Amount, err = strconv.ParseFloat(amount, 64)
	if err != nil {
		return Transaction{}, err
	}
	return trans, nil
}

// reads data from the specified file
// into the given channel
// sendChannel	: the channel to which to send data
func readData(filePaths []string, sendChannel chan Transaction) error {
	var err error
	defer close(sendChannel)
	for _, filePath := range filePaths{
		file, err := os.Open(filePath)
		defer file.Close()
		if err != nil{
			log.Printf("Failed to read file %s: %+v", filePath, err)
			continue
		}
		reader := csv.NewReader(file)
		for { // loop through all entries in file
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("Error reading %s: %+v", filePath, err)
			}
			trans, err := buildTransaction(record)
			if err != nil {
				log.Printf("Error parsing %s: %+v", filePath, err)
			}
			sendChannel <- trans
		}
	}
	return err
}

