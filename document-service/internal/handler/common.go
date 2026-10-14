package handler

import (
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/azuvicenna/selaras-services/document-service/internal/domain"
)

func toStatusError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrDocumentVoid):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrTemplateNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrCompilationFailed):
		return status.Error(codes.Internal, err.Error())
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

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func toTimestampPtr(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}
