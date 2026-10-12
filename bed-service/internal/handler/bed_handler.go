package handler

import (
	bedv1 "github.com/azuvicenna/selaras-services/bed-service/gen/bed/v1"
	"github.com/azuvicenna/selaras-services/bed-service/internal/domain"
)

type BedHandler struct {
	bedv1.UnimplementedBedServiceServer
	wardClass domain.WardClassUsecase
	ward      domain.WardUsecase
	bed       domain.BedUsecase
}

func NewBedHandler(
	wardClass domain.WardClassUsecase,
	ward domain.WardUsecase,
	bed domain.BedUsecase,
) *BedHandler {
	return &BedHandler{
		wardClass: wardClass,
		ward:      ward,
		bed:       bed,
	}
}
