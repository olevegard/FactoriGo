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
	inventory.items = map[string]InventoryItem{}
	return inventory
}

func AddInventoryItem(inventory Inventory, inventoryItem InventoryItem) Inventory {
	newInventory := deepCopyInventory(inventory)
	newInventory.items[inventoryItem.id] = inventoryItem
	return newInventory
}

type InventoryItemChange struct {
	invetoryItemId string
	changeAmount   int
}

func NewInventoryChange(inventoryItemId string, changeAmount int) InventoryItemChange {
	return InventoryItemChange{inventoryItemId, changeAmount}
}

func GetNewCountAfterInventoryItemChange(inventory Inventory, inventoryItemChange InventoryItemChange) int {
	return inventory.items[inventoryItemChange.invetoryItemId].count + inventoryItemChange.changeAmount
}

func ApplyInventoryItemChange(inventory Inventory, inventoryItemChange InventoryItemChange) (bool, Inventory) {
	return ApplyInventoryItemChangeSet(inventory, InventoryItemChangeSet{inventoryItemChange})
}

func ApplyInventoryItemChangeSet(inventory Inventory, inventoryItemChangeSet InventoryItemChangeSet) (bool, Inventory) {
	// We need to deep copy here, because map is a value type, and applyInventoryItemChange will change the calues in the map.
	// This means that we can't revert to the originals if they fail unless we copy
	newInventory := deepCopyInventory(inventory)

	for _, change := range inventoryItemChangeSet {
		newCount := GetNewCountAfterInventoryItemChange(newInventory, change)

		if newCount < 0 {
			return false, inventory
		}
		newInventory = changeCountOfInventoryItemWithId(newInventory, change.invetoryItemId, newCount)
	}
	return true, newInventory
}

// Internal
// =============================================================================
// TODO: Use the above functions to manipulate inventory
// Make it possible to read from JSON
// Make buildings use this when subtracting inventory

func deepCopyInventory(originalInventory Inventory) Inventory {
	newInventory := originalInventory

	newInventory.items = make(map[string]InventoryItem)
	for key, value := range originalInventory.items {
		newInventory.items[key] = value
	}

	return newInventory
}

func changeCountAndReturnNew(original InventoryItem, newCount int) InventoryItem {
	original.count = newCount
	return original
}

// NOTE: This will change inventory
func changeCountOfInventoryItemWithId(inventory Inventory, inventoryItemId string, newCount int) Inventory {
	newInventory := inventory
	newInventory.items[inventoryItemId] = changeCountAndReturnNew(newInventory.items[inventoryItemId], newCount)
	return inventory
}

func (inventoryItem InventoryItem) String() string {
	return inventoryItem.name
}

func (inventoryItem InventoryItem) Count() int {
	return inventoryItem.count
}
