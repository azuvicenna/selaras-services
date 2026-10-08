package handler

import (
	"context"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

func (h *EmrHandler) CreateObservation(ctx context.Context, req *emrv1.CreateObservationRequest) (*emrv1.CreateObservationResponse, error) {
	obs, err := h.observation.CreateObservation(ctx, fromProtoObservation(req.GetObservation()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.CreateObservationResponse{Observation: toProtoObservation(obs)}, nil
}

func (h *EmrHandler) GetEncounterObservations(ctx context.Context, req *emrv1.GetEncounterObservationsRequest) (*emrv1.GetEncounterObservationsResponse, error) {
	items, err := h.observation.GetEncounterObservations(ctx, req.GetEncounterId(), domain.ObservationCategory(req.GetCategory()))
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.GetEncounterObservationsResponse{
		Observations: toProtoList(items, toProtoObservation),
	}, nil
}

func (h *EmrHandler) GetPatientObservations(ctx context.Context, req *emrv1.GetPatientObservationsRequest) (*emrv1.GetPatientObservationsResponse, error) {
	items, total, err := h.observation.GetPatientObservations(ctx, domain.ObservationFilter{
		PatientID: req.GetPatientId(),
		Category:  domain.ObservationCategory(req.GetCategory()),
		Code:      req.GetCode(),
		Limit:     int(req.GetLimit()),
		Offset:    int(req.GetOffset()),
	})
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.GetPatientObservationsResponse{
		Observations: toProtoList(items, toProtoObservation),
		TotalCount:   total,
	}, nil
}

func fromProtoObservation(o *emrv1.Observation) domain.Observation {
	obs := domain.Observation{
		ID:                       o.GetId(),
		PatientID:                o.GetPatientId(),
		EncounterID:              o.GetEncounterId(),
		PractitionerID:           o.GetPractitionerId(),
		Status:                   domain.ObservationStatus(o.GetStatus()),
		Category:                 domain.ObservationCategory(o.GetCategory()),
		SatusehatID:              o.GetSatusehatId(),
		Code:                     o.GetCode(),
		Name:                     o.GetName(),
		Unit:                     o.GetUnit(),
		ReferenceRange:           o.GetReferenceRange(),
		Interpretation:           o.GetInterpretation(),
		BodySite:                 o.GetBodySite(),
		Method:                   o.GetMethod(),
		Notes:                    o.GetNotes(),
		AmendedFromObservationID: o.GetAmendedFromObservationId(),
		EffectiveTime:            fromTimestamp(o.GetEffectiveTime()),
	}

	switch v := o.GetValue().(type) {
	case *emrv1.Observation_ValueQuantity:
		val := v.ValueQuantity
		obs.ValueQuantity = &val
	case *emrv1.Observation_ValueString:
		obs.ValueString = v.ValueString
	case *emrv1.Observation_ValueBoolean:
		val := v.ValueBoolean
		obs.ValueBoolean = &val
	case *emrv1.Observation_ValueCode:
		obs.ValueCode = v.ValueCode
	}

	for _, c := range o.GetComponents() {
		comp := domain.ObservationComponent{
			Code:           c.GetCode(),
			Name:           c.GetName(),
			Unit:           c.GetUnit(),
			Interpretation: c.GetInterpretation(),
			ReferenceRange: c.GetReferenceRange(),
		}
		switch cv := c.GetValue().(type) {
		case *emrv1.ObservationComponent_ValueQuantity:
			val := cv.ValueQuantity
			comp.ValueQuantity = &val
		case *emrv1.ObservationComponent_ValueString:
			comp.ValueString = cv.ValueString
		case *emrv1.ObservationComponent_ValueBoolean:
			val := cv.ValueBoolean
			comp.ValueBoolean = &val
		}
		obs.Components = append(obs.Components, comp)
	}

	return obs
}

func toProtoObservation(o *domain.Observation) *emrv1.Observation {
	out := &emrv1.Observation{
		Id:                       o.ID,
		PatientId:                o.PatientID,
		EncounterId:              o.EncounterID,
		PractitionerId:           o.PractitionerID,
		Status:                   emrv1.ObservationStatus(o.Status),
		Category:                 emrv1.ObservationCategory(o.Category),
		SatusehatId:              o.SatusehatID,
		Code:                     o.Code,
		Name:                     o.Name,
		Unit:                     o.Unit,
		ReferenceRange:           o.ReferenceRange,
		Interpretation:           o.Interpretation,
		BodySite:                 o.BodySite,
		Method:                   o.Method,
		Notes:                    o.Notes,
		AmendedFromObservationId: o.AmendedFromObservationID,
		EffectiveTime:            toTimestamp(o.EffectiveTime),
		CreatedAt:                toTimestamp(o.CreatedAt),
		UpdatedAt:                toTimestamp(o.UpdatedAt),
	}

	switch {
	case o.ValueQuantity != nil:
		out.Value = &emrv1.Observation_ValueQuantity{ValueQuantity: *o.ValueQuantity}
	case o.ValueBoolean != nil:
		out.Value = &emrv1.Observation_ValueBoolean{ValueBoolean: *o.ValueBoolean}
	case o.ValueCode != "":
		out.Value = &emrv1.Observation_ValueCode{ValueCode: o.ValueCode}
	case o.ValueString != "":
		out.Value = &emrv1.Observation_ValueString{ValueString: o.ValueString}
	}

	for _, c := range o.Components {
		comp := &emrv1.ObservationComponent{
			Code:           c.Code,
			Name:           c.Name,
			Unit:           c.Unit,
			Interpretation: c.Interpretation,
			ReferenceRange: c.ReferenceRange,
		}
		switch {
		case c.ValueQuantity != nil:
			comp.Value = &emrv1.ObservationComponent_ValueQuantity{ValueQuantity: *c.ValueQuantity}
		case c.ValueBoolean != nil:
			comp.Value = &emrv1.ObservationComponent_ValueBoolean{ValueBoolean: *c.ValueBoolean}
		case c.ValueString != "":
			comp.Value = &emrv1.ObservationComponent_ValueString{ValueString: c.ValueString}
		}
		out.Components = append(out.Components, comp)
	}

	return out
}
