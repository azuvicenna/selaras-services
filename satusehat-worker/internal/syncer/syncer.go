package syncer

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/client"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/fhir"
	"github.com/azuvicenna/selaras-services/satusehat-worker/internal/metrics"
)

type Syncer interface {
	Process(ctx context.Context, task domain.SyncTask) (*domain.SyncResult, error)
}

type syncerImpl struct {
	satusehatClient client.SatusehatClient
	emrClient       client.EMRClient
	orgID           string
	logger          *slog.Logger
}

func NewSyncer(satusehatClient client.SatusehatClient, emrClient client.EMRClient, orgID string, logger *slog.Logger) Syncer {
	if logger == nil {
		logger = slog.Default()
	}
	return &syncerImpl{
		satusehatClient: satusehatClient,
		emrClient:       emrClient,
		orgID:           orgID,
		logger:          logger,
	}
}

func (s *syncerImpl) Process(ctx context.Context, task domain.SyncTask) (*domain.SyncResult, error) {
	start := time.Now()
	resType := strings.TrimSpace(task.ResourceType)

	s.logger.Info("Starting SATUSEHAT synchronization",
		slog.String("task_id", task.TaskID),
		slog.String("resource_type", resType),
		slog.String("resource_id", task.ResourceID),
		slog.String("action", string(task.Action)),
	)

	result := &domain.SyncResult{
		TaskID:       task.TaskID,
		ResourceType: resType,
		ResourceID:   task.ResourceID,
		Timestamp:    time.Now(),
	}

	if task.ResourceID == "" {
		err := domain.ErrMissingResourceID
		result.ErrorMessage = err.Error()
		metrics.RecordSync(resType, "error", time.Since(start).Seconds())
		return result, err
	}

	protoType, err := domain.ToProtoResourceType(resType)
	if err != nil {
		result.ErrorMessage = err.Error()
		metrics.RecordSync(resType, "error", time.Since(start).Seconds())
		return result, fmt.Errorf("%w: %s", domain.ErrUnsupportedResource, resType)
	}

	fhirType, err := domain.ToFHIRResourceType(resType)
	if err != nil {
		result.ErrorMessage = err.Error()
		metrics.RecordSync(resType, "error", time.Since(start).Seconds())
		return result, err
	}

	// 1. Prepare FHIR Payload
	payload, err := fhir.PrepareFHIRPayload(s.orgID, resType, task.Payload)
	if err != nil {
		errMsg := fmt.Sprintf("prepare FHIR payload: %v", err)
		result.ErrorMessage = errMsg
		metrics.RecordSync(resType, "error", time.Since(start).Seconds())
		s.logger.Error("Failed to prepare FHIR payload", slog.String("error", errMsg))
		return result, fmt.Errorf("%w: %v", domain.ErrInvalidTask, err)
	}

	// 2. Transmit to SATUSEHAT API
	var satusehatID string
	if task.Action == domain.ActionUpdate && task.SatusehatID != "" {
		satusehatID, err = s.satusehatClient.PutResource(ctx, fhirType, task.SatusehatID, payload)
		if err != nil {
			errMsg := fmt.Sprintf("put to satusehat: %v", err)
			result.ErrorMessage = errMsg
			metrics.RecordSync(resType, "error", time.Since(start).Seconds())
			s.logger.Error("Failed to update SATUSEHAT resource", slog.String("error", errMsg))
			return result, fmt.Errorf("%w: %v", domain.ErrSyncFailed, err)
		}
	} else {
		satusehatID, err = s.satusehatClient.PostResource(ctx, fhirType, payload)
		if err != nil {
			errMsg := fmt.Sprintf("post to satusehat: %v", err)
			result.ErrorMessage = errMsg
			metrics.RecordSync(resType, "error", time.Since(start).Seconds())
			s.logger.Error("Failed to create SATUSEHAT resource", slog.String("error", errMsg))
			return result, fmt.Errorf("%w: %v", domain.ErrSyncFailed, err)
		}
	}

	result.SatusehatID = satusehatID

	// 3. Callback EMR Service gRPC to persist the returned SATUSEHAT ID
	if s.emrClient != nil {
		if err := s.emrClient.SyncSatusehatID(ctx, protoType, task.ResourceID, satusehatID); err != nil {
			errMsg := fmt.Sprintf("emr sync callback failed: %v", err)
			result.ErrorMessage = errMsg
			metrics.RecordSync(resType, "error", time.Since(start).Seconds())
			s.logger.Error("Failed to sync SATUSEHAT ID back to EMR service",
				slog.String("satusehat_id", satusehatID),
				slog.String("error", errMsg),
			)
			return result, fmt.Errorf("%w: %v", domain.ErrEMRSyncFailed, err)
		}
	}

	duration := time.Since(start).Seconds()
	metrics.RecordSync(resType, "success", duration)
	result.Success = true

	s.logger.Info("Successfully synchronized with SATUSEHAT and updated EMR database",
		slog.String("task_id", task.TaskID),
		slog.String("resource_id", task.ResourceID),
		slog.String("satusehat_id", satusehatID),
		slog.Float64("duration_sec", duration),
	)

	return result, nil
}
