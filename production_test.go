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
	inventory = AddInventoryItem(inventory, NewInventoryItem(3, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(4, "Copper Ore", "copper_ore", false))

	newProductionUnit, newInventory := BuilNewProductionUnit(productionUnit, inventory)
	assert.Equal(t, 1, newProductionUnit.UnitCount)

	assert.Equal(t, 2, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 2, newInventory.Items["copper_ore"].ItemCount)

	assert.Equal(t, 3, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 4, inventory.Items["copper_ore"].ItemCount)
}

func TestThatMultiplyChageSetForProductionMultipliesCorrectly(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))

	newChangeSet := MultiplyChageSetForProduction(changeSet, 2)

	assert.Equal(t, "iron_ore", newChangeSet[0].InventoryItemId)
	assert.Equal(t, 2, newChangeSet[0].ChangeAmount)
	assert.Equal(t, "copper_ore", newChangeSet[1].InventoryItemId)
	assert.Equal(t, 4, newChangeSet[1].ChangeAmount)
}

func TestThatCreateNewProductionBatchCreatesCorrectBatch(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.UnitCount = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(0, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Copper Ore", "copper_ore", false))

	_, newInventory := CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, "Iron Ore", newInventory.Items["iron_ore"].Name)
	assert.Equal(t, 1, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, "Copper Ore", newInventory.Items["copper_ore"].Name)
	assert.Equal(t, 3, newInventory.Items["copper_ore"].ItemCount)
}

func TestThatCreateNewProductionBatchReturnsUpdatedProducionUnit(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(2, "Iron Mine", changeSet)
	productionUnit.UnitCount = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(0, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Copper Ore", "copper_ore", false))

	newProductionUnit, newInventory := CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, 1, int(newProductionUnit.TicksRemaining))
	assert.Equal(t, "Iron Ore", newInventory.Items["iron_ore"].Name)
	assert.Equal(t, 0, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, "Copper Ore", newInventory.Items["copper_ore"].Name)
	assert.Equal(t, 1, newInventory.Items["copper_ore"].ItemCount)

	newProductionUnit, newInventory = CreateNewBatchIfTimeBecomes0(newProductionUnit, inventory)

	assert.Equal(t, 2, int(newProductionUnit.TicksRemaining))
	assert.Equal(t, "Iron Ore", newInventory.Items["iron_ore"].Name)
	assert.Equal(t, 1, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, "Copper Ore", newInventory.Items["copper_ore"].Name)
	assert.Equal(t, 3, newInventory.Items["copper_ore"].ItemCount)
}

func TestThatCreateNewProductionBatchDoesntChangeBatchIfCountIs0(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.UnitCount = 0

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(0, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Copper Ore", "copper_ore", false))

	_, newInventory := CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, "Iron Ore", newInventory.Items["iron_ore"].Name)
	assert.Equal(t, 0, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, "Copper Ore", newInventory.Items["copper_ore"].Name)
	assert.Equal(t, 1, newInventory.Items["copper_ore"].ItemCount)
}

func TestThatCreateNewProductionBatchDoesntChangeBatchIfNotTimedOut(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(2, "Iron Mine", changeSet)
	productionUnit.UnitCount = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(0, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Copper Ore", "copper_ore", false))

	_, newInventory := CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, "Iron Ore", newInventory.Items["iron_ore"].Name)
	assert.Equal(t, 0, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, "Copper Ore", newInventory.Items["copper_ore"].Name)
	assert.Equal(t, 1, newInventory.Items["copper_ore"].ItemCount)
}

func TestThatCreateNewProductionBatchCanCreateManyNewItems(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.UnitCount = 1000

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(0, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Copper Ore", "copper_ore", false))

	_, newInventory := CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, "Iron Ore", newInventory.Items["iron_ore"].Name)
	assert.Equal(t, 1000, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, "Copper Ore", newInventory.Items["copper_ore"].Name)
	assert.Equal(t, 2001, newInventory.Items["copper_ore"].ItemCount)
}

func TestThatMultiplyChageSetForProductionDoesntChangeOriginal(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))

	MultiplyChageSetForProduction(changeSet, 2)

	assert.Equal(t, "iron_ore", changeSet[0].InventoryItemId)
	assert.Equal(t, 1, changeSet[0].ChangeAmount)
	assert.Equal(t, "copper_ore", changeSet[1].InventoryItemId)
	assert.Equal(t, 2, changeSet[1].ChangeAmount)
}

func TestThatMultiplyChageSetForProductionReturnsSameIf1(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	changeSet = append(changeSet, NewInventoryChange("copper_ore", 2))

	newChangeSet := MultiplyChageSetForProduction(changeSet, 1)

	assert.Equal(t, changeSet, newChangeSet)
	assert.Equal(t, "iron_ore", changeSet[0].InventoryItemId)
	assert.Equal(t, 1, changeSet[0].ChangeAmount)
	assert.Equal(t, "copper_ore", changeSet[1].InventoryItemId)
	assert.Equal(t, 2, changeSet[1].ChangeAmount)
}

func TestThatTestThatWeCanGetInventoryItemFailsWhenInventoryDoesntHaveItemsNeeded(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{}
	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("iron_ore", -1))
	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("copper_ore", -2))

	productionUnit := MakeNewProductionUnitWithNoRecipe(0, "Iron Mine", buildNewChangeSet)

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(0, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Copper Ore", "copper_ore", false))

	newProductionUnit, newInventory := BuilNewProductionUnit(productionUnit, inventory)
	assert.Equal(t, 0, newProductionUnit.UnitCount)

	assert.Equal(t, 0, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 1, newInventory.Items["copper_ore"].ItemCount)

	assert.Equal(t, 0, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 1, inventory.Items["copper_ore"].ItemCount)
}

func TestThatWeCanCheckThatWeCanBuilProductionUnit(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{}
	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("iron_ore", -1))

	productionUnit := MakeNewProductionUnitWithNoRecipe(0, "Iron Mine", buildNewChangeSet)

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron Ore", "iron_ore", false))

	assert.True(t, CanBuilNewProductionUnit(productionUnit, inventory))
}

func TestThatWeCanCheckThatWeCantBuilProductionUnit(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{NewInventoryChange("iron_ore", -1)}

	productionUnit := MakeNewProductionUnitWithNoRecipe(0, "Iron Mine", buildNewChangeSet)

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(0, "Iron Ore", "iron_ore", false))

	assert.False(t, CanBuilNewProductionUnit(productionUnit, inventory))
}

func TestThatCanBuildNewProductionUnitDoesntChangeInv(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{}
	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("iron_ore", -1))

	productionUnit := MakeNewProductionUnitWithNoRecipe(0, "Iron Mine", buildNewChangeSet)
	productionUnit.UnitCount = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron Ore", "iron_ore", false))

	assert.True(t, CanBuilNewProductionUnit(productionUnit, inventory))
	assert.Equal(t, 1, inventory.Items["iron_ore"].ItemCount)

	buildNewChangeSet = append(buildNewChangeSet, NewInventoryChange("copper_ore", -1))
	productionUnit.ChangeSetForBuildingNew = buildNewChangeSet

	assert.False(t, CanBuilNewProductionUnit(productionUnit, inventory))
	assert.Equal(t, 1, inventory.Items["iron_ore"].ItemCount)
}

func TestThatMakeProductionUnitCreatesProductionUnitCorrectly(t *testing.T) {
	buildNewChangeSet := InventoryItemChangeSet{NewInventoryChange("iron_ore", -1)}
	recipeChangeSet := InventoryItemChangeSet{NewInventoryChange("copper_ore", 1)}
	produciotUnit := MakeProductionUnit(1, "Iron Mine", recipeChangeSet, buildNewChangeSet)

	assert.Equal(t, 0, produciotUnit.UnitCount)
	assert.Equal(t, 1, produciotUnit.TicksPerCycle)
	assert.Equal(t, buildNewChangeSet, produciotUnit.ChangeSetForBuildingNew)
	assert.Equal(t, recipeChangeSet, produciotUnit.RecipeChangeSet)
	assert.Equal(t, 1, int(produciotUnit.TicksRemaining))
	assert.Equal(t, "Iron Mine", produciotUnit.Name)
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
	inventory := Inventory{}
	productionUnit := MakeNewProductionUnitWithNoChangeSet(1, "Iron Mine")

	wasReset := false
	wasReset, productionUnit = UpdateProductionUnitTimer(productionUnit, inventory)

	assert.True(t, wasReset)
	assert.Equal(t, 0, productionUnit.TicksRemaining)
	assert.Equal(t, 2, productionUnit.TicksPerCycle)
}

func TestThatMaybeResetTickReturnsTrueAndResetsIfTicksRemainingIs0(t *testing.T) {
	productionUnit := MakeNewProductionUnitWithNoChangeSet(1, "Iron Mine")
	productionUnit.UnitCount = 1
	wasReset := false
	wasReset, productionUnit = UpdateProductionUnitTimer(productionUnit, Inventory{})

	assert.True(t, wasReset)
	assert.Equal(t, 1, int(productionUnit.TicksRemaining))
	assert.Equal(t, 1, productionUnit.TicksPerCycle)
}

func TestThatMaybeResetTickReturnsDecrementedCount(t *testing.T) {
	wasReset := false
	productionUnit := MakeNewProductionUnitWithNoChangeSet(2, "Iron Mine")
	productionUnit.UnitCount = 1
	wasReset, productionUnit = UpdateProductionUnitTimer(productionUnit, Inventory{})

	assert.False(t, wasReset)
	assert.Equal(t, 2, productionUnit.TicksPerCycle)
	assert.Equal(t, 1, int(productionUnit.TicksRemaining))

	wasReset, productionUnit = UpdateProductionUnitTimer(productionUnit, Inventory{})
	assert.True(t, wasReset)
	assert.Equal(t, 2, productionUnit.TicksPerCycle)
	assert.Equal(t, 2, int(productionUnit.TicksRemaining))
}

func TestThatResetableIntResetsDoesntChangeOriginal(t *testing.T) {
	resetAble := ResetableInt(2)

	assert.Equal(t, 2, int(resetAble))

	resetAble.ResetIfValue(2, 3)
	assert.Equal(t, 2, int(resetAble))
}

func TestThatProductionUnitHasStringFunc(t *testing.T) {
	produciotUnit := MakeNewProductionUnitWithNoChangeSet(0, "Iron Mine")
	assert.Equal(t, "Iron Mine", produciotUnit.String())
}

func TestThatProductionUnitHasCountFunc(t *testing.T) {
	produciotUnit := MakeNewProductionUnitWithNoChangeSet(1, "Iron Mine")
	produciotUnit.UnitCount = 1
	assert.Equal(t, 1, produciotUnit.Count())
}

func TestThatWeCanCreatePartial(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_plate", 1))
	changeSet = append(changeSet, NewInventoryChange("iron_ore", -1))
	changeSet = append(changeSet, NewInventoryChange("coal", -2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.UnitCount = 2
	productionUnit.TicksRemaining = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(3, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(3, "Coal", "coal", false))

	productionUnit, inventory = CreateNewBatchIfTimeBecomes0(productionUnit, inventory)
	assert.Equal(t, 1, inventory.Items["iron_plate"].ItemCount)
	assert.Equal(t, 2, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 1, inventory.Items["coal"].ItemCount)
}

func TestThatWeCanGetMaxFactorIfNotAllCanBeMade(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_plate", 1))
	changeSet = append(changeSet, NewInventoryChange("iron_ore", -1))
	changeSet = append(changeSet, NewInventoryChange("coal", -2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.UnitCount = 2
	productionUnit.TicksRemaining = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(3, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(3, "Coal", "coal", false))

	factor := GetMaxFactorForProductionUnit(productionUnit, inventory)
	assert.Equal(t, 1, factor)
}

func TestThatWeCanGetMaxFactorIfAllCanBeMade(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_plate", 1))
	changeSet = append(changeSet, NewInventoryChange("iron_ore", -1))
	changeSet = append(changeSet, NewInventoryChange("coal", -2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.UnitCount = 2
	productionUnit.TicksRemaining = 0

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(2, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(4, "Coal", "coal", false))

	factor := GetMaxFactorForProductionUnit(productionUnit, inventory)
	assert.Equal(t, 2, factor)
}

func TestThatWeCanGetMaxFactorIfNothingCanBeMade(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_plate", 1))
	changeSet = append(changeSet, NewInventoryChange("iron_ore", -10))
	changeSet = append(changeSet, NewInventoryChange("coal", -2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.UnitCount = 2
	productionUnit.TicksRemaining = 0

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(3, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(2, "Coal", "coal", false))

	factor := GetMaxFactorForProductionUnit(productionUnit, inventory)
	productionUnit, inventory = CreateNewBatchIfTimeBecomes0(productionUnit, inventory)
	assert.Equal(t, 0, factor)
}

func TestThatWeCanGetMaxFactorCantExceedUnits(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_plate", 1))
	changeSet = append(changeSet, NewInventoryChange("iron_ore", -1))
	changeSet = append(changeSet, NewInventoryChange("coal", -2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.UnitCount = 2
	productionUnit.TicksRemaining = 0

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(3000, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(2000, "Coal", "coal", false))

	factor := GetMaxFactorForProductionUnit(productionUnit, inventory)

	productionUnit, inventory = CreateNewBatchIfTimeBecomes0(productionUnit, inventory)
	assert.Equal(t, 2, factor)
}

func TestThatMaxFactorIs0IfNoUnits(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_plate", 1))
	changeSet = append(changeSet, NewInventoryChange("iron_ore", -1))
	changeSet = append(changeSet, NewInventoryChange("coal", -2))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.UnitCount = 0
	productionUnit.TicksRemaining = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(3000, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(2000, "Coal", "coal", false))

	factor := GetMaxFactorForProductionUnit(productionUnit, inventory)

	productionUnit, inventory = CreateNewBatchIfTimeBecomes0(productionUnit, inventory)
	assert.Equal(t, 0, factor)
}

func TestThatMaxFactorIsCountOfUnitsIfReceipeDontNeedItems(t *testing.T) {
	changeSet := InventoryItemChangeSet{}
	changeSet = append(changeSet, NewInventoryChange("iron_ore", 1))
	productionUnit := MakeNewProductionUnitWithNoBuildNew(1, "Iron Mine", changeSet)
	productionUnit.UnitCount = 1000
	productionUnit.TicksRemaining = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(3000, "Iron Ore", "iron_ore", false))

	factor := GetMaxFactorForProductionUnit(productionUnit, inventory)

	productionUnit, inventory = CreateNewBatchIfTimeBecomes0(productionUnit, inventory)
	assert.Equal(t, 1000, factor)
}

func TestThatTickDoesntDecrementRemainingIfNoUnits(t *testing.T) {
	wasReset := false
	productionUnit := MakeNewProductionUnitWithNoChangeSet(2, "Iron Mine")
	productionUnit.UnitCount = 0
	wasReset, productionUnit = UpdateProductionUnitTimer(productionUnit, Inventory{})

	assert.False(t, wasReset)
	assert.Equal(t, 2, int(productionUnit.TicksRemaining))
	assert.Equal(t, 2, productionUnit.TicksPerCycle)

	wasReset, productionUnit = UpdateProductionUnitTimer(productionUnit, Inventory{})
	assert.False(t, wasReset)
	assert.Equal(t, 2, productionUnit.TicksPerCycle)
	assert.Equal(t, 2, int(productionUnit.TicksRemaining))
}

func TestThatCanCreateNewBatchDoesnChangeInv(t *testing.T) {
	recipeNewChangeSet := InventoryItemChangeSet{}
	recipeNewChangeSet = append(recipeNewChangeSet, NewInventoryChange("iron_ore", -2))

	productionUnit := MakeNewProductionUnitWithNoBuildNew(0, "Iron Mine", recipeNewChangeSet)
	productionUnit.UnitCount = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron Ore", "iron_ore", false))

	assert.False(t, CanCreateNewBatch(productionUnit, inventory))
	assert.Equal(t, 1, inventory.Items["iron_ore"].ItemCount)
}

func TestThatCanCreateNewBatchReturnsTrueIfEnoughInv(t *testing.T) {
	recipeNewChangeSet := InventoryItemChangeSet{}
	recipeNewChangeSet = append(recipeNewChangeSet, NewInventoryChange("iron_ore", -1))

	productionUnit := MakeNewProductionUnitWithNoBuildNew(0, "Iron Mine", recipeNewChangeSet)
	productionUnit.UnitCount = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron Ore", "iron_ore", false))

	assert.True(t, CanCreateNewBatch(productionUnit, inventory))
	assert.Equal(t, 1, inventory.Items["iron_ore"].ItemCount)
}

func TestThatCanCreateNewBatchReturnsFalseIfNotEnoughInv(t *testing.T) {
	recipeNewChangeSet := InventoryItemChangeSet{}
	recipeNewChangeSet = append(recipeNewChangeSet, NewInventoryChange("iron_ore", -2))

	productionUnit := MakeNewProductionUnitWithNoBuildNew(0, "Iron Mine", recipeNewChangeSet)
	productionUnit.UnitCount = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron Ore", "iron_ore", false))

	assert.False(t, CanCreateNewBatch(productionUnit, inventory))
	assert.Equal(t, 1, inventory.Items["iron_ore"].ItemCount)
}

func TestThatCanCreateNewBatchReturnsTrueIfPatial(t *testing.T) {
	recipeNewChangeSet := InventoryItemChangeSet{}
	recipeNewChangeSet = append(recipeNewChangeSet, NewInventoryChange("iron_ore", -1))

	productionUnit := MakeNewProductionUnitWithNoBuildNew(0, "Iron Mine", recipeNewChangeSet)
	productionUnit.UnitCount = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron Ore", "iron_ore", false))

	assert.True(t, CanCreateNewBatch(productionUnit, inventory))
	assert.Equal(t, 1, inventory.Items["iron_ore"].ItemCount)
}

func TestThatUpdateTimerDoesnUpdateIfCantBuildNew(t *testing.T) {
	recipeNewChangeSet := InventoryItemChangeSet{}
	recipeNewChangeSet = append(recipeNewChangeSet, NewInventoryChange("iron_ore", -2))

	productionUnit := MakeNewProductionUnitWithNoBuildNew(2, "Iron Mine", recipeNewChangeSet)
	productionUnit.UnitCount = 1

	wasUpdated := false
	wasUpdated, productionUnit = UpdateProductionUnitTimer(productionUnit, Inventory{})

	assert.False(t, wasUpdated)
	assert.Equal(t, 2, int(productionUnit.TicksRemaining))
}

func TestThatUpdateTimerUpdatesIfCantBuildPartial(t *testing.T) {
	recipeNewChangeSet := InventoryItemChangeSet{}
	recipeNewChangeSet = append(recipeNewChangeSet, NewInventoryChange("iron_ore", -1))

	productionUnit := MakeNewProductionUnitWithNoBuildNew(2, "Iron Mine", recipeNewChangeSet)
	productionUnit.UnitCount = 1

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron Ore", "iron_ore", false))

	assert.Equal(t, 1, inventory.Items["iron_ore"].Count())

	wasUpdated := false
	wasUpdated, productionUnit = UpdateProductionUnitTimer(productionUnit, inventory)

	assert.False(t, wasUpdated)

	assert.Equal(t, 1, int(productionUnit.TicksRemaining))
}

func TestThatTickDoesntDecrementRemainingIfNotProducing(t *testing.T) {
	recipeNewChangeSet := InventoryItemChangeSet{}
	recipeNewChangeSet = append(recipeNewChangeSet, NewInventoryChange("iron_ore", -2))

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron Ore", "iron_ore", false))

	productionUnit := MakeNewProductionUnitWithNoBuildNew(2, "Iron Mine", recipeNewChangeSet)
	productionUnit.UnitCount = 1

	assert.False(t, CanCreateNewBatch(productionUnit, inventory))
	assert.Equal(t, 1, inventory.Items["iron_ore"].ItemCount)

	productionUnit, inventory = CreateNewBatchIfTimeBecomes0(productionUnit, inventory)

	assert.Equal(t, 1, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 2, int(productionUnit.TicksRemaining))
	assert.Equal(t, 2, productionUnit.TicksPerCycle)
}
