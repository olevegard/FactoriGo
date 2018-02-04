package main

type Inventory struct {
	iron_ore   int
	copper_ore int

	iron_plates   int
	copper_plates int
}

type Production struct {
	// Harvested
	iron_mines   int
	copper_mines int

	// Smelted
	iron_smelters   int
	copper_smelters int
}

func MakeProduction() Production {
	return Production{0, 0, 0, 0}
}
