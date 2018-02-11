package main

import (
	"testing"

	"github.com/stvp/assert"
)

func MakeNewProductionUnitWithNoBuildNew(ticksPerCycle int, name string, recipeChangeSet InventoryItemChangeSet) ProductionUnit {
	return MakeProductionUnit(ticksPerCycle, name, recipeChangeSet, InventoryItemChangeSet{})
}

func MakeNewProductionUnitWithNoRecipe(ticksPerCycle int, name string, changeSetForBuildingNew InventoryItemChangeSet) ProductionUnit {
	return MakeProductionUnit(ticksPerCycle, name, InventoryItemChangeSet{}, changeSetForBuildingNew)
}

func MakeNewProductionUnitWithNoChangeSet(ticksPerCycle int, name string) ProductionUnit {
	return MakeProductionUnit(ticksPerCycle, name, InventoryItemChangeSet{}, InventoryItemChangeSet{})
}

func TestThatTestThatWeCanGetInventoryItemSucceedsWhenInventoryHasItemsNeeded(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{}
	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("iron_ore", -1))
	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("copper_ore", -2))

	productionUnit := MakeNewProductionUnitWithNoRecipe(0, "Iron Mine", buildNewChangeSet)

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{3, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{4, "Copper Ore", "copper_ore"})

	newProductionUnit, newInventory := BuilNewProductionUnit(productionUnit, inventory)
	assert.Equal(t, 1, newProductionUnit.count)

	assert.Equal(t, 2, newInventory.items["iron_ore"].count)
	assert.Equal(t, 2, newInventory.items["copper_ore"].count)

	assert.Equal(t, 3, inventory.items["iron_ore"].count)
	assert.Equal(t, 4, inventory.items["copper_ore"].count)
}

func TestThatMultiplyChageSetForProductionMultipliesCorrectly(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))

	newChangeSet := MultiplyChageSetForProduction(changeSet, 2)

	assert.Equal(t, "iron_ore", newChangeSet[0].invetoryItemId)
	assert.Equal(t, 2, newChangeSet[0].changeAmount)
	assert.Equal(t, "copper_ore", newChangeSet[1].invetoryItemId)
	assert.Equal(t, 4, newChangeSet[1].changeAmount)
}

func TestThatCreateNewProductionBatchCreatesCorrectBatch(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.count = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{0, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Copper Ore", "copper_ore"})

	_, newInventory := CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, "Iron Ore", newInventory.items["iron_ore"].name)
	assert.Equal(t, 1, newInventory.items["iron_ore"].count)
	assert.Equal(t, "Copper Ore", newInventory.items["copper_ore"].name)
	assert.Equal(t, 3, newInventory.items["copper_ore"].count)
}

func TestThatCreateNewProductionBatchReturnsUpdatedProducionUnit(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(2, "Iron Mine", changeSet)
	productionUnit.count = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{0, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Copper Ore", "copper_ore"})

	newProductionUnit, newInventory := CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, 1, int(newProductionUnit.ticks_remaining))
	assert.Equal(t, "Iron Ore", newInventory.items["iron_ore"].name)
	assert.Equal(t, 0, newInventory.items["iron_ore"].count)
	assert.Equal(t, "Copper Ore", newInventory.items["copper_ore"].name)
	assert.Equal(t, 1, newInventory.items["copper_ore"].count)

	newProductionUnit, newInventory = CreateNewBatchIfTimeBecomes0(newProductionUnit, inventory)

	assert.Equal(t, 2, int(newProductionUnit.ticks_remaining))
	assert.Equal(t, "Iron Ore", newInventory.items["iron_ore"].name)
	assert.Equal(t, 1, newInventory.items["iron_ore"].count)
	assert.Equal(t, "Copper Ore", newInventory.items["copper_ore"].name)
	assert.Equal(t, 3, newInventory.items["copper_ore"].count)
}

func TestThatCreateNewProductionBatchDoesntChangeBatchIfCountIs0(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.count = 0

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{0, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Copper Ore", "copper_ore"})

	_, newInventory := CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, "Iron Ore", newInventory.items["iron_ore"].name)
	assert.Equal(t, 0, newInventory.items["iron_ore"].count)
	assert.Equal(t, "Copper Ore", newInventory.items["copper_ore"].name)
	assert.Equal(t, 1, newInventory.items["copper_ore"].count)
}

func TestThatCreateNewProductionBatchDoesntChangeBatchIfNotTimedOut(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(2, "Iron Mine", changeSet)
	productionUnit.count = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{0, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Copper Ore", "copper_ore"})

	_, newInventory := CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, "Iron Ore", newInventory.items["iron_ore"].name)
	assert.Equal(t, 0, newInventory.items["iron_ore"].count)
	assert.Equal(t, "Copper Ore", newInventory.items["copper_ore"].name)
	assert.Equal(t, 1, newInventory.items["copper_ore"].count)
}

func TestThatCreateNewProductionBatchCanCreateManyNewItems(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.count = 1000

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{0, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Copper Ore", "copper_ore"})

	_, newInventory := CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, "Iron Ore", newInventory.items["iron_ore"].name)
	assert.Equal(t, 1000, newInventory.items["iron_ore"].count)
	assert.Equal(t, "Copper Ore", newInventory.items["copper_ore"].name)
	assert.Equal(t, 2001, newInventory.items["copper_ore"].count)
}

func TestThatMultiplyChageSetForProductionDoesntChangeOriginal(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))

	MultiplyChageSetForProduction(changeSet, 2)

	assert.Equal(t, "iron_ore", changeSet[0].invetoryItemId)
	assert.Equal(t, 1, changeSet[0].changeAmount)
	assert.Equal(t, "copper_ore", changeSet[1].invetoryItemId)
	assert.Equal(t, 2, changeSet[1].changeAmount)
}

func TestThatMultiplyChageSetForProductionReturnsSameIf1(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))

	newChangeSet := MultiplyChageSetForProduction(changeSet, 1)

	assert.Equal(t, changeSet, newChangeSet)
	assert.Equal(t, "iron_ore", changeSet[0].invetoryItemId)
	assert.Equal(t, 1, changeSet[0].changeAmount)
	assert.Equal(t, "copper_ore", changeSet[1].invetoryItemId)
	assert.Equal(t, 2, changeSet[1].changeAmount)
}

func TestThatTestThatWeCanGetInventoryItemFailsWhenInventoryDoesntHaveItemsNeeded(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{}
	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("iron_ore", -1))
	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("copper_ore", -2))

	productionUnit := MakeNewProductionUnitWithNoRecipe(0, "Iron Mine", buildNewChangeSet)

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{0, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Copper Ore", "copper_ore"})

	newProductionUnit, newInventory := BuilNewProductionUnit(productionUnit, inventory)
	assert.Equal(t, 0, newProductionUnit.count)

	assert.Equal(t, 0, newInventory.items["iron_ore"].count)
	assert.Equal(t, 1, newInventory.items["copper_ore"].count)

	assert.Equal(t, 0, inventory.items["iron_ore"].count)
	assert.Equal(t, 1, inventory.items["copper_ore"].count)
}

func TestThatWeCanCheckThatWeCanBuilProductionUnit(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{}
	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("iron_ore", -1))

	productionUnit := MakeNewProductionUnitWithNoRecipe(0, "Iron Mine", buildNewChangeSet)

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Iron Ore", "iron_ore"})

	assert.True(t, CanBuilNewProductionUnit(productionUnit, inventory))
}

func TestThatWeCanCheckThatWeCantBuilProductionUnit(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{NewInventoryChange("iron_ore", -1)}

	productionUnit := MakeNewProductionUnitWithNoRecipe(0, "Iron Mine", buildNewChangeSet)

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{0, "Iron Ore", "iron_ore"})

	assert.False(t, CanBuilNewProductionUnit(productionUnit, inventory))
}

func TestThatWeCanBuildNewProductionUnitDoesntChangeInv(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{}
	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("iron_ore", -1))

	productionUnit := MakeNewProductionUnitWithNoRecipe(0, "Iron Mine", buildNewChangeSet)

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Iron Ore", "iron_ore"})

	assert.True(t, CanBuilNewProductionUnit(productionUnit, inventory))
	assert.Equal(t, 1, inventory.items["iron_ore"].count)

	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("copper_ore", -1))
	productionUnit.changeSetForBuildingNew = buildNewChangeSet

	assert.False(t, CanBuilNewProductionUnit(productionUnit, inventory))
	assert.Equal(t, 1, inventory.items["iron_ore"].count)
}

func TestThatMakeProductionUnitCreatesProductionUnitCorrectly(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{NewInventoryChange("iron_ore", -1)}
	recipeChangeSet := InventoryItemChangeSet{NewInventoryChange("copper_ore", 1)}
	produciotUnit := MakeProductionUnit(1, "Iron Mine", recipeChangeSet, buildNewChangeSet)

	assert.Equal(t, 0, produciotUnit.count)
	assert.Equal(t, 1, produciotUnit.ticks_per_cycle)
	assert.Equal(t, buildNewChangeSet, produciotUnit.changeSetForBuildingNew)
	assert.Equal(t, recipeChangeSet, produciotUnit.recipeChangeSet)
	assert.Equal(t, 1, int(produciotUnit.ticks_remaining))
	assert.Equal(t, "Iron Mine", produciotUnit.name)
}

func TestThatResetableIntResetsIfValueIsEqualToResetValue(t *testing.T) {
	resetAble := ResetableInt(2)

	assert.Equal(t, 2, int(resetAble))

	wasReset, newValue := resetAble.ResetIfValue(2, 3)
	assert.Equal(t, 3, int(newValue))
	assert.True(t, wasReset)
}

func TestThatResetableIntDoesntResetsIfValueIsNotEqualToResetValue(t *testing.T) {
	resetAble := ResetableInt(2)

	assert.Equal(t, 2, int(resetAble))

	wasReset, newValue := resetAble.ResetIfValue(4, 3)
	assert.Equal(t, 2, int(newValue))
	assert.False(t, wasReset)
}

func CheckThatUpdateProductionUnitTimerUpdateTickReturnsTrueAndResetsIfTicksRemainingIs0(t *testing.T) {
	productionUnit := MakeNewProductionUnitWithNoChangeSet(1, "Iron Mine")
	wasReset := false
	wasReset, productionUnit = UpdateProductionUnitTimer(productionUnit)

	assert.True(t, wasReset)
	assert.Equal(t, 0, productionUnit.ticks_remaining)
	assert.Equal(t, 2, productionUnit.ticks_per_cycle)

}

func TestThatMaybeResetTickReturnsTrueAndResetsIfTicksRemainingIs0(t *testing.T) {
	productionUnit := MakeNewProductionUnitWithNoChangeSet(1, "Iron Mine")
	wasReset := false
	wasReset, productionUnit = UpdateProductionUnitTimer(productionUnit)

	assert.True(t, wasReset)
	assert.Equal(t, 1, int(productionUnit.ticks_remaining))
	assert.Equal(t, 1, productionUnit.ticks_per_cycle)
}

func TestThatMaybeResetTickReturnsDecrementsCount(t *testing.T) {
	wasReset := false
	productionUnit := MakeNewProductionUnitWithNoChangeSet(2, "Iron Mine")
	wasReset, productionUnit = UpdateProductionUnitTimer(productionUnit)

	assert.Equal(t, 2, productionUnit.ticks_per_cycle)
	assert.Equal(t, 2, productionUnit.ticks_per_cycle)

	wasReset, productionUnit = UpdateProductionUnitTimer(productionUnit)
	assert.True(t, wasReset)
	assert.Equal(t, 2, productionUnit.ticks_per_cycle)
	assert.Equal(t, 2, int(productionUnit.ticks_remaining))
}

func TestThatResetableIntResetsDoesntChangeOriginal(t *testing.T) {
	resetAble := ResetableInt(2)

	assert.Equal(t, 2, int(resetAble))

	resetAble.ResetIfValue(2, 3)
	assert.Equal(t, 2, int(resetAble))
}

func TestThatPRoductionUnitHasStringFunc(t *testing.T) {
	produciotUnit := MakeNewProductionUnitWithNoChangeSet(0, "Iron Mine")
	assert.Equal(t, "Iron Mine", produciotUnit.String())
}

func TestThatPRoductionUnitHasCountFunc(t *testing.T) {
	produciotUnit := MakeNewProductionUnitWithNoChangeSet(1, "Iron Mine")
	produciotUnit.count = 1
	assert.Equal(t, 1, produciotUnit.Count())
}
