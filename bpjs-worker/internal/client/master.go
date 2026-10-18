package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
)

// FetchHospitalProfile retrieves the canonical hospital profile from master-service via gRPC.
func FetchHospitalProfile(ctx context.Context, addr string) (*masterv1.HospitalProfile, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return nil, nil
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master-service at %s: %w", addr, err)
	}
	defer func() { _ = conn.Close() }()

	client := masterv1.NewMasterServiceClient(conn)
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := client.GetHospitalProfile(reqCtx, &masterv1.GetHospitalProfileRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get hospital profile: %w", err)
	}

	return resp.GetHospitalProfile(), nil
}
