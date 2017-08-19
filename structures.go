package main

import "time"

// stores details of a transaction
type Transaction struct{
	Status,
	BusinessName,
	Category string
	Date 	time.Time
	Amount 	float64
	Id		int
}

// stores command line arguments
type ConfigCLI struct{
	StartDate,
	EndDate,
	DatabaseFile,
	OutputDirectory,
	GraphTypes,
	Categories,
	LogLevel 	string
	CSV			[]string
	Version 	bool
	Interval 	time.Duration
	MaxBufferedTransactions int
}

// implements the Compare for Transactions
func (a Transaction) Compare(b Transaction) int {
	if a.Date.Before(b.Date){
		return 1
	} else if a.Date.After(b.Date) {
		return -1
	}
	if a.BusinessName > b.BusinessName{
		return 1
	} else if a.BusinessName < b.BusinessName{
		return -1
	}
	if a.Amount > b.Amount {
		return 1
	} else if a.Amount < b.Amount {
		return -1
	}
	return 0
}

