package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/client"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/metrics"
)

type Processor interface {
	Process(ctx context.Context, task domain.BPJSTask) (*domain.BPJSResult, error)
}

type processorImpl struct {
	bpjsClient client.BPJSClient
	logger     *slog.Logger
}

func NewProcessor(bpjsClient client.BPJSClient, logger *slog.Logger) Processor {
	if logger == nil {
		logger = slog.Default()
	}
	return &processorImpl{
		bpjsClient: bpjsClient,
		logger:     logger,
	}
}

func (p *processorImpl) Process(ctx context.Context, task domain.BPJSTask) (*domain.BPJSResult, error) {
	start := time.Now()
	actionStr := string(task.Action)

	p.logger.Info("Processing BPJS task",
		slog.String("task_id", task.TaskID),
		slog.String("action", actionStr),
		slog.String("resource_id", task.ResourceID),
	)

	result := &domain.BPJSResult{
		TaskID:     task.TaskID,
		Action:     task.Action,
		ResourceID: task.ResourceID,
		Timestamp:  time.Now(),
	}

	var err error
	switch task.Action {
	case domain.ActionCheckParticipant:
		var pReq struct {
			NIK     string `json:"nik"`
			NoKartu string `json:"no_kartu"`
			TglSEP  string `json:"tgl_sep"`
		}
		if err = json.Unmarshal(task.Payload, &pReq); err != nil {
			err = fmt.Errorf("%w: %v", domain.ErrInvalidTaskPayload, err)
			break
		}

		var peserta *domain.Peserta
		if pReq.NIK != "" {
			peserta, err = p.bpjsClient.CariPesertaByNIK(ctx, pReq.NIK, pReq.TglSEP)
		} else if pReq.NoKartu != "" {
			peserta, err = p.bpjsClient.CariPesertaByNoKartu(ctx, pReq.NoKartu, pReq.TglSEP)
		} else {
			err = fmt.Errorf("%w: either nik or no_kartu is required", domain.ErrInvalidTaskPayload)
			break
		}

		if err == nil {
			result.Data = peserta
			result.ReferenceNo = peserta.NoKartu
		}

	case domain.ActionCheckRujukan:
		var rReq struct {
			NoRujukan string `json:"no_rujukan"`
		}
		if err = json.Unmarshal(task.Payload, &rReq); err != nil || rReq.NoRujukan == "" {
			err = fmt.Errorf("%w: no_rujukan is required", domain.ErrInvalidTaskPayload)
			break
		}

		var rujukan *domain.Rujukan
		rujukan, err = p.bpjsClient.CariRujukanByNoRujukan(ctx, rReq.NoRujukan)
		if err == nil {
			result.Data = rujukan
			result.ReferenceNo = rujukan.NoKunjungan
		}

	case domain.ActionCreateSEP:
		var sepReq domain.InsertSEPRequest
		if err = json.Unmarshal(task.Payload, &sepReq); err != nil {
			err = fmt.Errorf("%w: %v", domain.ErrInvalidTaskPayload, err)
			break
		}

		var sep *domain.SEPResponseData
		sep, err = p.bpjsClient.InsertSEP(ctx, &sepReq)
		if err == nil {
			result.Data = sep
			result.ReferenceNo = sep.NoSep
		}

	case domain.ActionDeleteSEP:
		var delReq struct {
			NoSEP string `json:"no_sep"`
			User  string `json:"user"`
		}
		if err = json.Unmarshal(task.Payload, &delReq); err != nil || delReq.NoSEP == "" {
			err = fmt.Errorf("%w: no_sep is required", domain.ErrInvalidTaskPayload)
			break
		}

		err = p.bpjsClient.DeleteSEP(ctx, delReq.NoSEP, delReq.User)
		if err == nil {
			result.ReferenceNo = delReq.NoSEP
		}

	case domain.ActionUpdateTglPulang:
		var pulangReq domain.UpdateTglPulangRequest
		if err = json.Unmarshal(task.Payload, &pulangReq); err != nil {
			err = fmt.Errorf("%w: %v", domain.ErrInvalidTaskPayload, err)
			break
		}

		err = p.bpjsClient.UpdateTglPulang(ctx, &pulangReq)
		if err == nil {
			result.ReferenceNo = pulangReq.Request.TSEP.NoSep
		}

	case domain.ActionAplicaresSyncBed:
		var bedReq domain.AplicaresBedRequest
		if err = json.Unmarshal(task.Payload, &bedReq); err != nil {
			err = fmt.Errorf("%w: %v", domain.ErrInvalidTaskPayload, err)
			break
		}

		err = p.bpjsClient.AplicaresUpdateBed(ctx, &bedReq)

	case domain.ActionAntreanUpdate:
		var antreanReq domain.UpdateAntreanRequest
		if err = json.Unmarshal(task.Payload, &antreanReq); err != nil {
			err = fmt.Errorf("%w: %v", domain.ErrInvalidTaskPayload, err)
			break
		}

		err = p.bpjsClient.AntreanUpdateWaktu(ctx, &antreanReq)
		if err == nil {
			result.ReferenceNo = antreanReq.KodeBooking
		}

	default:
		err = fmt.Errorf("%w: %s", domain.ErrUnsupportedAction, task.Action)
	}

	duration := time.Since(start).Seconds()

	if err != nil {
		result.Success = false
		result.ErrorMessage = err.Error()
		metrics.RecordRequest(actionStr, "error", duration)
		p.logger.Error("BPJS task processing failed",
			slog.String("task_id", task.TaskID),
			slog.String("action", actionStr),
			slog.String("error", err.Error()),
		)
		return result, err
	}

	result.Success = true
	metrics.RecordRequest(actionStr, "success", duration)
	p.logger.Info("BPJS task completed successfully",
		slog.String("task_id", task.TaskID),
		slog.String("action", actionStr),
		slog.String("reference_no", result.ReferenceNo),
		slog.Float64("duration_sec", duration),
	)

	return result, nil
}
