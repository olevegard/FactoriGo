package main

func UpdateTimeRemaining(productionUnit ProductionUnit) ProductionUnit {
	if productionUnit.ticks_remaining == 0 {
		productionUnit.ticks_remaining = productionUnit.ticks_per_cycle
	}

	return productionUnit
}

func GetProductionIfTimeout(productionUnit ProductionUnit, inventory Inventory) Inventory {
	if productionUnit.ticks_remaining == 0 {
		inventory = productionUnit.doRecipe(inventory, productionUnit.count)
	}

	return inventory
}

func CheckProductionUnitForNewBatch(productionUnit ProductionUnit, inventory Inventory) (Inventory, ProductionUnit) {
	productionUnit.ticks_remaining--

	return GetProductionIfTimeout(productionUnit, inventory), UpdateTimeRemaining(productionUnit)
}

func UpdateInventory(gameState GameState) GameState {
	gameState.inventory, gameState.production.iron_mines =
		CheckProductionUnitForNewBatch(gameState.production.iron_mines, gameState.inventory)

	gameState.inventory, gameState.production.copper_mines =
		CheckProductionUnitForNewBatch(gameState.production.copper_mines, gameState.inventory)

	gameState.inventory, gameState.production.iron_smelters =
		CheckProductionUnitForNewBatch(gameState.production.iron_smelters, gameState.inventory)

	gameState.inventory, gameState.production.copper_smelters =
		CheckProductionUnitForNewBatch(gameState.production.copper_smelters, gameState.inventory)

	return gameState
}

func MakeDefaultGameState() GameState {
	production := Production{}

	production.iron_mines = MakeProductionUnit(1,
		func(inventory Inventory, count int) Inventory {
			inventory.iron_ore.count += count
			return inventory
		}, "Iron mines")

	production.copper_mines = MakeProductionUnit(1,
		func(inventory Inventory, count int) Inventory {
			inventory.copper_ore.count += count
			return inventory
		}, "Copper mines")

	production.iron_smelters = MakeProductionUnit(2,
		func(inventory Inventory, count int) Inventory {
			newItems := min(inventory.iron_ore.count, count)
			inventory.iron_ore.count -= newItems
			inventory.iron_plates.count += newItems
			return inventory
		}, "Iron furnaces")

	production.copper_smelters = MakeProductionUnit(2,
		func(inventory Inventory, count int) Inventory {
			newItems := min(inventory.copper_ore.count, count)
			inventory.copper_ore.count -= newItems
			inventory.copper_plates.count += newItems
			return inventory
		}, "Copper furnaces")

	return GameState{MakeDefaultInventory(), production}
}

func MakeDefaultInventory() Inventory {
	io := InventoryItem{0, "Iron ore"}
	co := InventoryItem{0, "Copper ore"}

	ip := InventoryItem{0, "Iron plastes"}
	cp := InventoryItem{0, "Copper plates"}
	return Inventory{io, co, ip, cp}

}
