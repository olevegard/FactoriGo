package main

import (
	"testing"

	"github.com/stvp/assert"
)

func TestThatWeCanGetInventoryItem(t *testing.T) {
	inventory := NewInventory()
	inventoryItem := NewInventoryItem(1, "Iron ore", "iron_ore", false)
	inventory = AddInventoryItem(inventory, inventoryItem)

	returnedItem := inventory.Items[inventoryItem.Id]

	if returnedItem.Id != "iron_ore" {
		t.Fatal("Couldn't get item from Inventory")
	}

	assert.Equal(t, returnedItem, inventoryItem)
}

func TestThatChangeCountAndReturnNewReturnsUpdateInventoryItem(t *testing.T) {
	inventory := NewInventory()
	inventoryItem := NewInventoryItem(1, "Iron ore", "iron_ore", false)
	inventory = AddInventoryItem(inventory, inventoryItem)

	inventory.Items[inventoryItem.Id] = changeCountAndReturnNew(inventory.Items["iron_ore"], 4)

	assert.Equal(t, 4, inventory.Items["iron_ore"].ItemCount)
}

func TestThatChangeCountOfInventoryItemWithIdUpdatesInventory(t *testing.T) {
	inventory := NewInventory()
	inventoryItem := NewInventoryItem(1, "Iron ore", "iron_ore", false)
	inventory = AddInventoryItem(inventory, inventoryItem)
	newInventory := changeCountOfInventoryItemWithId(inventory, "iron_ore", 2)

	assert.Equal(t, 2, newInventory.Items["iron_ore"].ItemCount)
}

func TestThatChangeCountOfInventoryItemWithIdDoesntChangeOriginalInv(t *testing.T) {
	inventory := NewInventory()
	inventoryItem := NewInventoryItem(1, "Iron ore", "iron_ore", false)
	inventory = AddInventoryItem(inventory, inventoryItem)
	newInventory := changeCountOfInventoryItemWithId(inventory, "iron_ore", 2)

	assert.Equal(t, 2, newInventory.Items["iron_ore"].ItemCount)
}

func TestThatCreateInventoryChangeReturnsExpectedInventoryItemChange(t *testing.T) {
	expectedInventoryChange := InventoryItemChange{}
	expectedInventoryChange.ChangeAmount = 3
	expectedInventoryChange.InventoryItemId = "iron_ore"
	actuaInventoryChange := NewInventoryChange("iron_ore", 3)

	assert.Equal(t, expectedInventoryChange, actuaInventoryChange)
}

func TestThatApplyInventoryItemChangeDoesntChangeOtherValues(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", 3)
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(5, "Iron ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(10, "Copper ore", "copper_ore", false))

	wasApplied, newInventory := ApplyInventoryItemChange(inventory, inventoryChange)

	assert.True(t, wasApplied)
	assert.Equal(t, 8, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 10, newInventory.Items["copper_ore"].ItemCount)
	assert.Equal(t, 5, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 10, inventory.Items["copper_ore"].ItemCount)
}

func TestThatApplyInventoryItemChangeAppliesPositiveChange(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", 3)
	inventoryItem := NewInventoryItem(5, "Iron ore", "iron_ore", false)
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, inventoryItem)

	_, inventory = ApplyInventoryItemChange(inventory, inventoryChange)

	assert.Equal(t, 8, inventory.Items[inventoryItem.Id].ItemCount)
}

func TestThatApplyInventoryItemChangeAppliesNegativeChange(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", -3)
	inventoryItem := NewInventoryItem(5, "Iron ore", "iron_ore", false)
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, inventoryItem)

	_, inventory = ApplyInventoryItemChange(inventory, inventoryChange)

	assert.Equal(t, 2, inventory.Items[inventoryItem.Id].ItemCount)
}

func TestThatWeCanApplyASetOfChangesToInventory(t *testing.T) {
	inventoryItemChangeSet := InventoryItemChangeSet{}
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("iron_ore", 1))
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("copper_ore", -2))

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(3, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(4, "Copper Ore", "copper_ore", false))

	assert.Equal(t, 3, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 4, inventory.Items["copper_ore"].ItemCount)

	_, inventory = ApplyInventoryItemChangeSet(inventory, inventoryItemChangeSet)

	assert.Equal(t, 4, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 2, inventory.Items["copper_ore"].ItemCount)
}

func TestThatWeCantApplyASetOfChangesIfCountWillBeNegative(t *testing.T) {
	inventoryItemChangeSet := InventoryItemChangeSet{}
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("iron_ore", 1))
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("copper_ore", -10))

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(3, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(4, "Copper Ore", "copper_ore", false))

	assert.Equal(t, 3, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 4, inventory.Items["copper_ore"].ItemCount)

	wasApplied, newInventory := ApplyInventoryItemChangeSet(inventory, inventoryItemChangeSet)

	assert.False(t, wasApplied)
	assert.Equal(t, 3, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 4, newInventory.Items["copper_ore"].ItemCount)
}

func TestThatApplyInventoryItemChangeReturnsFalseIfFailed(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", -2)
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron ore", "iron_ore", false))

	wasChanged, newInventory := ApplyInventoryItemChange(inventory, inventoryChange)

	assert.False(t, wasChanged)
	assert.Equal(t, 1, newInventory.Items["iron_ore"].ItemCount)
}

func TestThatApplyInventoryItemChangeReturnsCorrectResultIfNewCountIs0(t *testing.T) {
	inventoryChange := NewInventoryChange("iron_ore", -1)
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron ore", "iron_ore", false))

	wasChanged, newInventory := ApplyInventoryItemChange(inventory, inventoryChange)

	assert.True(t, wasChanged)
	assert.Equal(t, 0, newInventory.Items["iron_ore"].ItemCount)
}

func TestThatGetNewCountAfterInventoryItemChangeReturnsCorrectCount(t *testing.T) {
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Iron ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(2, "Copper ore", "copper_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(3, "Copper plates", "copper_plate", false))

	inventoryChange := NewInventoryChange("iron_ore", -1)
	assert.Equal(t, 0, GetNewCountAfterInventoryItemChange(inventory, inventoryChange))

	inventoryChange = NewInventoryChange("copper_ore", -3)
	assert.Equal(t, -1, GetNewCountAfterInventoryItemChange(inventory, inventoryChange))

	inventoryChange = NewInventoryChange("copper_plate", 0)
	assert.Equal(t, 3, GetNewCountAfterInventoryItemChange(inventory, inventoryChange))
}

func TestThatApplyInventoryItemChangeSetDoesntChangeOriginal(t *testing.T) {
	inventoryItemChangeSet := InventoryItemChangeSet{}
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("iron_ore", 1))
	inventoryItemChangeSet = append(inventoryItemChangeSet, NewInventoryChange("copper_ore", -20))

	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(2, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(30, "Copper Ore", "copper_ore", false))

	assert.Equal(t, 2, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 30, inventory.Items["copper_ore"].ItemCount)

	wasChanged, newInventory := ApplyInventoryItemChangeSet(inventory, inventoryItemChangeSet)

	assert.True(t, wasChanged)
	assert.Equal(t, 2, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 30, inventory.Items["copper_ore"].ItemCount)
	assert.Equal(t, 3, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 10, newInventory.Items["copper_ore"].ItemCount)
}

func TestThatTestDeepCopyInventoryItemActuallyCopiesMap(t *testing.T) {
	inventory := NewInventory()
	inventory = AddInventoryItem(inventory, NewInventoryItem(0, "Iron Ore", "iron_ore", false))
	inventory = AddInventoryItem(inventory, NewInventoryItem(1, "Copper Ore", "copper_ore", false))

	newInventory := deepCopyInventory(inventory)

	assert.Equal(t, 0, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 1, newInventory.Items["copper_ore"].ItemCount)
	assert.Equal(t, 0, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 1, inventory.Items["copper_ore"].ItemCount)

	newInventory = changeCountOfInventoryItemWithId(newInventory, "iron_ore", 3)

	assert.Equal(t, 3, newInventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 1, newInventory.Items["copper_ore"].ItemCount)

	assert.Equal(t, 0, inventory.Items["iron_ore"].ItemCount)
	assert.Equal(t, 1, inventory.Items["copper_ore"].ItemCount)
}

func TestThatInventoryItemHasStringFunc(t *testing.T) {
	inventoryItem := NewInventoryItem(0, "Iron Ore", "iron_ore", false)
	assert.Equal(t, "Iron Ore", inventoryItem.String())
}

func TestThatInventoryItemHasCountFunc(t *testing.T) {
	inventoryItem := NewInventoryItem(1, "Iron Ore", "iron_ore", false)
	assert.Equal(t, 1, inventoryItem.Count())
}
