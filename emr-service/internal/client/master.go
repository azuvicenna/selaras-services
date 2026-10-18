package client

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
)

var _ domain.ICDVerifier = (*MasterClient)(nil)

type MasterClient struct {
	conn   *grpc.ClientConn
	client masterv1.MasterServiceClient
}

func NewMasterClient(addr string) (*MasterClient, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return nil, nil
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master-service at %s: %w", addr, err)
	}

	return &MasterClient{
		conn:   conn,
		client: masterv1.NewMasterServiceClient(conn),
	}, nil
}

func (c *MasterClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *MasterClient) VerifyAndLookupICD10(ctx context.Context, code string) (string, error) {
	if c == nil || c.client == nil {
		return "", nil
	}

	code = strings.TrimSpace(code)
	if code == "" {
		return "", fmt.Errorf("%w: icd10 code is required", domain.ErrInvalidInput)
	}

	resp, err := c.client.GetICD10ByCode(ctx, &masterv1.GetICD10ByCodeRequest{Code: code})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", fmt.Errorf("%w: icd10 code %s not found", domain.ErrNotFound, code)
		}
		return "", fmt.Errorf("failed to verify icd10 code: %w", err)
	}

	if resp == nil || resp.Icd10 == nil {
		return "", fmt.Errorf("%w: icd10 code %s not found", domain.ErrNotFound, code)
	}

	if !resp.Icd10.IsActive {
		return "", fmt.Errorf("%w: icd10 code %s is inactive", domain.ErrInvalidInput, code)
	}

	name := strings.TrimSpace(resp.Icd10.DescriptionId)
	if name == "" {
		name = strings.TrimSpace(resp.Icd10.DescriptionEn)
	}

	return name, nil
}

func (c *MasterClient) VerifyAndLookupICD9CM(ctx context.Context, code string) (string, error) {
	if c == nil || c.client == nil {
		return "", nil
	}

	code = strings.TrimSpace(code)
	if code == "" {
		return "", fmt.Errorf("%w: icd9cm code is required", domain.ErrInvalidInput)
	}

	resp, err := c.client.GetICD9CMByCode(ctx, &masterv1.GetICD9CMByCodeRequest{Code: code})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", fmt.Errorf("%w: icd9cm code %s not found", domain.ErrNotFound, code)
		}
		return "", fmt.Errorf("failed to verify icd9cm code: %w", err)
	}

	if resp == nil || resp.Icd9Cm == nil {
		return "", fmt.Errorf("%w: icd9cm code %s not found", domain.ErrNotFound, code)
	}

	if !resp.Icd9Cm.IsActive {
		return "", fmt.Errorf("%w: icd9cm code %s is inactive", domain.ErrInvalidInput, code)
	}

	return strings.TrimSpace(resp.Icd9Cm.Description), nil
}
