package main

type GameState struct {
	CurrentInventory  Inventory
	CurrentProduction Production
}

// Interface used for all printable things
type Printable interface {
	String() string
	Count() int
}
