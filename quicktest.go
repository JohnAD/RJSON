package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	interestRate, err := bson.ParseDecimal128("3.0200") // accurate to 4 decimal places
	if err != nil {
		panic(err)
	}
	fmt.Println("val: ", interestRate) // val: 3.0200
}
