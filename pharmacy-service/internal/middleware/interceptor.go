package middleware

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	var reqID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ids := md.Get("x-request-id"); len(ids) > 0 {
			reqID = ids[0]
		}
	}

	resp, err := handler(ctx, req)

	attrs := []any{
		"method", info.FullMethod,
		"code", status.Code(err).String(),
		"duration", time.Since(start),
	}
	if reqID != "" {
		attrs = append(attrs, "request_id", reqID)
	}

	slog.Info("grpc request", attrs...)
	return resp, err
}

func Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic recovered",
				"method", info.FullMethod,
				"panic", r,
				"stack", string(debug.Stack()))
			err = status.Error(codes.Internal, "internal error")
		}
	}()
	return handler(ctx, req)
}
