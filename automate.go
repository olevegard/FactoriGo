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

	production.iron_mines = MakleProductionUnit(1,
		func(inventory Inventory, count int) Inventory {
			inventory.iron_ore += count
			return inventory
		})

	production.copper_mines = MakleProductionUnit(1,
		func(inventory Inventory, count int) Inventory {
			inventory.copper_ore += count
			return inventory
		})

	production.iron_smelters = MakleProductionUnit(2,
		func(inventory Inventory, count int) Inventory {
			newItems := min(inventory.iron_ore, count)
			inventory.iron_ore -= newItems
			inventory.iron_plates += newItems
			return inventory
		})

	production.copper_smelters = MakleProductionUnit(2,
		func(inventory Inventory, count int) Inventory {
			newItems := min(inventory.copper_ore, count)
			inventory.copper_ore -= newItems
			inventory.copper_plates += newItems
			return inventory
		})

	return GameState{Inventory{}, production}
}
