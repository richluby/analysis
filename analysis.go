// this package implements pretty printing of data
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// Set at compile time
var (
	Version string
	BuildTime string
)

var config ConfigCLI

// displays the help and exits the program
func printHelp(){
	fmt.Printf("\nData presenter, Version %s, Build %s\n", Version, BuildTime)
	fmt.Printf("Usage: %s [options] [csv_files]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

// initializes the command-line arguments for the progam
func initCLIArgs() ConfigCLI{
	config := ConfigCLI{}
	// informational args
	flag.Usage = printHelp
	flag.StringVar(&config.DatabaseFile, "database", "data/transaction.db","the sqlite3 file that contains the transactions")
	// date/time args
	defaultInterval := time.Duration(3) * time.Hour * 24
	flag.StringVar(&config.StartDate, "start-date", "", "the first date to include the display")
	flag.StringVar(&config.EndDate, "end-date", "", "the last date to include the display")
	flag.DurationVar(&config.Interval, "time-slice", defaultInterval, "the interval at which to aggregate disparate data points for graphing. accepts any format parsable by time.ParseDuration")
	// performance args
	flag.IntVar(&config.MaxBufferedTransactions, "max-buffer", 5000, "the maximum number of transactions to buffer in memory while reading CSV entries")
	// helper args
	flag.BoolVar(&config.Version, "version", false, "print version information and exit")
	flag.Parse()
	config.CSV = flag.Args()
	return config
}

// the entry point of the program
func main(){
	config = initCLIArgs()
	if config.Version {
		printHelp()
	}
	transactionChannel := make(chan Transaction, config.MaxBufferedTransactions)
	go readData(config.CSV, transactionChannel)
	initDB()
	enterData(transactionChannel)
}
