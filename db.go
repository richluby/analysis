package main

import "fmt"

// initializes the database and its connection
func initDB(){

}

// enters data into the database
// via the provided channel
func enterData(receiveChannel chan Transaction){
	for trans := range receiveChannel{
		fmt.Printf("Transaction: %+v\n", trans)
	}
}

