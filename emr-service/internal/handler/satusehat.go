package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

func (h *EmrHandler) SyncSatusehatID(ctx context.Context, req *emrv1.SyncSatusehatIDRequest) (*emrv1.SyncSatusehatIDResponse, error) {
	if h.satusehat == nil {
		return nil, status.Error(codes.Unimplemented, "satusehat synchronization is not configured")
	}

	err := h.satusehat.SyncSatusehatID(
		ctx,
		domain.SatusehatResourceType(req.GetResourceType()),
		req.GetResourceId(),
		req.GetSatusehatId(),
	)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &emrv1.SyncSatusehatIDResponse{Success: true}, nil
}
