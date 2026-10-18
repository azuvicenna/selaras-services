package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SeedWardClass struct {
	Code                string `json:"code"`
	Name                string `json:"name"`
	BPJSEquivalentClass string `json:"bpjs_equivalent_class"`
	IsActive            bool   `json:"is_active"`
}

type SeedWard struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	WardClassCode string `json:"ward_class_code"`
	DepartmentID  string `json:"department_id"`
	Floor         int32  `json:"floor"`
	Capacity      int32  `json:"capacity"`
	IsActive      bool   `json:"is_active"`
}

type SeedBed struct {
	Code     string `json:"code"`
	WardCode string `json:"ward_code"`
	Status   int32  `json:"status"`
	IsActive bool   `json:"is_active"`
}

type SeedData struct {
	WardClasses []SeedWardClass `json:"ward_classes"`
	Wards       []SeedWard      `json:"wards"`
	Beds        []SeedBed       `json:"beds"`
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
