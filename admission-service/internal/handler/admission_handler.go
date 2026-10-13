package handler

import (
	admissionv1 "github.com/azuvicenna/selaras-services/admission-service/gen/admission/v1"
	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

type AdmissionHandler struct {
	admissionv1.UnimplementedAdmissionServiceServer
	schedule    domain.ScheduleUsecase
	appointment domain.AppointmentUsecase
	admission   domain.AdmissionUsecase
}

func NewAdmissionHandler(
	schedule domain.ScheduleUsecase,
	appointment domain.AppointmentUsecase,
	admission domain.AdmissionUsecase,
) *AdmissionHandler {
	return &AdmissionHandler{
		schedule:    schedule,
		appointment: appointment,
		admission:   admission,
	}
}
