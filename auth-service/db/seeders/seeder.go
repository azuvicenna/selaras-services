package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SeedUser struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FullName       string `json:"full_name"`
	Role           int32  `json:"role"`
	PractitionerID string `json:"practitioner_id"`
	IsActive       bool   `json:"is_active"`
}

type SeedData struct {
	Users []SeedUser `json:"users"`
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
