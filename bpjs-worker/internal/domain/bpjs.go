package domain

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrInvalidTaskPayload   = errors.New("invalid bpjs task payload")
	ErrUnsupportedAction    = errors.New("unsupported bpjs action")
	ErrBPJSRequestFailed    = errors.New("bpjs web service request failed")
	ErrDecryptionFailed     = errors.New("failed to decrypt bpjs response")
	ErrSEPGenerationFailed  = errors.New("failed to generate sep")
	ErrParticipantNotFound  = errors.New("bpjs participant not found")
)

type ActionType string

const (
	ActionCheckParticipant ActionType = "CHECK_PARTICIPANT"
	ActionCheckRujukan     ActionType = "CHECK_RUJUKAN"
	ActionCreateSEP        ActionType = "CREATE_SEP"
	ActionUpdateSEP        ActionType = "UPDATE_SEP"
	ActionDeleteSEP        ActionType = "DELETE_SEP"
	ActionUpdateTglPulang  ActionType = "UPDATE_TGL_PULANG"
	ActionAplicaresSyncBed ActionType = "APLICARES_SYNC_BED"
	ActionAntreanUpdate    ActionType = "ANTREAN_UPDATE"
)

// Common BPJS Metadata & Response Wrapper
type BPJSMetaData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type BPJSBaseResponse struct {
	MetaData BPJSMetaData    `json:"metaData"`
	Response json.RawMessage `json:"response,omitempty"`
}

// 1. Peserta BPJS Model
type StatusPeserta struct {
	Kode       string `json:"kode"`
	Keterangan string `json:"keterangan"`
}

type JenisPeserta struct {
	Kode       string `json:"kode"`
	Keterangan string `json:"keterangan"`
}

type HakKelas struct {
	Kode       string `json:"kode"`
	Keterangan string `json:"keterangan"`
}

type ProvUmum struct {
	KdProvider string `json:"kdProvider"`
	NmProvider string `json:"nmProvider"`
}

type Peserta struct {
	NoKartu       string        `json:"noKartu"`
	NIK           string        `json:"nik"`
	Nama          string        `json:"nama"`
	Pisa          string        `json:"pisa"`
	Sex           string        `json:"sex"`
	TglLahir      string        `json:"tglLahir"`
	TglCetakKartu string        `json:"tglCetakKartu,omitempty"`
	StatusPeserta StatusPeserta `json:"statusPeserta"`
	JenisPeserta  JenisPeserta  `json:"jenisPeserta"`
	HakKelas      HakKelas      `json:"hakKelas"`
	ProvUmum      ProvUmum      `json:"provUmum"`
}

type PesertaResponse struct {
	Peserta Peserta `json:"peserta"`
}

// 2. Rujukan Model
type Rujukan struct {
	NoKunjungan string  `json:"noKunjungan"`
	TglKunjungan string `json:"tglKunjungan"`
	ProvPerujuk ProvUmum `json:"provPerujuk"`
	Peserta     Peserta `json:"peserta"`
	PoliRujukan struct {
		Kode string `json:"kode"`
		Nama string `json:"nama"`
	} `json:"poliRujukan"`
	Diagnosa struct {
		Kode string `json:"kode"`
		Nama string `json:"nama"`
	} `json:"diagnosa"`
	Keluhan string `json:"keluhan"`
}

type RujukanResponse struct {
	Rujukan Rujukan `json:"rujukan"`
}

// 3. SEP (Surat Eligibilitas Peserta) Models
type InsertSEPRujukan struct {
	AsalRujukan string `json:"asalRujukan"` // "1" Faskes 1, "2" Faskes 2
	TglRujukan  string `json:"tglRujukan"`
	NoRujukan   string `json:"noRujukan"`
	PPKRujukan  string `json:"ppkRujukan"`
}

type InsertSEPPoli struct {
	Tujuan    string `json:"tujuan"`    // Kode Poli
	Eksekutif string `json:"eksekutif"` // "0" Reguler, "1" Eksekutif
}

type InsertSEPCob struct {
	Cob string `json:"cob"` // "0" Tidak, "1" Ya
}

type InsertSEPKatarak struct {
	Katarak string `json:"katarak"` // "0" Tidak, "1" Ya
}

type InsertSEPJaminan struct {
	LakaLantas string `json:"lakaLantas"` // "0" Bukan, "1" Ya
	Penjamin   string `json:"penjamin,omitempty"`
}

type InsertSEPDetail struct {
	NoKartu          string           `json:"noKartu"`
	TglSep           string           `json:"tglSep"`
	PpkPelayanan     string           `json:"ppkPelayanan"`
	JnsPelayanan     string           `json:"jnsPelayanan"` // "1" Rawat Inap, "2" Rawat Jalan
	KlsRawat         string           `json:"klsRawat"`     // "1", "2", "3"
	NoMR             string           `json:"noMR"`
	Rujukan          InsertSEPRujukan `json:"rujukan"`
	Catatan          string           `json:"catatan"`
	DiagAwal         string           `json:"diagAwal"` // ICD-10
	Poli             InsertSEPPoli    `json:"poli"`
	Cob              InsertSEPCob     `json:"cob"`
	Katarak          InsertSEPKatarak `json:"katarak"`
	Jaminan          InsertSEPJaminan `json:"jaminan"`
	TujuanKunjungan  string           `json:"tujuanKunjungan"` // "0" Normal
	FlagProcedure    string           `json:"flagProcedure,omitempty"`
	KdPenunjang      string           `json:"kdPenunjang,omitempty"`
	AssesmentPel     string           `json:"assesmentPel,omitempty"`
	NoSurat          string           `json:"noSurat,omitempty"`
	DpjpLayan        string           `json:"dpjpLayan,omitempty"`
	NoTelp           string           `json:"noTelp"`
	User             string           `json:"user"`
}

type InsertSEPRequest struct {
	Request struct {
		TSEP InsertSEPDetail `json:"t_sep"`
	} `json:"request"`
}

type SEPResponseData struct {
	NoSep        string `json:"noSep"`
	TglSep       string `json:"tglSep"`
	JnsPelayanan string `json:"jnsPelayanan"`
	KlsRawat     string `json:"klsRawat"`
	Diagnosa     string `json:"diagnosa"`
	Poli         string `json:"poli"`
	NoRujukan    string `json:"noRujukan"`
	Catatan      string `json:"catatan"`
	Peserta struct {
		NoKartu    string `json:"noKartu"`
		Nama       string `json:"nama"`
		NoMr       string `json:"noMr"`
		JnsPeserta string `json:"jnsPeserta"`
		HakKelas   string `json:"hakKelas"`
		Kelamin    string `json:"kelamin"`
	} `json:"peserta"`
}

type SEPResponse struct {
	Sep SEPResponseData `json:"sep"`
}

type UpdateTglPulangRequest struct {
	Request struct {
		TSEP struct {
			NoSep            string `json:"noSep"`
			StatusPulang     string `json:"statusPulang"` // 1=Sembuh, 2=Rujuk, 3=APS, 4=Meninggal, 5=Lain-lain
			TglPulang        string `json:"tglPulang"`
			NoSuratMeninggal string `json:"noSuratMeninggal,omitempty"`
			TglMeninggal     string `json:"tglMeninggal,omitempty"`
			NoLPManual       string `json:"noLPManual,omitempty"`
			User             string `json:"user"`
		} `json:"t_sep"`
	} `json:"request"`
}

// 4. Aplicares (Ketersediaan Tempat Tidur)
type AplicaresBedRequest struct {
	KodeKelas         string `json:"kodekelas"`
	KodeRuang         string `json:"koderuang"`
	NamaRuang         string `json:"namaruang"`
	Kapasitas         int    `json:"kapasitas"`
	Tersedia          int    `json:"tersedia"`
	TersediaPria      int    `json:"tersediapria"`
	TersediaWanita    int    `json:"tersediawanita"`
	TersediaPriaWanita int   `json:"tersediapriawanita"`
}

// 5. Antrean RS (Mobile JKN Task ID)
type UpdateAntreanRequest struct {
	KodeBooking string `json:"kodebooking"`
	TaskID      int    `json:"taskid"` // 1: Mulai tunggu admisi, 2: Selesai admisi, 3: Tunggu poli, 4: Pelayanan poli, 5: Selesai poli/resep, 6: Tunggu farmasi, 7: Selesai farmasi
	Waktu       int64  `json:"waktu"`  // Epoch in milliseconds
}

type BatalAntreanRequest struct {
	KodeBooking string `json:"kodebooking"`
	Keterangan  string `json:"keterangan"`
}

// 6. BPJS Task Message for RabbitMQ
type BPJSTask struct {
	TaskID     string          `json:"task_id"`
	Action     ActionType      `json:"action"`
	ResourceID string          `json:"resource_id,omitempty"` // Local ID (e.g. Encounter ID)
	Payload    json.RawMessage `json:"payload"`
	Timestamp  time.Time       `json:"timestamp"`
}

type BPJSResult struct {
	TaskID       string    `json:"task_id"`
	Action       ActionType `json:"action"`
	ResourceID   string    `json:"resource_id,omitempty"`
	ReferenceNo  string    `json:"reference_no,omitempty"` // e.g. No SEP
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"error_message,omitempty"`
	Data         any       `json:"data,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}
