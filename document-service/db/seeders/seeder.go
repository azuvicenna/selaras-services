package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SeedTemplate struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         int32  `json:"type"`
	TypstContent string `json:"typst_content"`
}

type SeedData struct {
	Templates []SeedTemplate `json:"templates"`
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
