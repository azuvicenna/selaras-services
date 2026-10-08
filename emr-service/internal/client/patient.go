package client

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/azuvicenna/selaras-services/emr-service/internal/domain"
)

const (
	getPatientByIDMethod  = "/patient.v1.PatientService/GetPatientByID"
	patientStatusFieldNum = 27
	patientStatusInactive = 2
	patientStatusDeceased = 3
)

var _ domain.PatientVerifier = (*PatientClient)(nil)

type PatientClient struct {
	conn *grpc.ClientConn
}

func NewPatientClient(addr string) (*PatientClient, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return nil, nil
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to patient-service at %s: %w", addr, err)
	}
	return &PatientClient{conn: conn}, nil
}

func (c *PatientClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *PatientClient) VerifyActivePatient(ctx context.Context, patientID string) error {
	if c == nil || c.conn == nil {
		return nil
	}

	req := &wrapperspb.StringValue{Value: patientID}
	resp := &emptypb.Empty{}

	if err := c.conn.Invoke(ctx, getPatientByIDMethod, req, resp); err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("%w: patient %s not found", domain.ErrNotFound, patientID)
		}
		return fmt.Errorf("patient verification failed: %w", err)
	}

	patientStatus := extractPatientStatus(resp.ProtoReflect().GetUnknown())
	switch patientStatus {
	case patientStatusInactive:
		return fmt.Errorf("%w: patient %s is inactive", domain.ErrInvalidState, patientID)
	case patientStatusDeceased:
		return fmt.Errorf("%w: patient %s is deceased", domain.ErrInvalidState, patientID)
	}

	return nil
}

func extractPatientStatus(raw []byte) uint64 {
	for len(raw) > 0 {
		num, typ, n := protowire.ConsumeTag(raw)
		if n < 0 {
			return 0
		}
		raw = raw[n:]

		if num == 1 && typ == protowire.BytesType {
			patientBytes, m := protowire.ConsumeBytes(raw)
			if m < 0 {
				return 0
			}
			return parsePatientMessageStatus(patientBytes)
		}

		skip := protowire.ConsumeFieldValue(num, typ, raw)
		if skip < 0 {
			return 0
		}
		raw = raw[skip:]
	}
	return 0
}

func parsePatientMessageStatus(b []byte) uint64 {
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return 0
		}
		b = b[n:]

		if num == patientStatusFieldNum && typ == protowire.VarintType {
			v, m := protowire.ConsumeVarint(b)
			if m < 0 {
				return 0
			}
			return v
		}

		skip := protowire.ConsumeFieldValue(num, typ, b)
		if skip < 0 {
			return 0
		}
		b = b[skip:]
	}
	return 0
}
