package main

type GameState struct {
	inventory  Inventory
	production Production
}

// Interface used for all printable things
type Printable interface {
	String() string
	Count() int
}
