package main

import (
	"encoding/json"
	"fmt"
	"os"

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
)

type SeedData struct {
	ICD9CM          []*masterv1.ICD9CM          `json:"icd9cm"`
	ICD10           []*masterv1.ICD10           `json:"icd10"`
	Provinces       []*masterv1.Province        `json:"provinces"`
	Cities          []*masterv1.City            `json:"cities"`
	Districts       []*masterv1.District        `json:"districts"`
	Villages        []*masterv1.Village         `json:"villages"`
	Departments     []*masterv1.Department      `json:"departments"`
	HospitalProfile *masterv1.HospitalProfile   `json:"hospital_profile"`
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
