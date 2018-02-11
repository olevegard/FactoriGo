package main

import "testing"

func TestThatWeCanUpdateInventory(t *testing.T) {
	gameState := MakeDefaultGameState()

	gameState = UpdateInventory(gameState)
}
