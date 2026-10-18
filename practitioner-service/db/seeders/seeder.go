package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SeedPractitioner struct {
	ID               string `json:"id"`
	NIK              string `json:"nik"`
	Name             string `json:"name"`
	Gender           int32  `json:"gender"`
	BirthDate        string `json:"birth_date"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	PractitionerType int32  `json:"practitioner_type"`
	Specialization   string `json:"specialization"`
	SIPNumber        string `json:"sip_number"`
	SIPExpiryDate    string `json:"sip_expiry_date"`
	STRNumber        string `json:"str_number"`
	STRExpiryDate    string `json:"str_expiry_date"`
	DepartmentID     string `json:"department_id"`
	IHSNumber        string `json:"ihs_number"`
	IsActive         bool   `json:"is_active"`
}

type SeedData struct {
	Practitioners []SeedPractitioner `json:"practitioners"`
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
