// TODO:
// 1. Inventory items should be refered to by an Enum or something similar.
// 2. Find a way to represent change in inventory that can be easily converted to JSON
// 		i.e. uint inventoryId, int inventoryChange { 'inventory_changes': [{'inventory_id':3, 'inventory_change':-3}, {'inventory_id':2, 'inventory_change:5}]}
// 3. Have a function that takes a list of inventory changes and an inventory, that returns an inventory updated according to the list of inventory changes

/*
struct InvetoryItem
	count int
	name  string
	id  "iron_ore"

*/
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

// Production
// ============================================================================
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

	ticks_remaining ResetableInt // Time left until next batch of units
	ticks_per_cycle int          // Time it takes to complete a full cycle

	doRecipe DoRecipe // The function that updates the recipe
	name     string
}

func MakeProductionUnit(ticksPerCycle int, recipe DoRecipe, name string) ProductionUnit {
	unit := ProductionUnit{}
	unit.count = 0
	unit.ticks_per_cycle = ticksPerCycle
	unit.ticks_remaining = ResetableInt(ticksPerCycle)

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

func (inventoryItem InventoryItem) String() string {
	return inventoryItem.name
}

func (inventoryItem InventoryItem) Count() int {
	return inventoryItem.count
}

type ResetableInt int

func (value ResetableInt) ResetIfValue(conditionValue, resetValue int) (ResetableInt, bool) {
	if int(value) == conditionValue {
		return ResetableInt(resetValue), true
	}
	return value, false
}

type ProductionTimer struct {
	ticks_remaining int
	ticks_per_cycle int
}

func (unit ProductionUnit) MaybeResetTick() (ProductionUnit, bool) {
	unit.ticks_remaining--
	wasReset := false
	unit.ticks_remaining, wasReset = unit.ticks_remaining.ResetIfValue(0, unit.ticks_per_cycle)

	return unit, wasReset
}
