package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
	"github.com/azuvicenna/selaras-services/bff-web/internal/handler"
	"google.golang.org/grpc"
)

type mockMasterClient struct {
	masterv1.MasterServiceClient
}

func (m *mockMasterClient) ListProvinces(ctx context.Context, in *masterv1.ListProvincesRequest, opts ...grpc.CallOption) (*masterv1.ListProvincesResponse, error) {
	return &masterv1.ListProvincesResponse{
		Items: []*masterv1.Province{
			{Code: "31", Name: "DKI JAKARTA", IsActive: true},
			{Code: "32", Name: "JAWA BARAT", IsActive: true},
		},
		TotalCount: 2,
	}, nil
}

func (m *mockMasterClient) GetProvinceByCode(ctx context.Context, in *masterv1.GetProvinceByCodeRequest, opts ...grpc.CallOption) (*masterv1.GetProvinceByCodeResponse, error) {
	return &masterv1.GetProvinceByCodeResponse{
		Province: &masterv1.Province{Code: in.Code, Name: "DKI JAKARTA", IsActive: true},
	}, nil
}

func (m *mockMasterClient) ListCities(ctx context.Context, in *masterv1.ListCitiesRequest, opts ...grpc.CallOption) (*masterv1.ListCitiesResponse, error) {
	return &masterv1.ListCitiesResponse{
		Items: []*masterv1.City{
			{Code: "3171", ProvinceCode: "31", Name: "KOTA JAKARTA PUSAT", IsActive: true},
		},
		TotalCount: 1,
	}, nil
}

func (m *mockMasterClient) GetCityByCode(ctx context.Context, in *masterv1.GetCityByCodeRequest, opts ...grpc.CallOption) (*masterv1.GetCityByCodeResponse, error) {
	return &masterv1.GetCityByCodeResponse{
		City: &masterv1.City{Code: in.Code, ProvinceCode: "31", Name: "KOTA JAKARTA PUSAT", IsActive: true},
	}, nil
}

func (m *mockMasterClient) ListDistricts(ctx context.Context, in *masterv1.ListDistrictsRequest, opts ...grpc.CallOption) (*masterv1.ListDistrictsResponse, error) {
	return &masterv1.ListDistrictsResponse{
		Items: []*masterv1.District{
			{Code: "317101", CityCode: "3171", Name: "GAMBIR", IsActive: true},
		},
		TotalCount: 1,
	}, nil
}

func (m *mockMasterClient) GetDistrictByCode(ctx context.Context, in *masterv1.GetDistrictByCodeRequest, opts ...grpc.CallOption) (*masterv1.GetDistrictByCodeResponse, error) {
	return &masterv1.GetDistrictByCodeResponse{
		District: &masterv1.District{Code: in.Code, CityCode: "3171", Name: "GAMBIR", IsActive: true},
	}, nil
}

func (m *mockMasterClient) ListVillages(ctx context.Context, in *masterv1.ListVillagesRequest, opts ...grpc.CallOption) (*masterv1.ListVillagesResponse, error) {
	return &masterv1.ListVillagesResponse{
		Items: []*masterv1.Village{
			{Code: "3171011001", DistrictCode: "317101", Name: "GAMBIR", PostalCode: "10110", IsActive: true},
		},
		TotalCount: 1,
	}, nil
}

func (m *mockMasterClient) GetVillageByCode(ctx context.Context, in *masterv1.GetVillageByCodeRequest, opts ...grpc.CallOption) (*masterv1.GetVillageByCodeResponse, error) {
	return &masterv1.GetVillageByCodeResponse{
		Village: &masterv1.Village{Code: in.Code, DistrictCode: "317101", Name: "GAMBIR", PostalCode: "10110", IsActive: true},
	}, nil
}

func (m *mockMasterClient) ListICD10(ctx context.Context, in *masterv1.ListICD10Request, opts ...grpc.CallOption) (*masterv1.ListICD10Response, error) {
	return &masterv1.ListICD10Response{
		Items: []*masterv1.ICD10{
			{Code: "A09", DescriptionId: "Gastroenteritis", IsActive: true},
		},
		TotalCount: 1,
	}, nil
}

func (m *mockMasterClient) GetICD10ByCode(ctx context.Context, in *masterv1.GetICD10ByCodeRequest, opts ...grpc.CallOption) (*masterv1.GetICD10ByCodeResponse, error) {
	return &masterv1.GetICD10ByCodeResponse{
		Icd10: &masterv1.ICD10{Code: in.Code, DescriptionId: "Gastroenteritis", IsActive: true},
	}, nil
}

func (m *mockMasterClient) ListDepartments(ctx context.Context, in *masterv1.ListDepartmentsRequest, opts ...grpc.CallOption) (*masterv1.ListDepartmentsResponse, error) {
	return &masterv1.ListDepartmentsResponse{
		Items: []*masterv1.Department{
			{Id: "01J9B1X2K3M4N5P6Q7R8S9T0D1", Code: "POLI_UMUM", Name: "Poliklinik Umum", IsActive: true},
		},
		TotalCount: 1,
	}, nil
}

func (m *mockMasterClient) GetDepartmentByID(ctx context.Context, in *masterv1.GetDepartmentByIDRequest, opts ...grpc.CallOption) (*masterv1.GetDepartmentByIDResponse, error) {
	return &masterv1.GetDepartmentByIDResponse{
		Department: &masterv1.Department{Id: in.Id, Code: "POLI_UMUM", Name: "Poliklinik Umum", IsActive: true},
	}, nil
}

func (m *mockMasterClient) GetHospitalProfile(ctx context.Context, in *masterv1.GetHospitalProfileRequest, opts ...grpc.CallOption) (*masterv1.GetHospitalProfileResponse, error) {
	return &masterv1.GetHospitalProfileResponse{
		HospitalProfile: &masterv1.HospitalProfile{
			Name: "RS Selaras Sehat Medika",
		},
	}, nil
}

func setupMasterTest() *handler.MasterHandler {
	cm := &client.ClientManager{
		Master: &mockMasterClient{},
	}
	return handler.NewMasterHandler(cm)
}

func TestMasterHandler_Regions(t *testing.T) {
	h := setupMasterTest()

	t.Run("ListProvinces", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/provinces?page=1&page_size=10", nil)
		rec := httptest.NewRecorder()
		h.ListProvinces(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var resp map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}
		if resp["success"] != true {
			t.Errorf("expected success true, got %v", resp["success"])
		}
	})

	t.Run("GetProvinceByCode", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/provinces/31", nil)
		req.SetPathValue("code", "31")
		rec := httptest.NewRecorder()
		h.GetProvinceByCode(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("ListCities", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/cities?province_code=31", nil)
		rec := httptest.NewRecorder()
		h.ListCities(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("ListDistricts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/districts?city_code=3171", nil)
		rec := httptest.NewRecorder()
		h.ListDistricts(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("ListVillages", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/villages?district_code=317101", nil)
		rec := httptest.NewRecorder()
		h.ListVillages(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})
}

func TestMasterHandler_DepartmentsAndICD(t *testing.T) {
	h := setupMasterTest()

	t.Run("ListDepartments", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/departments", nil)
		rec := httptest.NewRecorder()
		h.ListDepartments(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("GetDepartmentByID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/departments/01J9B1X2K3M4N5P6Q7R8S9T0D1", nil)
		req.SetPathValue("id", "01J9B1X2K3M4N5P6Q7R8S9T0D1")
		rec := httptest.NewRecorder()
		h.GetDepartmentByID(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("ListICD10", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/icd10", nil)
		rec := httptest.NewRecorder()
		h.ListICD10(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("GetICD10ByCode", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/icd10/A09", nil)
		req.SetPathValue("code", "A09")
		rec := httptest.NewRecorder()
		h.GetICD10ByCode(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("GetHospitalProfile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/master/hospital-profile", nil)
		rec := httptest.NewRecorder()
		h.GetHospitalProfile(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})
}
