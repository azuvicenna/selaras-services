package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SeedTest struct {
	Code                 string  `json:"code"`
	Name                 string  `json:"name"`
	Category             int32   `json:"category"`
	LOINCCode            string  `json:"loinc_code"`
	SpecimenType         int32   `json:"specimen_type"`
	Unit                 string  `json:"unit"`
	ReferenceRangeMale   string  `json:"reference_range_male"`
	ReferenceRangeFemale string  `json:"reference_range_female"`
	RefLow               float64 `json:"ref_low"`
	RefHigh              float64 `json:"ref_high"`
	CriticalLow          float64 `json:"critical_low"`
	CriticalHigh         float64 `json:"critical_high"`
	Price                float64 `json:"price"`
}

type SeedData struct {
	Tests []SeedTest `json:"tests"`
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
