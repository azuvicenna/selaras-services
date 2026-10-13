package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SeedTariff struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Category    int32   `json:"category"`
	ClassID     string  `json:"class_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type SeedData struct {
	Tariffs []SeedTariff `json:"tariffs"`
}

func loadSeedData(path string) (*SeedData, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	var data SeedData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("unmarshal seed json: %w", err)
	}

	return &data, nil
}
