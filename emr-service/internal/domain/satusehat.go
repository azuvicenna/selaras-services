package domain

import "context"

type SatusehatResourceType int32

const (
	SatusehatResourceTypeUnspecified       SatusehatResourceType = 0
	SatusehatResourceTypeEncounter         SatusehatResourceType = 1
	SatusehatResourceTypeObservation       SatusehatResourceType = 2
	SatusehatResourceTypeClinicalNote      SatusehatResourceType = 3
	SatusehatResourceTypeCondition         SatusehatResourceType = 4
	SatusehatResourceTypeProcedure         SatusehatResourceType = 5
	SatusehatResourceTypeMedicationRequest SatusehatResourceType = 6
	SatusehatResourceTypeDiagnosticReport  SatusehatResourceType = 7
	SatusehatResourceTypeAllergy           SatusehatResourceType = 8
)

type SatusehatUsecase interface {
	SyncSatusehatID(ctx context.Context, resourceType SatusehatResourceType, resourceID, satusehatID string) error
}
