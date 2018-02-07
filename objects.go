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

type Inventory struct {
	iron_ore   InventoryItem
	copper_ore InventoryItem

	iron_plates   InventoryItem
	copper_plates InventoryItem
}

type Production struct {
	// Harvested
	iron_mines   ProductionUnit
	copper_mines ProductionUnit

	// Smelted
	iron_smelters   ProductionUnit
	copper_smelters ProductionUnit
}

// Used to update Inventory for a recipe
type DoRecipe func(Inventory, int) Inventory

// Any thing that can produce/mine/smelt etc...
type ProductionUnit struct {
	count int

	ticks_remaining int // Time left until next batch of units
	ticks_per_cycle int // Time it takes to complete a full cycle

	doRecipe DoRecipe // The function that updates the recipe
	name     string
}

func MakeProductionUnit(ticksPerCycle int, recipe DoRecipe, name string) ProductionUnit {
	unit := ProductionUnit{}
	unit.count = 0
	unit.ticks_per_cycle = ticksPerCycle
	unit.ticks_remaining = ticksPerCycle

	unit.doRecipe = recipe
	unit.name = name

	return unit
}

func (production ProductionUnit) String() string {
	return production.name
}

func (production ProductionUnit) Count() int {
	return production.count
}

type InventoryItem struct {
	count int
	name  string
}

func (inventoryItem InventoryItem) String() string {
	return inventoryItem.name
}

func (inventoryItem InventoryItem) Count() int {
	return inventoryItem.count
}
