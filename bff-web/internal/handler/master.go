package handler

import (
	"net/http"
	"strconv"

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
)

type MasterHandler struct {
	cm *client.ClientManager
}

func NewMasterHandler(cm *client.ClientManager) *MasterHandler {
	return &MasterHandler{cm: cm}
}

// --- 1. ICD-10 (Diagnosa) ---

func (h *MasterHandler) ListICD10(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	query := r.URL.Query().Get("query")
	category := r.URL.Query().Get("category")

	resp, err := h.cm.Master.ListICD10(r.Context(), &masterv1.ListICD10Request{
		Query:    query,
		Category: category,
		Limit:    int32(pageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *MasterHandler) GetICD10ByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "code is required")
		return
	}

	resp, err := h.cm.Master.GetICD10ByCode(r.Context(), &masterv1.GetICD10ByCodeRequest{Code: code})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Icd10)
}

// --- 2. ICD-9-CM (Prosedur / Tindakan) ---

func (h *MasterHandler) ListICD9CM(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	query := r.URL.Query().Get("query")

	resp, err := h.cm.Master.ListICD9CM(r.Context(), &masterv1.ListICD9CMRequest{
		Query:  query,
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *MasterHandler) GetICD9CMByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "code is required")
		return
	}

	resp, err := h.cm.Master.GetICD9CMByCode(r.Context(), &masterv1.GetICD9CMByCodeRequest{Code: code})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Icd9Cm)
}

// --- 3. Departemen / Unit / Poliklinik ---

func (h *MasterHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	query := r.URL.Query().Get("query")
	typeInt, _ := strconv.Atoi(r.URL.Query().Get("type"))

	resp, err := h.cm.Master.ListDepartments(r.Context(), &masterv1.ListDepartmentsRequest{
		Query:          query,
		DepartmentType: masterv1.DepartmentType(typeInt),
		Limit:          int32(pageSize),
		Offset:         int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *MasterHandler) GetDepartmentByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Master.GetDepartmentByID(r.Context(), &masterv1.GetDepartmentByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Department)
}

// --- 4. Hospital Profile ---

func (h *MasterHandler) GetHospitalProfile(w http.ResponseWriter, r *http.Request) {
	resp, err := h.cm.Master.GetHospitalProfile(r.Context(), &masterv1.GetHospitalProfileRequest{})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.HospitalProfile)
}

// --- 5. Wilayah: Provinsi ---

func (h *MasterHandler) ListProvinces(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	query := r.URL.Query().Get("query")

	resp, err := h.cm.Master.ListProvinces(r.Context(), &masterv1.ListProvincesRequest{
		Query:  query,
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *MasterHandler) GetProvinceByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "code is required")
		return
	}

	resp, err := h.cm.Master.GetProvinceByCode(r.Context(), &masterv1.GetProvinceByCodeRequest{Code: code})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Province)
}

// --- 6. Wilayah: Kota / Kabupaten ---

func (h *MasterHandler) ListCities(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	query := r.URL.Query().Get("query")
	provinceCode := r.URL.Query().Get("province_code")

	resp, err := h.cm.Master.ListCities(r.Context(), &masterv1.ListCitiesRequest{
		ProvinceCode: provinceCode,
		Query:        query,
		Limit:        int32(pageSize),
		Offset:       int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *MasterHandler) GetCityByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "code is required")
		return
	}

	resp, err := h.cm.Master.GetCityByCode(r.Context(), &masterv1.GetCityByCodeRequest{Code: code})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.City)
}

// --- 7. Wilayah: Kecamatan ---

func (h *MasterHandler) ListDistricts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	query := r.URL.Query().Get("query")
	cityCode := r.URL.Query().Get("city_code")

	resp, err := h.cm.Master.ListDistricts(r.Context(), &masterv1.ListDistrictsRequest{
		CityCode: cityCode,
		Query:    query,
		Limit:    int32(pageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *MasterHandler) GetDistrictByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "code is required")
		return
	}

	resp, err := h.cm.Master.GetDistrictByCode(r.Context(), &masterv1.GetDistrictByCodeRequest{Code: code})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.District)
}

// --- 8. Wilayah: Kelurahan / Desa ---

func (h *MasterHandler) ListVillages(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	query := r.URL.Query().Get("query")
	districtCode := r.URL.Query().Get("district_code")
	postalCode := r.URL.Query().Get("postal_code")

	resp, err := h.cm.Master.ListVillages(r.Context(), &masterv1.ListVillagesRequest{
		DistrictCode: districtCode,
		PostalCode:   postalCode,
		Query:        query,
		Limit:        int32(pageSize),
		Offset:       int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Items, page, pageSize, resp.TotalCount)
}

func (h *MasterHandler) GetVillageByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "code is required")
		return
	}

	resp, err := h.cm.Master.GetVillageByCode(r.Context(), &masterv1.GetVillageByCodeRequest{Code: code})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.Village)
}
