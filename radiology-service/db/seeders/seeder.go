package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SeedExam struct {
	Code                    string  `json:"code"`
	Name                    string  `json:"name"`
	Modality                int32   `json:"modality"`
	BodySite                int32   `json:"body_site"`
	LOINCCode               string  `json:"loinc_code"`
	RequiresContrast        bool    `json:"requires_contrast"`
	Price                   float64 `json:"price"`
	PreparationInstructions string  `json:"preparation_instructions"`
}

type SeedData struct {
	Exams []SeedExam `json:"exams"`
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
