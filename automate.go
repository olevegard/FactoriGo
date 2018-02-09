package main

func DoRecipdeIfTimedOut(productionUnit ProductionUnit, inventory Inventory) (Inventory, ProductionUnit) {
	newProductionUnit, wasTimedOut := productionUnit.MaybeResetTick()

	if wasTimedOut {
		inventory = productionUnit.doRecipe(inventory, productionUnit.count)
	}

	return inventory, newProductionUnit
}

func UpdateInventory(gameState GameState) GameState {
	gameState.inventory, gameState.production.iron_mines =
		DoRecipdeIfTimedOut(gameState.production.iron_mines, gameState.inventory)

	gameState.inventory, gameState.production.copper_mines =
		DoRecipdeIfTimedOut(gameState.production.copper_mines, gameState.inventory)

	gameState.inventory, gameState.production.iron_smelters =
		DoRecipdeIfTimedOut(gameState.production.iron_smelters, gameState.inventory)

	gameState.inventory, gameState.production.copper_smelters =
		DoRecipdeIfTimedOut(gameState.production.copper_smelters, gameState.inventory)

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
