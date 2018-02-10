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
			return ApplyInventoryItemChange(inventory, NewInventoryChange("iron_ore", count))
		}, "Iron mines")

	production.copper_mines = MakeProductionUnit(1,
		func(inventory Inventory, count int) Inventory {

			return ApplyInventoryItemChange(inventory, NewInventoryChange("copper_ore", count))
		}, "Copper mines")

	production.iron_smelters = MakeProductionUnit(2,
		func(inventory Inventory, count int) Inventory {
			newItems := min(inventory.items["iron_ore"].count, count)
			inventory = ApplyInventoryItemChange(inventory, NewInventoryChange("iron_ore", newItems*-1))

			inventory = ApplyInventoryItemChange(inventory, NewInventoryChange("iron_plates", newItems))
			return inventory
		}, "Iron furnaces")

	production.copper_smelters = MakeProductionUnit(2,
		func(inventory Inventory, count int) Inventory {
			newItems := min(inventory.items["copper_ore"].count, count)
			inventory = ApplyInventoryItemChange(inventory, NewInventoryChange("copper_ore", newItems*-1))

			inventory = ApplyInventoryItemChange(inventory, NewInventoryChange("copper_plates", newItems))
			return inventory
		}, "Copper furnaces")

	return GameState{MakeDefaultInventory(), production}
}

func MakeDefaultInventory() Inventory {
	inventory := NewInventory()
	io := InventoryItem{0, "Iron ore", "iron_ore"}
	co := InventoryItem{0, "Copper ore", "copper_ore"}

	ip := InventoryItem{0, "Iron plates", "iron_plates"}
	cp := InventoryItem{0, "Copper plates", "copper_plates"}
	inventory = AddInventoryItem(inventory, io)
	inventory = AddInventoryItem(inventory, co)
	inventory = AddInventoryItem(inventory, ip)
	inventory = AddInventoryItem(inventory, cp)
	return inventory

}
