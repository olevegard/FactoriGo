package main

type Inventory struct {
	iron_ore   int
	copper_ore int

	iron_plates   int
	copper_plates int
}

type Production struct {
	// Harvested
	iron_mines   ProductionUnit
	copper_mines ProductionUnit

	// Smelted
	iron_smelters   ProductionUnit
	copper_smelters ProductionUnit
}

type GameState struct {
	inventory  Inventory
	production Production
}

type DoRecipe func(Inventory, int) Inventory

func MakeProduction() Production {
	return Production{}
}

type ProductionUnit struct {
	count int

	ticks_remaining int // Time left until next batch of units
	ticks_per_cycle int // Time it takes to complete a full cycle

	doRecipe DoRecipe
}

func MakleProductionUnit(ticksPerCycle int, recipe DoRecipe) ProductionUnit {
	unit := ProductionUnit{}
	unit.count = 0
	unit.ticks_per_cycle = ticksPerCycle
	unit.ticks_remaining = ticksPerCycle

	unit.doRecipe = recipe

	return unit
}
