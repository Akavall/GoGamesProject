package main

import ( "fmt"
	 "testing"
)

func TestInitializeDeck(t *testing.T) {
	deck := initialize_dec()

	for _, d := range deck {
		fmt.Println(d.name)
	}
}
