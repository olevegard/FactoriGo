package main

func UpdateInventory(gameState GameState) GameState {
	gameState.production.iron_mines, gameState.inventory = CreateNewBatchIfTimeBecomes0(gameState.production.iron_mines, gameState.inventory)
	gameState.production.copper_mines, gameState.inventory = CreateNewBatchIfTimeBecomes0(gameState.production.copper_mines, gameState.inventory)

	gameState.production.iron_smelters, gameState.inventory = CreateNewBatchIfTimeBecomes0(gameState.production.iron_smelters, gameState.inventory)
	gameState.production.copper_smelters, gameState.inventory = CreateNewBatchIfTimeBecomes0(gameState.production.copper_smelters, gameState.inventory)
	return gameState
}

func MakeDefaultGameState() GameState {
	production := Production{}

	createNewChangeSet := InventoryItemChangeSet{}
	createNewChangeSet = append(createNewChangeSet, NewInventoryChange("iron_plates", -1))
	createNewChangeSet = append(createNewChangeSet, NewInventoryChange("copper_plates", -2))

	recipeChangeSet := InventoryItemChangeSet{NewInventoryChange("iron_ore", 1)}
	production.iron_mines = MakeProductionUnit(1, "Iron Mines", recipeChangeSet, createNewChangeSet)

	recipeChangeSet = InventoryItemChangeSet{NewInventoryChange("copper_ore", 1)}
	production.copper_mines = MakeProductionUnit(1, "Copper Mines", recipeChangeSet, createNewChangeSet)

	recipeChangeSet = InventoryItemChangeSet{NewInventoryChange("iron_plates", 1)}
	production.iron_smelters = MakeProductionUnit(1, "Iron Smelters", recipeChangeSet, createNewChangeSet)

	recipeChangeSet = InventoryItemChangeSet{NewInventoryChange("copper_plates", 1)}
	production.copper_smelters = MakeProductionUnit(1, "Copper Smelters", recipeChangeSet, createNewChangeSet)

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
