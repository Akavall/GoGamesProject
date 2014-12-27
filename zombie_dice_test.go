package main

import ( "fmt"
	 "testing"

	 "github.com/Akavall/GoGamesProject/dice"
)

func TestInitializeDeck(t *testing.T) {
	deck := initialize_dec()

	for _, d := range deck {
		fmt.Println(d.name)
	}
}
