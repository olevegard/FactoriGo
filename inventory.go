package main

type InventoryItemChangeSet []InventoryItemChange

type Inventory struct {
	items map[string]InventoryItem
}

type InventoryItem struct {
	count int
	name  string
	id    string
}

func NewInventory() Inventory {
	inventory := Inventory{}
	inventory.items = make(map[string]InventoryItem)
	return inventory
}

func AddInventoryItem(inventory Inventory, inventoryItem InventoryItem) Inventory {
	inventory.items[inventoryItem.id] = inventoryItem
	return inventory
}

func ChangeCountAndReturnNew(original InventoryItem, newCount int) InventoryItem {
	original.count = newCount
	return original
}

func ChangeCountOfInventoryItemWithId(inventory Inventory, inventoryItemId string, newCount int) Inventory {
	inventory.items[inventoryItemId] = ChangeCountAndReturnNew(inventory.items[inventoryItemId], newCount)
	return inventory
}

type InventoryItemChange struct {
	invetoryItemId string
	changeAmount   int
}

func NewInventoryChange(inventoryItemId string, changeAmount int) InventoryItemChange {
	return InventoryItemChange{inventoryItemId, changeAmount}
}

func ApplyInventoryItemChange(inventory Inventory, inventoryItemChange InventoryItemChange) Inventory {
	newCount := inventory.items[inventoryItemChange.invetoryItemId].count + inventoryItemChange.changeAmount

	return ChangeCountOfInventoryItemWithId(inventory, inventoryItemChange.invetoryItemId, newCount)
}

func ApplyInventoryItemChangeSet(inventory Inventory, inventoryItemChangeSet InventoryItemChangeSet) Inventory {
	for _, change := range inventoryItemChangeSet {
		inventory = ApplyInventoryItemChange(inventory, change)
	}

	return inventory
}
