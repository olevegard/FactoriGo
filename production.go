package main

type Production struct {
	// Harvested
	IronMines   ProductionUnit
	CopperMines ProductionUnit

	// Smelted
	IronSmelters   ProductionUnit
	CopperSmelters ProductionUnit
}

// Used to update Inventory for a recipe
type DoRecipe func(Inventory, int) Inventory

// Any thing that can produce/mine/smelt etc...
type ProductionUnit struct {
	UnitCount int

	TicksRemaining ResetableInt // Time left until next batch of units
	TicksPerCycle  int          // Time it takes to complete a full cycle

	RecipeChangeSet         InventoryItemChangeSet // Use InventoryItemChangeSet here? We then only need one function to do recipe that'll work for all ProductionUnits
	ChangeSetForBuildingNew InventoryItemChangeSet // Use InventoryItemChangeSet here? We then only need one function to build new ProductionUnit that'll work for all ProductionUnits
	Name                    string
}

func MakeProductionUnit(ticksPerCycle int, name string, recipeChangeSet, changeSetForBuildingNew InventoryItemChangeSet) ProductionUnit {
	unit := ProductionUnit{}
	unit.UnitCount = 0
	unit.TicksPerCycle = ticksPerCycle
	unit.TicksRemaining = ResetableInt(ticksPerCycle)
	unit.RecipeChangeSet = recipeChangeSet

	unit.ChangeSetForBuildingNew = changeSetForBuildingNew

	// unit.doRecipe = recipe
	unit.Name = name

	return unit
}

func CanBuilNewProductionUnit(productionUnit ProductionUnit, inventory Inventory) bool {
	canBuild, _ := ApplyInventoryItemChangeSet(inventory, productionUnit.ChangeSetForBuildingNew)
	return canBuild
}

func BuilNewProductionUnit(productionUnit ProductionUnit, inventory Inventory) (ProductionUnit, Inventory) {
	canBuild, newInventory := ApplyInventoryItemChangeSet(inventory, productionUnit.ChangeSetForBuildingNew)

	if canBuild {
		productionUnit.UnitCount++
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
		newChangeSet[index].ChangeAmount *= factor
	}

	return newChangeSet
}

func CreateNewBatchIfTimeBecomes0(productionUnit ProductionUnit, inventory Inventory) (ProductionUnit, Inventory) {
	shouldBuildNew, newProductionUnit := UpdateProductionUnitTimer(productionUnit)

	if !shouldBuildNew {
		return newProductionUnit, inventory
	}

	productionChangeSet := MultiplyChageSetForProduction(productionUnit.RecipeChangeSet, productionUnit.UnitCount)
	_, newInventory := ApplyInventoryItemChangeSet(inventory, productionChangeSet)
	return newProductionUnit, newInventory
}

func (production ProductionUnit) String() string {
	return production.Name
}

func (production ProductionUnit) Count() int {
	return production.UnitCount
}

type ResetableInt int

func (value ResetableInt) ResetIfValue(conditionValue, resetValue int) (bool, ResetableInt) {
	if int(value) == conditionValue {
		return true, ResetableInt(resetValue)
	}
	return false, value
}

func UpdateProductionUnitTimer(unit ProductionUnit) (bool, ProductionUnit) {
	unit.TicksRemaining--
	wasReset := false
	wasReset, unit.TicksRemaining = unit.TicksRemaining.ResetIfValue(0, unit.TicksPerCycle)

	return wasReset, unit
}
