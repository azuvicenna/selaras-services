package handler

import (
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
)

func toStatusError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrQuotaExceeded):
		return status.Error(codes.ResourceExhausted, err.Error())
	case errors.Is(err, domain.ErrInvalidStatusTransition):
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	slog.Error("unexpected internal error", "err", err)
	return status.Error(codes.Internal, "internal error")
}

func toProtoList[D, P any](items []D, convert func(*D) *P) []*P {
	out := make([]*P, len(items))
	for i := range items {
		out[i] = convert(&items[i])
	}
	return out
}

func fromTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
