package main

import (
	"encoding/json"
	"io/ioutil"
)

func UpdateInventory(gameState GameState) GameState {
	for index, _ := range gameState.CurrentProduction.ProductionUnits {
		gameState.CurrentProduction.ProductionUnits[index], gameState.CurrentInventory = CreateNewBatchIfTimeBecomes0(gameState.CurrentProduction.ProductionUnits[index], gameState.CurrentInventory)
	}

	return gameState
}

func ReadDefaultGameState() GameState {
	gameState := GameState{}
	data, _ := ioutil.ReadFile("defaultState.json")
	json.Unmarshal(data, &gameState)

	return gameState
}
