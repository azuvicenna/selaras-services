package handler

import (
	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	"github.com/azuvicenna/selaras-services/master-service/internal/domain"
)

type MasterHandler struct {
	masterv1.UnimplementedMasterServiceServer
	icd9cm     domain.ICD9CMUsecase
	icd10      domain.ICD10Usecase
	region     domain.RegionUsecase
	department domain.DepartmentUsecase
	hospital   domain.HospitalProfileUsecase
}

func NewMasterHandler(
	icd9cm domain.ICD9CMUsecase,
	icd10 domain.ICD10Usecase,
	region domain.RegionUsecase,
	department domain.DepartmentUsecase,
	hospital domain.HospitalProfileUsecase,
) *MasterHandler {
	return &MasterHandler{
		icd9cm:     icd9cm,
		icd10:      icd10,
		region:     region,
		department: department,
		hospital:   hospital,
	}
}
