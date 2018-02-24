package main

type Icon interface{}

type InventoryItemChangeSet []InventoryItemChange

type Inventory struct {
	Items map[string]InventoryItem `json:"items"`
}

type InventoryItem struct {
	ItemCount         int    `json:"count"`
	Name              string `json:"name"`
	Id                string `json:"id,id"`
	CanBeMadeManually bool   `json:"can_be_made_manually"`
	Icon              Icon   `json:"-"`
}

func NewInventoryItem(itemCount int, name string, id string, canBeMadeManually bool) InventoryItem {
	return InventoryItem{itemCount, name, id, canBeMadeManually, nil}
}

func NewInventory() Inventory {
	inventory := Inventory{}
	inventory.Items = map[string]InventoryItem{}
	return inventory
}

func AddInventoryItem(inventory Inventory, inventoryItem InventoryItem) Inventory {
	newInventory := deepCopyInventory(inventory)
	newInventory.Items[inventoryItem.Id] = inventoryItem
	return newInventory
}

type InventoryItemChange struct {
	InventoryItemId string
	ChangeAmount    int
}

func NewInventoryChange(inventoryItemId string, changeAmount int) InventoryItemChange {
	return InventoryItemChange{inventoryItemId, changeAmount}
}

func GetNewCountAfterInventoryItemChange(inventory Inventory, inventoryItemChange InventoryItemChange) int {
	return inventory.Items[inventoryItemChange.InventoryItemId].ItemCount + inventoryItemChange.ChangeAmount
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
		newInventory = changeCountOfInventoryItemWithId(newInventory, change.InventoryItemId, newCount)
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

	newInventory.Items = make(map[string]InventoryItem)
	for key, value := range originalInventory.Items {
		newInventory.Items[key] = value
	}

	return newInventory
}

func changeCountAndReturnNew(original InventoryItem, newCount int) InventoryItem {
	original.ItemCount = newCount
	return original
}

// NOTE: This will change inventory
func changeCountOfInventoryItemWithId(inventory Inventory, inventoryItemId string, newCount int) Inventory {
	newInventory := inventory
	newInventory.Items[inventoryItemId] = changeCountAndReturnNew(newInventory.Items[inventoryItemId], newCount)
	return inventory
}

func (inventoryItem InventoryItem) String() string {
	return inventoryItem.Name
}

func (inventoryItem InventoryItem) Count() int {
	return inventoryItem.ItemCount
}
