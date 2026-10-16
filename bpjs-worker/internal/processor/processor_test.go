package processor_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/domain"
	"github.com/azuvicenna/selaras-services/bpjs-worker/internal/processor"
)

type mockBPJSClient struct {
	pesertaResponse *domain.Peserta
	pesertaErr      error
	rujukanResponse *domain.Rujukan
	rujukanErr      error
	sepResponse     *domain.SEPResponseData
	sepErr          error
	deleteSEPErr    error
	pulangErr       error
	aplicaresErr    error
	antreanErr      error
}

func (m *mockBPJSClient) CariPesertaByNIK(ctx context.Context, nik, tglSEP string) (*domain.Peserta, error) {
	return m.pesertaResponse, m.pesertaErr
}

func (m *mockBPJSClient) CariPesertaByNoKartu(ctx context.Context, noKartu, tglSEP string) (*domain.Peserta, error) {
	return m.pesertaResponse, m.pesertaErr
}

func (m *mockBPJSClient) CariRujukanByNoRujukan(ctx context.Context, noRujukan string) (*domain.Rujukan, error) {
	return m.rujukanResponse, m.rujukanErr
}

func (m *mockBPJSClient) InsertSEP(ctx context.Context, req *domain.InsertSEPRequest) (*domain.SEPResponseData, error) {
	return m.sepResponse, m.sepErr
}

func (m *mockBPJSClient) DeleteSEP(ctx context.Context, noSEP, user string) error {
	return m.deleteSEPErr
}

func (m *mockBPJSClient) UpdateTglPulang(ctx context.Context, req *domain.UpdateTglPulangRequest) error {
	return m.pulangErr
}

func (m *mockBPJSClient) AplicaresUpdateBed(ctx context.Context, req *domain.AplicaresBedRequest) error {
	return m.aplicaresErr
}

func (m *mockBPJSClient) AntreanUpdateWaktu(ctx context.Context, req *domain.UpdateAntreanRequest) error {
	return m.antreanErr
}

func TestProcessor_CheckParticipant(t *testing.T) {
	mockCli := &mockBPJSClient{
		pesertaResponse: &domain.Peserta{
			NoKartu: "0001234567890",
			NIK:     "3171010000000001",
			Nama:    "Dewi Lestari",
		},
	}

	p := processor.NewProcessor(mockCli, nil)

	task := domain.BPJSTask{
		TaskID:  "task-p1",
		Action:  domain.ActionCheckParticipant,
		Payload: json.RawMessage(`{"nik":"3171010000000001"}`),
	}

	res, err := p.Process(context.Background(), task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success")
	}
	if res.ReferenceNo != "0001234567890" {
		t.Errorf("expected reference_no 0001234567890, got %s", res.ReferenceNo)
	}
}

func TestProcessor_CreateSEP(t *testing.T) {
	mockCli := &mockBPJSClient{
		sepResponse: &domain.SEPResponseData{
			NoSep:  "0115R0011026V000999",
			TglSep: "2026-10-16",
		},
	}

	p := processor.NewProcessor(mockCli, nil)

	task := domain.BPJSTask{
		TaskID:  "task-sep-1",
		Action:  domain.ActionCreateSEP,
		Payload: json.RawMessage(`{"request":{"t_sep":{"noKartu":"0001234567890"}}}`),
	}

	res, err := p.Process(context.Background(), task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.Success || res.ReferenceNo != "0115R0011026V000999" {
		t.Errorf("expected noSep 0115R0011026V000999, got %v", res)
	}
}

func TestProcessor_UnsupportedAction(t *testing.T) {
	p := processor.NewProcessor(&mockBPJSClient{}, nil)

	task := domain.BPJSTask{
		TaskID: "task-unknown",
		Action: "UNKNOWN_ACTION",
	}

	_, err := p.Process(context.Background(), task)
	if !errors.Is(err, domain.ErrUnsupportedAction) {
		t.Errorf("expected ErrUnsupportedAction, got %v", err)
	}
}
