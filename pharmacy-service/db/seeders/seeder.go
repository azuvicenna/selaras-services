package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SeedMedication struct {
	KFACode          string  `json:"kfa_code"`
	Name             string  `json:"name"`
	Form             int32   `json:"form"`
	Strength         string  `json:"strength"`
	Unit             string  `json:"unit"`
	PackageUnit      string  `json:"package_unit"`
	ConversionFactor int32   `json:"conversion_factor"`
	SellingPrice     float64 `json:"selling_price"`
	PurchasePrice    float64 `json:"purchase_price"`
	MinStock         int32   `json:"min_stock"`
}

type SeedData struct {
	Medications []SeedMedication `json:"medications"`
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
