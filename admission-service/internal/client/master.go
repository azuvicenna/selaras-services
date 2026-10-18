package client

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/azuvicenna/selaras-services/admission-service/internal/domain"
	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
)

var _ domain.DepartmentVerifier = (*MasterClient)(nil)

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

func (c *MasterClient) VerifyActiveDepartment(ctx context.Context, departmentID string) error {
	if c == nil || c.client == nil {
		return nil
	}

	departmentID = strings.TrimSpace(departmentID)
	if departmentID == "" {
		return fmt.Errorf("%w: department id is required", domain.ErrInvalidInput)
	}

	resp, err := c.client.GetDepartmentByID(ctx, &masterv1.GetDepartmentByIDRequest{Id: departmentID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("%w: department %s not found", domain.ErrNotFound, departmentID)
		}
		return fmt.Errorf("failed to verify department: %w", err)
	}

	if resp == nil || resp.Department == nil {
		return fmt.Errorf("%w: department %s not found", domain.ErrNotFound, departmentID)
	}

	if !resp.Department.IsActive {
		return fmt.Errorf("%w: department %s is inactive", domain.ErrInvalidInput, departmentID)
	}

	return nil
}
