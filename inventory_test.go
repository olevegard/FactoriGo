package main

import (
	"github.com/stvp/assert"
	"testing"
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

	inventory.items[inventoryItem.id] = ChangeCountAndReturnNew(inventory.items["iron_ore"], 4)

	assert.Equal(t, 4, inventory.items["iron_ore"].count)
}

func TestThatChangeCountOfInventoryItemWithIdUpdatesInventory(t *testing.T) {
	inventory := NewInventory()
	inventoryItem := InventoryItem{1, "Iron ore", "iron_ore"}
	inventory = AddInventoryItem(inventory, inventoryItem)
	inventory = ChangeCountOfInventoryItemWithId(inventory, "iron_ore", 2)

	assert.Equal(t, 2, inventory.items["iron_ore"].count)
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

	inventory = ApplyInventoryItemChange(inventory, inventoryChange)

	assert.Equal(t, 8, inventory.items["iron_ore"].count)
	assert.Equal(t, 10, inventory.items["copper_ore"].count)
}

func TestThatApplyInventoryItemChangeAppliesPositiveChange(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", 3)
	inventoryItem := InventoryItem{5, "Iron ore", "iron_ore"}
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, inventoryItem)

	inventory = ApplyInventoryItemChange(inventory, inventoryChange)

	assert.Equal(t, 8, inventory.items[inventoryItem.id].count)
}

func TestThatApplyInventoryItemChangeAppliesNegativeChange(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", -3)
	inventoryItem := InventoryItem{5, "Iron ore", "iron_ore"}
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, inventoryItem)

	inventory = ApplyInventoryItemChange(inventory, inventoryChange)

	assert.Equal(t, 2, inventory.items[inventoryItem.id].count)
}
