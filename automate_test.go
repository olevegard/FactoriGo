package main

import "testing"

func TestThatWeCanUpdateInventory(t *testing.T) {
	gameState := ReadDefaultGameState()
	gameState = UpdateInventory(gameState)
}
