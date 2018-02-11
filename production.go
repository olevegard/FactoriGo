package main

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

	recipeChangeSet         InventoryItemChangeSet // Use InventoryItemChangeSet here? We then only need one function to do recipe that'll work for all ProductionUnits
	changeSetForBuildingNew InventoryItemChangeSet // Use InventoryItemChangeSet here? We then only need one function to build new ProductionUnit that'll work for all ProductionUnits
	name                    string
}

func MakeProductionUnit(ticksPerCycle int, name string, recipeChangeSet, changeSetForBuildingNew InventoryItemChangeSet) ProductionUnit {
	unit := ProductionUnit{}
	unit.count = 0
	unit.ticks_per_cycle = ticksPerCycle
	unit.ticks_remaining = ResetableInt(ticksPerCycle)
	unit.recipeChangeSet = recipeChangeSet

	unit.changeSetForBuildingNew = changeSetForBuildingNew

	// unit.doRecipe = recipe
	unit.name = name

	return unit
}

func CanBuilNewProductionUnit(productionUnit ProductionUnit, inventory Inventory) bool {
	canBuild, _ := ApplyInventoryItemChangeSet(inventory, productionUnit.changeSetForBuildingNew)
	return canBuild
}

func BuilNewProductionUnit(productionUnit ProductionUnit, inventory Inventory) (ProductionUnit, Inventory) {
	canBuild, newInventory := ApplyInventoryItemChangeSet(inventory, productionUnit.changeSetForBuildingNew)

	if canBuild {
		productionUnit.count++
		return productionUnit, newInventory
	}
	return productionUnit, inventory
}

func MultiplyChageSetForProduction(changeSet InventoryItemChangeSet, factor int) InventoryItemChangeSet {
	if factor == 1 {
		return changeSet
	}
	newChangeSet := make(InventoryItemChangeSet, len(changeSet))

	for index, _ := range changeSet {
		// Since we allocate a new InventoryItemChangeSet, we need to set the entire item first
		newChangeSet[index] = changeSet[index]
		newChangeSet[index].changeAmount *= factor
	}

	return newChangeSet
}

func CreateNewBatchIfTimeBecomes0(productionUnit ProductionUnit, inventory Inventory) (ProductionUnit, Inventory) {
	shouldBuildNew, newProductionUnit := UpdateProductionUnitTimer(productionUnit)

	if !shouldBuildNew {
		return newProductionUnit, inventory
	}

	productionChangeSet := MultiplyChageSetForProduction(productionUnit.recipeChangeSet, productionUnit.count)
	_, newInventory := ApplyInventoryItemChangeSet(inventory, productionChangeSet)
	return newProductionUnit, newInventory
}

func (production ProductionUnit) String() string {
	return production.name
}

func (production ProductionUnit) Count() int {
	return production.count
}

type ResetableInt int

func (value ResetableInt) ResetIfValue(conditionValue, resetValue int) (bool, ResetableInt) {
	if int(value) == conditionValue {
		return true, ResetableInt(resetValue)
	}
	return false, value
}

func UpdateProductionUnitTimer(unit ProductionUnit) (bool, ProductionUnit) {
	unit.ticks_remaining--
	wasReset := false
	wasReset, unit.ticks_remaining = unit.ticks_remaining.ResetIfValue(0, unit.ticks_per_cycle)

	return wasReset, unit
}
