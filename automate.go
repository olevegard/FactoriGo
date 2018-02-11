package main

import (
	"encoding/json"
	"io/ioutil"
)

func UpdateInventory(gameState GameState) GameState {
	gameState.CurrentProduction.IronMines, gameState.CurrentInventory = CreateNewBatchIfTimeBecomes0(gameState.CurrentProduction.IronMines, gameState.CurrentInventory)
	gameState.CurrentProduction.CopperMines, gameState.CurrentInventory = CreateNewBatchIfTimeBecomes0(gameState.CurrentProduction.CopperMines, gameState.CurrentInventory)

	gameState.CurrentProduction.IronSmelters, gameState.CurrentInventory = CreateNewBatchIfTimeBecomes0(gameState.CurrentProduction.IronSmelters, gameState.CurrentInventory)
	gameState.CurrentProduction.CopperSmelters, gameState.CurrentInventory = CreateNewBatchIfTimeBecomes0(gameState.CurrentProduction.CopperSmelters, gameState.CurrentInventory)
	return gameState
}

func ReadDefaultGameState() GameState {
	gameState := GameState{}
	data, _ := ioutil.ReadFile("defaultState.json")
	json.Unmarshal(data, &gameState)

	return gameState
}
