package main

import (
	"testing"
)

func TestThatMinReturnsSmalles(t *testing.T) {
	assertEqual_int(0, min(0, 1), t, "min")
	assertEqual_int(0, min(1, 0), t, "min")
	assertEqual_int(1, min(1, 1), t, "min")
}

func MakeProductionForTest(iron_mines, copper_mines, iron_smelters, copper_smelters int) Production {
	return Production{iron_mines, copper_mines, iron_smelters, copper_smelters}
}

// Test Smelted
// =============================================================================
func createProductionForTestingSmelted(iron_smelters, copper_smelters int) Production {
	return MakeProductionForTest(0, 0, iron_smelters, copper_smelters)
}

func createInventoryForTestingSmelted(iron_ore, copper_ore int) Inventory {
	inv := Inventory{}
	inv.iron_ore = iron_ore
	inv.copper_ore = copper_ore
	return inv
}

func testCreateInventoryForTestingSmelted(t *testing.T) {
	assertEqual_Inventory(Inventory{1, 2, 0, 0}, createInventoryForTestingSmelted(1, 2), t)
}

func testCreateProductionForTestingSmelted(t *testing.T) {
	assertEqual_Production(MakeProductionForTest(1, 2, 0, 0), createProductionForTestingSmelted(1, 2), t)
}

func TestSmelted_ThatUpdateDoesntChangeProd(t *testing.T) {
	var prod = MakeProduction()
	update_smelted(Inventory{}, prod)

	assertEqual_Production(MakeProduction(), prod, t)
}

func TestSmelted_ThatUpdateReturnsIdenticalInventoryIfNoProd(t *testing.T) {
	var inv = Inventory{}
	var new_inv = update_smelted(inv, MakeProduction())

	assertEqual_Inventory(new_inv, inv, t)
}

func TestSmelted_ThatUpdateReturnsIdenticalInventoryIfNoOre(t *testing.T) {
	var inv = createInventoryForTestingSmelted(0, 0)
	var new_inv = update_smelted(inv, createProductionForTestingSmelted(1, 1))

	assertEqual_Inventory(new_inv, inv, t)
}

func TestSmel_DoesnChangeOriginal(t *testing.T) {
	var inv = Inventory{}
	update_smelted(inv, MakeProductionForTest(1, 2, 3, 4))

	assertEqual_Inventory(Inventory{0, 0, 0, 0}, inv, t)
}

func TestSmelt_ProductionIncreasesInventoryIfBothSmeltersAndOre(t *testing.T) {
	var inv = createInventoryForTestingSmelted(1, 2)
	var new_inv = update_smelted(inv, createProductionForTestingSmelted(1, 2))

	assertEqual_Inventory(Inventory{0, 0, 1, 2}, new_inv, t)
}

func TestSmelt_SeveralUpdates(t *testing.T) {
	var inv = createInventoryForTestingSmelted(5, 10)
	var prod = createProductionForTestingSmelted(1, 2)

	var new_inv = update_smelted(inv, prod)
	assertEqual_Inventory(Inventory{4, 8, 1, 2}, new_inv, t)

	new_inv = update_smelted(new_inv, prod)
	assertEqual_Inventory(Inventory{3, 6, 2, 4}, new_inv, t)

	new_inv = update_smelted(new_inv, prod)
	assertEqual_Inventory(Inventory{2, 4, 3, 6}, new_inv, t)

	new_inv = update_smelted(new_inv, prod)
	assertEqual_Inventory(Inventory{1, 2, 4, 8}, new_inv, t)

	new_inv = update_smelted(new_inv, prod)
	assertEqual_Inventory(Inventory{0, 0, 5, 10}, new_inv, t)
}

func TestSmelt_ProdChange(t *testing.T) {
	var inv = createInventoryForTestingSmelted(10, 10)
	var prod = createProductionForTestingSmelted(2, 1)

	var new_inv = update_smelted(inv, prod)
	assertEqual_Inventory(Inventory{8, 9, 2, 1}, new_inv, t)

	prod = createProductionForTestingSmelted(3, 2)
	new_inv = update_smelted(new_inv, prod)
	assertEqual_Inventory(Inventory{5, 7, 5, 3}, new_inv, t)

	prod = createProductionForTestingSmelted(1, 0)
	new_inv = update_smelted(new_inv, prod)
	assertEqual_Inventory(Inventory{4, 7, 6, 3}, new_inv, t)

	prod = createProductionForTestingSmelted(4, 7)
	new_inv = update_smelted(new_inv, prod)
	assertEqual_Inventory(Inventory{0, 0, 10, 10}, new_inv, t)
}

func TestCantSmeltWhenNoInventory(t *testing.T) {
	var inv = createInventoryForTestingSmelted(0, 0)
	var prod = createProductionForTestingSmelted(2, 1)

	var new_inv = update_smelted(inv, prod)
	assertEqual_Inventory(Inventory{0, 0, 0, 0}, new_inv, t)
}

// Test Harvested
// =============================================================================
func createProductionForTestingHarvested(iron_mines, copper_mines int) Production {
	prod := MakeProduction()
	prod.iron_mines = iron_mines
	prod.copper_mines = copper_mines
	return prod
}

func createInventoryForTestingHarvested(iron_ore, copper_ore int) Inventory {
	inv := Inventory{}
	inv.iron_ore = iron_ore
	inv.copper_ore = copper_ore
	return inv
}

func TestHarvest_ThatUpdateDoesntChangeProd(t *testing.T) {
	var prod = MakeProduction()
	update_harvested(Inventory{}, prod)

	assertEqual_Production(MakeProduction(), prod, t)
}

func TestHarvest_ThatUpdateReturnsIdenticalInventoryIfNoProd(t *testing.T) {
	var inv = Inventory{}
	var new_inv = update_harvested(inv, MakeProduction())

	assertEqual_Inventory(new_inv, inv, t)
}

func TestHarvest_ThatUpdateDoesntChangeOriginalInventoryIfProd(t *testing.T) {
	var inv = Inventory{}
	update_harvested(inv, createProductionForTestingHarvested(1, 2))

	assertEqual_Inventory(Inventory{}, inv, t)
}

func TestHarvest_ThatUpdateReturnsUpdatedInventory(t *testing.T) {
	var inv = Inventory{}
	var new_inv = update_harvested(inv, createProductionForTestingHarvested(1, 2))

	assertEqual_Inventory(createInventoryForTestingHarvested(1, 2), new_inv, t)
}

func TestHarvest_ThatWeCanUpdateSeveralTimes(t *testing.T) {
	var inv = createInventoryForTestingHarvested(1, 2)
	var prod = createProductionForTestingHarvested(2, 1)

	var new_inv = update_harvested(inv, prod)
	assertEqual_Inventory(Inventory{3, 3, 0, 0}, new_inv, t)

	new_inv = update_harvested(new_inv, prod)
	assertEqual_Inventory(Inventory{5, 4, 0, 0}, new_inv, t)

	new_inv = update_harvested(new_inv, prod)
	assertEqual_Inventory(Inventory{7, 5, 0, 0}, new_inv, t)

	new_inv = update_harvested(new_inv, prod)
	assertEqual_Inventory(Inventory{9, 6, 0, 0}, new_inv, t)

	new_inv = update_harvested(new_inv, prod)
	assertEqual_Inventory(Inventory{11, 7, 0, 0}, new_inv, t)
}

func TestHarvest_ThatWeCanUpdateProduction(t *testing.T) {
	var inv = Inventory{}
	var prod = MakeProductionForTest(1, 2, 0, 0)

	var new_inv = update_harvested(inv, prod)
	assertEqual_Inventory(Inventory{1, 2, 0, 0}, new_inv, t)

	prod = MakeProductionForTest(2, 4, 0, 0)
	new_inv = update_harvested(new_inv, prod)
	assertEqual_Inventory(Inventory{3, 6, 0, 0}, new_inv, t)

	new_inv = update_harvested(new_inv, prod)
	assertEqual_Inventory(Inventory{5, 10, 0, 0}, new_inv, t)

	prod = MakeProductionForTest(1, 0, 5, 5)
	new_inv = update_harvested(new_inv, prod)
	assertEqual_Inventory(Inventory{6, 10, 0, 0}, new_inv, t)
}

// Harvested + Smelted
// =============================================================================
func TestHarvestSmelt_ThatWeCanUpdateProduction(t *testing.T) {
	var inv = Inventory{}
	var prod = MakeProductionForTest(1, 2, 3, 4)

	var new_inv = update_harvested(inv, prod)
	assertEqual_Inventory(Inventory{1, 2, 0, 0}, new_inv, t)

	new_inv = update_smelted(new_inv, prod)
	assertEqual_Inventory(Inventory{0, 0, 1, 2}, new_inv, t)
}

// Main update
// =============================================================================
func TestUpdate_UpdatesSmelt(t *testing.T) {
	var inv = createInventoryForTestingSmelted(1, 2)
	var prod = MakeProductionForTest(0, 0, 3, 4)

	var new_inv = update_inventory(inv, prod)
	assertEqual_Inventory(Inventory{0, 0, 1, 2}, new_inv, t)

	new_inv = update_inventory(inv, prod)
	assertEqual_Inventory(Inventory{0, 0, 1, 2}, new_inv, t)

	new_inv.iron_ore = 1
	new_inv.copper_ore = 2
	new_inv = update_inventory(new_inv, prod)
	assertEqual_Inventory(Inventory{0, 0, 2, 4}, new_inv, t)
}

func TestUpdate_UpdatesHarvest(t *testing.T) {
	var inv = createInventoryForTestingSmelted(0, 0)
	var prod = MakeProductionForTest(1, 2, 0, 0)

	var new_inv = update_inventory(inv, prod)
	assertEqual_Inventory(Inventory{1, 2, 0, 0}, new_inv, t)

	new_inv = update_inventory(new_inv, prod)
	assertEqual_Inventory(Inventory{2, 4, 0, 0}, new_inv, t)

	new_inv = update_inventory(new_inv, prod)
	assertEqual_Inventory(Inventory{3, 6, 0, 0}, new_inv, t)
}

func TestUpdate_UpdatesHarvestAndSmeltAllValuesSame(t *testing.T) {
	var inv = createInventoryForTestingSmelted(0, 0)
	var prod = MakeProductionForTest(1, 1, 1, 1)

	var new_inv = update_inventory(inv, prod)
	assertEqual_Inventory(Inventory{0, 0, 1, 1}, new_inv, t)

	new_inv = update_inventory(new_inv, prod)
	assertEqual_Inventory(Inventory{0, 0, 2, 2}, new_inv, t)

	new_inv = update_inventory(new_inv, prod)
	assertEqual_Inventory(Inventory{0, 0, 3, 3}, new_inv, t)
}

func TestUpdate_UpdatesHarvestAndSmeltSameValuesForMineAndSmelt(t *testing.T) {
	var inv = createInventoryForTestingSmelted(0, 0)
	var prod = MakeProductionForTest(2, 2, 1, 1)

	var new_inv = update_inventory(inv, prod)
	assertEqual_Inventory(Inventory{1, 1, 1, 1}, new_inv, t)

	new_inv = update_inventory(new_inv, prod)
	assertEqual_Inventory(Inventory{2, 2, 2, 2}, new_inv, t)

	new_inv = update_inventory(new_inv, prod)
	assertEqual_Inventory(Inventory{3, 3, 3, 3}, new_inv, t)
}

func TestUpdate_UpdatesHarvestAndSmelt(t *testing.T) {
	var inv = createInventoryForTestingSmelted(0, 0)
	var prod = MakeProductionForTest(3, 4, 1, 2)

	var new_inv = update_inventory(inv, prod)
	assertEqual_Inventory(Inventory{2, 2, 1, 2}, new_inv, t)

	new_inv = update_inventory(new_inv, prod)
	assertEqual_Inventory(Inventory{4, 4, 2, 4}, new_inv, t)

	new_inv = update_inventory(new_inv, prod)
	assertEqual_Inventory(Inventory{6, 6, 3, 6}, new_inv, t)
}

// Helper asserts
// =============================================================================
func assertEqual_Inventory(expected Inventory, actual Inventory, t *testing.T) {
	assertEqual_int(expected.iron_ore, actual.iron_ore, t, "Iron Ore")
	assertEqual_int(expected.copper_ore, actual.copper_ore, t, "Copper Ore")

	assertEqual_int(expected.iron_plates, actual.iron_plates, t, "Iron Plate")
	assertEqual_int(expected.copper_plates, actual.copper_plates, t, "Copper Plate")
}

func assertEqual_Production(expected Production, actual Production, t *testing.T) {
	assertEqual_int(expected.iron_mines, actual.iron_mines, t, "Iron Mines")
	assertEqual_int(expected.copper_mines, actual.copper_mines, t, "Copper Mines")

	assertEqual_int(expected.iron_smelters, actual.iron_smelters, t, "Iron Smelters")
	assertEqual_int(expected.copper_smelters, actual.copper_smelters, t, "Copper Smelters")
}
