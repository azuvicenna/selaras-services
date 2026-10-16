package client

import (
	"context"
	"fmt"

	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EMRClient interface {
	SyncSatusehatID(ctx context.Context, resType emrv1.SatusehatResourceType, localID, satusehatID string) error
	Close() error
}

type emrClientImpl struct {
	conn      *grpc.ClientConn
	grpcEmr   emrv1.EmrServiceClient
}

func NewEMRClient(addr string) (EMRClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to emr service at %s: %w", addr, err)
	}

	return &emrClientImpl{
		conn:    conn,
		grpcEmr: emrv1.NewEmrServiceClient(conn),
	}, nil
}

func (c *emrClientImpl) SyncSatusehatID(ctx context.Context, resType emrv1.SatusehatResourceType, localID, satusehatID string) error {
	req := &emrv1.SyncSatusehatIDRequest{
		ResourceType: resType,
		ResourceId:   localID,
		SatusehatId:  satusehatID,
	}

	res, err := c.grpcEmr.SyncSatusehatID(ctx, req)
	if err != nil {
		return fmt.Errorf("grpc SyncSatusehatID failed: %w", err)
	}

	if !res.GetSuccess() {
		return fmt.Errorf("emr service reported failure for resource %s", localID)
	}

	return nil
}

func (c *emrClientImpl) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
