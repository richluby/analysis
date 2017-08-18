// this package implements pretty printing of data
package main

import (
	"fmt"
	"os"
)

// the entry point of the program
func main(){
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	readData()
}

