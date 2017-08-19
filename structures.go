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
	LogLevel 	string
	CSV			[]string
	Version 	bool
	Interval 	time.Duration
	MaxBufferedTransactions int
}

