package main

import (
	"fmt"

	"github.com/fofoRS/go-tutorial/own_deck"
)

func main() {
	deck := own_deck.New()
	for _, card := range deck {
		fmt.Printf("(%s,%s)\n", card.Family.String(), card.Name.String())
	}
}
