package main

import (
	"testing"

	"github.com/stvp/assert"
)

func TestThatWeCanGetInventoryItem(t *testing.T) {
	inventory := NewInventory()
	inventoryItem := InventoryItem{1, "Iron ore", "iron_ore"}
	inventory = AddInventoryItem(inventory, inventoryItem)

	returnedItem := inventory.items[inventoryItem.id]

	if returnedItem.id != "iron_ore" {
		t.Fatal("Couldn't get item from Inventory")
	}

	assert.Equal(t, returnedItem, inventoryItem)
}

func TestThatChangeCountAndReturnNewReturnsUpdateInventoryItem(t *testing.T) {
	inventory := NewInventory()
	inventoryItem := InventoryItem{1, "Iron ore", "iron_ore"}
	inventory = AddInventoryItem(inventory, inventoryItem)

	inventory.items[inventoryItem.id] = changeCountAndReturnNew(inventory.items["iron_ore"], 4)

	assert.Equal(t, 4, inventory.items["iron_ore"].count)
}

func TestThatChangeCountOfInventoryItemWithIdUpdatesInventory(t *testing.T) {
	inventory := NewInventory()
	inventoryItem := InventoryItem{1, "Iron ore", "iron_ore"}
	inventory = AddInventoryItem(inventory, inventoryItem)
	newInventory := changeCountOfInventoryItemWithId(inventory, "iron_ore", 2)

	assert.Equal(t, 2, newInventory.items["iron_ore"].count)
}

func TestThatChangeCountOfInventoryItemWithIdDoesntChangeOriginalInv(t *testing.T) {
	inventory := NewInventory()
	inventoryItem := InventoryItem{1, "Iron ore", "iron_ore"}
	inventory = AddInventoryItem(inventory, inventoryItem)
	newInventory := changeCountOfInventoryItemWithId(inventory, "iron_ore", 2)

	assert.Equal(t, 2, newInventory.items["iron_ore"].count)
}

func TestThatCreateInventoryChangeReturnsExpectedInventoryItemChange(t *testing.T) {
	expectedInventoryChange := InventoryItemChange{}
	expectedInventoryChange.changeAmount = 3
	expectedInventoryChange.invetoryItemId = "iron_ore"
	actuaInventoryChange := NewInventoryChange("iron_ore", 3)

	assert.Equal(t, expectedInventoryChange, actuaInventoryChange)
}

func TestThatApplyInventoryItemChangeDoesntChangeOtherValues(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", 3)
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{5, "Iron ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{10, "Copper ore", "copper_ore"})

	wasApplied, newInventory := ApplyInventoryItemChange(inventory, inventoryChange)

	assert.True(t, wasApplied)
	assert.Equal(t, 8, newInventory.items["iron_ore"].count)
	assert.Equal(t, 10, newInventory.items["copper_ore"].count)
	assert.Equal(t, 5, inventory.items["iron_ore"].count)
	assert.Equal(t, 10, inventory.items["copper_ore"].count)
}

func TestThatApplyInventoryItemChangeAppliesPositiveChange(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", 3)
	inventoryItem := InventoryItem{5, "Iron ore", "iron_ore"}
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, inventoryItem)

	_, inventory = ApplyInventoryItemChange(inventory, inventoryChange)

	assert.Equal(t, 8, inventory.items[inventoryItem.id].count)
}

func TestThatApplyInventoryItemChangeAppliesNegativeChange(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", -3)
	inventoryItem := InventoryItem{5, "Iron ore", "iron_ore"}
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, inventoryItem)

	_, inventory = ApplyInventoryItemChange(inventory, inventoryChange)

	assert.Equal(t, 2, inventory.items[inventoryItem.id].count)
}

func TestThatWeCanApplyASetOfChangesToInventory(t *testing.T) {
	inventoryItemChangeSet := InventoryItemChangeSet{}
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("iron_ore", 1))
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("copper_ore", -2))

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{3, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{4, "Copper Ore", "copper_ore"})

	assert.Equal(t, 3, inventory.items["iron_ore"].count)
	assert.Equal(t, 4, inventory.items["copper_ore"].count)

	_, inventory = ApplyInventoryItemChangeSet(inventory, inventoryItemChangeSet)

	assert.Equal(t, 4, inventory.items["iron_ore"].count)
	assert.Equal(t, 2, inventory.items["copper_ore"].count)
}

func TestThatWeCantApplyASetOfChangesIfCountWillBeNegative(t *testing.T) {
	inventoryItemChangeSet := InventoryItemChangeSet{}
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("iron_ore", 1))
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("copper_ore", -10))

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{3, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{4, "Copper Ore", "copper_ore"})

	assert.Equal(t, 3, inventory.items["iron_ore"].count)
	assert.Equal(t, 4, inventory.items["copper_ore"].count)

	wasApplied, newInventory := ApplyInventoryItemChangeSet(inventory, inventoryItemChangeSet)

	assert.False(t, wasApplied)
	assert.Equal(t, 3, newInventory.items["iron_ore"].count)
	assert.Equal(t, 4, newInventory.items["copper_ore"].count)
}

func TestThatApplyInventoryItemChangeReturnsFalseIfFailed(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", -2)
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Iron ore", "iron_ore"})

	wasChanged, newInventory := ApplyInventoryItemChange(inventory, inventoryChange)

	assert.False(t, wasChanged)
	assert.Equal(t, 1, newInventory.items["iron_ore"].count)
}

func TestThatApplyInventoryItemChangeReturnsCorrectResultIfNewCountIs0(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", -1)
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Iron ore", "iron_ore"})

	wasChanged, newInventory := ApplyInventoryItemChange(inventory, inventoryChange)

	assert.True(t, wasChanged)
	assert.Equal(t, 0, newInventory.items["iron_ore"].count)
}

func TestThatGetNewCountAfterInventoryItemChangeReturnsCorrectCount(t *testing.T) {
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Iron ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{2, "Copper ore", "copper_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{3, "Copper plates", "copper_plates"})

	inventoryChange := NewInventoryChange("iron_ore", -1)
	assert.Equal(t, 0, GetNewCountAfterInventoryItemChange(inventory, inventoryChange))

	inventoryChange = NewInventoryChange("copper_ore", -3)
	assert.Equal(t, -1, GetNewCountAfterInventoryItemChange(inventory, inventoryChange))

	inventoryChange = NewInventoryChange("copper_plates", 0)
	assert.Equal(t, 3, GetNewCountAfterInventoryItemChange(inventory, inventoryChange))
}

func TestThatApplyInventoryItemChangeSetDoesntChangeOriginal(t *testing.T) {
	inventoryItemChangeSet := InventoryItemChangeSet{}
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("iron_ore", 1))
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("copper_ore", -20))

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{2, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{30, "Copper Ore", "copper_ore"})

	assert.Equal(t, 2, inventory.items["iron_ore"].count)
	assert.Equal(t, 30, inventory.items["copper_ore"].count)

	wasChanged, newInventory := ApplyInventoryItemChangeSet(inventory, inventoryItemChangeSet)

	assert.True(t, wasChanged)
	assert.Equal(t, 2, inventory.items["iron_ore"].count)
	assert.Equal(t, 30, inventory.items["copper_ore"].count)
	assert.Equal(t, 3, newInventory.items["iron_ore"].count)
	assert.Equal(t, 10, newInventory.items["copper_ore"].count)
}

func TestThatTestDeepCopyInventoryItemActuallyCopiesMap(t *testing.T) {
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, InventoryItem{0, "Iron Ore", "iron_ore"})
	inventory = AddInventoryItem(inventory, InventoryItem{1, "Copper Ore", "copper_ore"})

	newInventory := deepCopyInventory(inventory)

	assert.Equal(t, 0, newInventory.items["iron_ore"].count)
	assert.Equal(t, 1, newInventory.items["copper_ore"].count)
	assert.Equal(t, 0, inventory.items["iron_ore"].count)
	assert.Equal(t, 1, inventory.items["copper_ore"].count)

	newInventory = changeCountOfInventoryItemWithId(newInventory, "iron_ore", 3)

	assert.Equal(t, 3, newInventory.items["iron_ore"].count)
	assert.Equal(t, 1, newInventory.items["copper_ore"].count)

	assert.Equal(t, 0, inventory.items["iron_ore"].count)
	assert.Equal(t, 1, inventory.items["copper_ore"].count)
}
