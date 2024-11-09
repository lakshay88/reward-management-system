package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type CategoryMultipliers map[string]int

var multipliers CategoryMultipliers

func init() {
	err := loadMultipliers()
	if err != nil {
		log.Fatalf("Error loading multipliers: %v", err)
	}
}

func loadMultipliers() error {
	file, err := os.Open("list.json")
	if err != nil {
		return fmt.Errorf("could not open multipliers file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read multipliers file: %v", err)
	}

	err = json.Unmarshal(data, &multipliers)
	if err != nil {
		return fmt.Errorf("could not parse multipliers JSON: %v", err)
	}

	fmt.Println("Multipliers Loaded:", multipliers)

	return nil
}

func GetCategoryMultiplier(category string) int {

	if multiplier, exists := multipliers[category]; exists {
		return multiplier
	}
	// default 1
	return 1
}
