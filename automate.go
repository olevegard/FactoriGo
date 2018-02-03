package main

func update_harvested(inventory Inventory, production Production) Inventory {
	new_inventory := inventory

	new_inventory.iron_ore += production.iron_mines
	new_inventory.copper_ore += production.copper_mines

	return new_inventory
}

func update_smelted(inventory Inventory, production Production) Inventory {
	new_inventory := inventory

	smelted := min(new_inventory.copper_ore, production.copper_smelters)
	new_inventory.copper_ore -= smelted
	new_inventory.copper_plates += smelted

	smelted = min(new_inventory.iron_ore, production.iron_smelters)
	new_inventory.iron_ore -= smelted
	new_inventory.iron_plates += smelted

	return new_inventory
}

func update_inventory(inventory Inventory, production Production) Inventory {
	new_inventory := update_harvested(inventory, production)
	new_inventory = update_smelted(new_inventory, production)

	return new_inventory
}
