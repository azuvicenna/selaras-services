package domain

import "time"

type Gender int32

const (
	GenderUnspecified Gender = 0
	GenderMale        Gender = 1
	GenderFemale      Gender = 2
)

type BloodType int32

const (
	BloodTypeUnspecified BloodType = 0
	BloodTypeA           BloodType = 1
	BloodTypeB           BloodType = 2
	BloodTypeAB          BloodType = 3
	BloodTypeO           BloodType = 4
)

type MaritalStatus int32

const (
	MaritalStatusUnspecified MaritalStatus = 0
	MaritalStatusSingle      MaritalStatus = 1
	MaritalStatusMarried     MaritalStatus = 2
	MaritalStatusDivorced    MaritalStatus = 3
	MaritalStatusWidowed     MaritalStatus = 4
)

type Religion int32

const (
	ReligionUnspecified Religion = 0
	ReligionIslam       Religion = 1
	ReligionProtestant  Religion = 2
	ReligionCatholic    Religion = 3
	ReligionHindu       Religion = 4
	ReligionBuddhist    Religion = 5
	ReligionConfucian   Religion = 6
	ReligionOther       Religion = 7
)

type DisabilityType int32

const (
	DisabilityTypeUnspecified    DisabilityType = 0
	DisabilityTypeNone           DisabilityType = 1
	DisabilityTypeHearing        DisabilityType = 2
	DisabilityTypeVisual         DisabilityType = 3
	DisabilityTypeSpeech         DisabilityType = 4
	DisabilityTypePhysical       DisabilityType = 5
	DisabilityTypeMental         DisabilityType = 6
	DisabilityTypeColorBlindness DisabilityType = 7
	DisabilityTypeOther          DisabilityType = 8
)

type PatientStatus int32

const (
	PatientStatusUnspecified PatientStatus = 0
	PatientStatusActive      PatientStatus = 1
	PatientStatusInactive    PatientStatus = 2
	PatientStatusDeceased    PatientStatus = 3
)

type Patient struct {
	ID                  string         `json:"id"`
	MedicalRecordNo     string         `json:"medical_record_no"`
	SatusehatID         string         `json:"satusehat_id"`
	NIK                 string         `json:"nik"`
	Name                string         `json:"name"`
	MotherName          string         `json:"mother_name"`
	BirthPlace          string         `json:"birth_place"`
	BirthDate           time.Time      `json:"birth_date"`
	Gender              Gender         `json:"gender"`
	BloodType           BloodType      `json:"blood_type"`
	MaritalStatus       MaritalStatus  `json:"marital_status"`
	Religion            Religion       `json:"religion"`
	Phone               string         `json:"phone"`
	Email               string         `json:"email"`
	Address             string         `json:"address"`
	VillageCode         string         `json:"village_code"`
	DistrictCode        string         `json:"district_code"`
	CityCode            string         `json:"city_code"`
	ProvinceCode        string         `json:"province_code"`
	PostalCode          string         `json:"postal_code"`
	RT                  string         `json:"rt"`
	RW                  string         `json:"rw"`
	Occupation          string         `json:"occupation"`
	Education           string         `json:"education"`
	PreferredLanguage   string         `json:"preferred_language"`
	DisabilityType      DisabilityType `json:"disability_type"`
	SpecialNeedsNote    string         `json:"special_needs_note"`
	PhotoURL            string         `json:"photo_url"`
	FingerprintTemplate []byte         `json:"-"`
	Status              PatientStatus  `json:"status"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}