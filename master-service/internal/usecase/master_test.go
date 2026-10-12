package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

// --- Mock ICD9CM Repository ---
type mockICD9CMRepo struct {
	items map[string]*domain.ICD9CM
}

func newMockICD9CMRepo() *mockICD9CMRepo {
	return &mockICD9CMRepo{items: make(map[string]*domain.ICD9CM)}
}

func (m *mockICD9CMRepo) Create(ctx context.Context, item *domain.ICD9CM) error {
	for _, it := range m.items {
		if it.Code == item.Code {
			return domain.ErrAlreadyExists
		}
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockICD9CMRepo) GetByID(ctx context.Context, id string) (*domain.ICD9CM, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (m *mockICD9CMRepo) GetByCode(ctx context.Context, code string) (*domain.ICD9CM, error) {
	for _, it := range m.items {
		if it.Code == code {
			return it, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockICD9CMRepo) List(ctx context.Context, filter domain.ICD9CMFilter) ([]domain.ICD9CM, int64, error) {
	var res []domain.ICD9CM
	for _, it := range m.items {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}

func (m *mockICD9CMRepo) Update(ctx context.Context, item *domain.ICD9CM) error {
	if _, ok := m.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockICD9CMRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

// --- Mock ICD10 Repository ---
type mockICD10Repo struct {
	items map[string]*domain.ICD10
}

func newMockICD10Repo() *mockICD10Repo {
	return &mockICD10Repo{items: make(map[string]*domain.ICD10)}
}

func (m *mockICD10Repo) Create(ctx context.Context, item *domain.ICD10) error {
	for _, it := range m.items {
		if it.Code == item.Code {
			return domain.ErrAlreadyExists
		}
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockICD10Repo) GetByID(ctx context.Context, id string) (*domain.ICD10, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (m *mockICD10Repo) GetByCode(ctx context.Context, code string) (*domain.ICD10, error) {
	for _, it := range m.items {
		if it.Code == code {
			return it, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockICD10Repo) List(ctx context.Context, filter domain.ICD10Filter) ([]domain.ICD10, int64, error) {
	var res []domain.ICD10
	for _, it := range m.items {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}

func (m *mockICD10Repo) Update(ctx context.Context, item *domain.ICD10) error {
	if _, ok := m.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockICD10Repo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

// --- Mock Department Repository ---
type mockDeptRepo struct {
	items map[string]*domain.Department
}

func newMockDeptRepo() *mockDeptRepo {
	return &mockDeptRepo{items: make(map[string]*domain.Department)}
}

func (m *mockDeptRepo) Create(ctx context.Context, item *domain.Department) error {
	for _, it := range m.items {
		if it.Code == item.Code {
			return domain.ErrAlreadyExists
		}
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockDeptRepo) GetByID(ctx context.Context, id string) (*domain.Department, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

func (m *mockDeptRepo) GetByCode(ctx context.Context, code string) (*domain.Department, error) {
	for _, it := range m.items {
		if it.Code == code {
			return it, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockDeptRepo) List(ctx context.Context, filter domain.DepartmentFilter) ([]domain.Department, int64, error) {
	var res []domain.Department
	for _, it := range m.items {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}

func (m *mockDeptRepo) Update(ctx context.Context, item *domain.Department) error {
	if _, ok := m.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockDeptRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.items, id)
	return nil
}

// --- Mock Region Repository ---
type mockRegionRepo struct {
	provinces map[string]*domain.Province
	cities    map[string]*domain.City
	districts map[string]*domain.District
	villages  map[string]*domain.Village
}

func newMockRegionRepo() *mockRegionRepo {
	return &mockRegionRepo{
		provinces: make(map[string]*domain.Province),
		cities:    make(map[string]*domain.City),
		districts: make(map[string]*domain.District),
		villages:  make(map[string]*domain.Village),
	}
}

func (m *mockRegionRepo) CreateProvince(ctx context.Context, item *domain.Province) error {
	if _, ok := m.provinces[item.Code]; ok {
		return domain.ErrAlreadyExists
	}
	m.provinces[item.Code] = item
	return nil
}
func (m *mockRegionRepo) GetProvinceByCode(ctx context.Context, code string) (*domain.Province, error) {
	item, ok := m.provinces[code]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}
func (m *mockRegionRepo) ListProvinces(ctx context.Context, filter domain.ProvinceFilter) ([]domain.Province, int64, error) {
	var res []domain.Province
	for _, it := range m.provinces {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}
func (m *mockRegionRepo) UpdateProvince(ctx context.Context, item *domain.Province) error {
	if _, ok := m.provinces[item.Code]; !ok {
		return domain.ErrNotFound
	}
	m.provinces[item.Code] = item
	return nil
}
func (m *mockRegionRepo) DeleteProvince(ctx context.Context, code string) error {
	if _, ok := m.provinces[code]; !ok {
		return domain.ErrNotFound
	}
	delete(m.provinces, code)
	return nil
}

func (m *mockRegionRepo) CreateCity(ctx context.Context, item *domain.City) error {
	if _, ok := m.cities[item.Code]; ok {
		return domain.ErrAlreadyExists
	}
	m.cities[item.Code] = item
	return nil
}
func (m *mockRegionRepo) GetCityByCode(ctx context.Context, code string) (*domain.City, error) {
	item, ok := m.cities[code]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}
func (m *mockRegionRepo) ListCities(ctx context.Context, filter domain.CityFilter) ([]domain.City, int64, error) {
	var res []domain.City
	for _, it := range m.cities {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}
func (m *mockRegionRepo) UpdateCity(ctx context.Context, item *domain.City) error {
	if _, ok := m.cities[item.Code]; !ok {
		return domain.ErrNotFound
	}
	m.cities[item.Code] = item
	return nil
}
func (m *mockRegionRepo) DeleteCity(ctx context.Context, code string) error {
	if _, ok := m.cities[code]; !ok {
		return domain.ErrNotFound
	}
	delete(m.cities, code)
	return nil
}

func (m *mockRegionRepo) CreateDistrict(ctx context.Context, item *domain.District) error {
	if _, ok := m.districts[item.Code]; ok {
		return domain.ErrAlreadyExists
	}
	m.districts[item.Code] = item
	return nil
}
func (m *mockRegionRepo) GetDistrictByCode(ctx context.Context, code string) (*domain.District, error) {
	item, ok := m.districts[code]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}
func (m *mockRegionRepo) ListDistricts(ctx context.Context, filter domain.DistrictFilter) ([]domain.District, int64, error) {
	var res []domain.District
	for _, it := range m.districts {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}
func (m *mockRegionRepo) UpdateDistrict(ctx context.Context, item *domain.District) error {
	if _, ok := m.districts[item.Code]; !ok {
		return domain.ErrNotFound
	}
	m.districts[item.Code] = item
	return nil
}
func (m *mockRegionRepo) DeleteDistrict(ctx context.Context, code string) error {
	if _, ok := m.districts[code]; !ok {
		return domain.ErrNotFound
	}
	delete(m.districts, code)
	return nil
}

func (m *mockRegionRepo) CreateVillage(ctx context.Context, item *domain.Village) error {
	if _, ok := m.villages[item.Code]; ok {
		return domain.ErrAlreadyExists
	}
	m.villages[item.Code] = item
	return nil
}
func (m *mockRegionRepo) GetVillageByCode(ctx context.Context, code string) (*domain.Village, error) {
	item, ok := m.villages[code]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}
func (m *mockRegionRepo) ListVillages(ctx context.Context, filter domain.VillageFilter) ([]domain.Village, int64, error) {
	var res []domain.Village
	for _, it := range m.villages {
		res = append(res, *it)
	}
	return res, int64(len(res)), nil
}
func (m *mockRegionRepo) UpdateVillage(ctx context.Context, item *domain.Village) error {
	if _, ok := m.villages[item.Code]; !ok {
		return domain.ErrNotFound
	}
	m.villages[item.Code] = item
	return nil
}
func (m *mockRegionRepo) DeleteVillage(ctx context.Context, code string) error {
	if _, ok := m.villages[code]; !ok {
		return domain.ErrNotFound
	}
	delete(m.villages, code)
	return nil
}

// --- Mock Hospital Profile Repository ---
type mockHospitalProfileRepo struct {
	profile *domain.HospitalProfile
}

func (m *mockHospitalProfileRepo) Get(ctx context.Context, id string) (*domain.HospitalProfile, error) {
	if m.profile == nil {
		return nil, domain.ErrNotFound
	}
	return m.profile, nil
}

func (m *mockHospitalProfileRepo) Upsert(ctx context.Context, item *domain.HospitalProfile) error {
	m.profile = item
	return nil
}

func (m *mockHospitalProfileRepo) Update(ctx context.Context, item *domain.HospitalProfile) error {
	if m.profile == nil {
		return domain.ErrNotFound
	}
	m.profile = item
	return nil
}

// -----------------------------------------------------------------------------
// TESTS
// -----------------------------------------------------------------------------

func TestICD9CMUsecase(t *testing.T) {
	repo := newMockICD9CMRepo()
	uc := NewICD9CMUsecase(repo)
	ctx := context.Background()

	// 1. Validation fail
	_, err := uc.CreateICD9CM(ctx, domain.ICD9CM{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	// 2. Success with ULID generation
	created, err := uc.CreateICD9CM(ctx, domain.ICD9CM{
		Code:        "89.07",
		Description: "General Medical Consultation",
		IsActive:    true,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if len(created.ID) != 26 {
		t.Fatalf("expected ULID length 26, got %d (id=%s)", len(created.ID), created.ID)
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatalf("timestamps should not be zero")
	}

	// 3. Get by ID
	got, err := uc.GetICD9CMByID(ctx, created.ID)
	if err != nil || got.Code != "89.07" {
		t.Fatalf("get by id failed: %v, got: %+v", err, got)
	}

	// 4. Duplicate code error
	_, err = uc.CreateICD9CM(ctx, domain.ICD9CM{
		Code:        "89.07",
		Description: "Duplicate",
	})
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}

	// 5. Update
	created.Description = "Updated Consultation"
	updated, err := uc.UpdateICD9CM(ctx, *created)
	if err != nil || updated.Description != "Updated Consultation" {
		t.Fatalf("update failed: %v", err)
	}

	// 6. Delete
	if err := uc.DeleteICD9CM(ctx, created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	_, err = uc.GetICD9CMByID(ctx, created.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestICD10Usecase(t *testing.T) {
	repo := newMockICD10Repo()
	uc := NewICD10Usecase(repo)
	ctx := context.Background()

	// 1. Validation fail
	_, err := uc.CreateICD10(ctx, domain.ICD10{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	// 2. Success with ULID generation
	created, err := uc.CreateICD10(ctx, domain.ICD10{
		Code:          "A09.0",
		DescriptionID: "Gastroenteritis infeksi",
		DescriptionEN: "Infectious gastroenteritis",
		Category:      "Infectious diseases",
		IsActive:      true,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if len(created.ID) != 26 {
		t.Fatalf("expected ULID length 26, got %d", len(created.ID))
	}

	// 3. Get by code
	byCode, err := uc.GetICD10ByCode(ctx, "A09.0")
	if err != nil || byCode.DescriptionID != "Gastroenteritis infeksi" {
		t.Fatalf("get by code failed: %v", err)
	}
}

func TestDepartmentUsecase(t *testing.T) {
	repo := newMockDeptRepo()
	uc := NewDepartmentUsecase(repo)
	ctx := context.Background()

	// 1. Validation fail
	_, err := uc.CreateDepartment(ctx, domain.Department{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	// 2. Create with ULID
	created, err := uc.CreateDepartment(ctx, domain.Department{
		Code:           "POLI_UMUM",
		Name:           "Poliklinik Umum",
		DepartmentType: domain.DepartmentTypeOutpatient,
		IsActive:       true,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if len(created.ID) != 26 {
		t.Fatalf("expected ULID length 26, got %d", len(created.ID))
	}

	// 3. Get by code
	byCode, err := uc.GetDepartmentByCode(ctx, "POLI_UMUM")
	if err != nil || byCode.Name != "Poliklinik Umum" {
		t.Fatalf("get by code failed: %v", err)
	}
}

func TestRegionUsecase(t *testing.T) {
	repo := newMockRegionRepo()
	uc := NewRegionUsecase(repo)
	ctx := context.Background()

	// 1. Province validation
	_, err := uc.CreateProvince(ctx, domain.Province{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	prov, err := uc.CreateProvince(ctx, domain.Province{
		Code:     "31",
		Name:     "DKI JAKARTA",
		IsActive: true,
	})
	if err != nil || prov.Code != "31" {
		t.Fatalf("create province failed: %v", err)
	}

	// 2. City validation
	city, err := uc.CreateCity(ctx, domain.City{
		Code:         "3171",
		ProvinceCode: "31",
		Name:         "KOTA JAKARTA PUSAT",
		IsActive:     true,
	})
	if err != nil || city.Code != "3171" {
		t.Fatalf("create city failed: %v", err)
	}

	// 3. District validation
	dist, err := uc.CreateDistrict(ctx, domain.District{
		Code:     "317101",
		CityCode: "3171",
		Name:     "GAMBIR",
		IsActive: true,
	})
	if err != nil || dist.Code != "317101" {
		t.Fatalf("create district failed: %v", err)
	}

	// 4. Village validation
	vil, err := uc.CreateVillage(ctx, domain.Village{
		Code:         "3171011001",
		DistrictCode: "317101",
		Name:         "GAMBIR",
		PostalCode:   "10110",
		IsActive:     true,
	})
	if err != nil || vil.PostalCode != "10110" {
		t.Fatalf("create village failed: %v", err)
	}
}

func TestHospitalProfileUsecase(t *testing.T) {
	repo := &mockHospitalProfileRepo{}
	uc := NewHospitalProfileUsecase(repo)
	ctx := context.Background()

	// 1. Validation fail
	_, err := uc.UpsertHospitalProfile(ctx, domain.HospitalProfile{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}

	// 2. Upsert with ULID
	created, err := uc.UpsertHospitalProfile(ctx, domain.HospitalProfile{
		Name:             "RS Selaras Sehat",
		BPJSFacilityCode: "0115R001",
		Address:          "Jakarta",
	})
	if err != nil {
		t.Fatalf("upsert failed: %v", err)
	}
	if len(created.ID) != 26 {
		t.Fatalf("expected ULID length 26, got %d", len(created.ID))
	}

	// 3. Get
	got, err := uc.GetHospitalProfile(ctx, created.ID)
	if err != nil || got.BPJSFacilityCode != "0115R001" {
		t.Fatalf("get profile failed: %v", err)
	}
}
