package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GameState struct {
}

func (state *GameState) Save(filename string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	fmt.Println("New High Score! Saving...")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (state *GameState) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, state)
}
