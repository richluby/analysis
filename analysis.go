// this package implements pretty printing of data
package main

import (
	"fmt"
	"os"
)

// Set at compile time
var (
	Version string
	BuildTime string
)

func printHelp(){
	fmt.Printf("Data presenter, Version %s, Build %s\n", Version, BuildTime)
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
}

// the entry point of the program
func main(){
	printHelp()
	transactionChannel := make(chan Transaction)
	if err := readData(transactionChannel); err != nil{
		fmt.Printf("Error while reading data: %+v", err)
		os.Exit(1)
	}
}

