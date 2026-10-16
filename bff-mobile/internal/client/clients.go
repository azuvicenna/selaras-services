package client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	admissionv1 "github.com/azuvicenna/selaras-services/admission-service/gen/admission/v1"
	authv1 "github.com/azuvicenna/selaras-services/auth-service/gen/auth/v1"
	bedv1 "github.com/azuvicenna/selaras-services/bed-service/gen/bed/v1"
	billingv1 "github.com/azuvicenna/selaras-services/billing-service/gen/billing/v1"
	documentv1 "github.com/azuvicenna/selaras-services/document-service/gen/document/v1"
	emrv1 "github.com/azuvicenna/selaras-services/emr-service/gen/emr/v1"
	laboratoryv1 "github.com/azuvicenna/selaras-services/laboratory-service/gen/laboratory/v1"
	masterv1 "github.com/azuvicenna/selaras-services/master-service/gen/master/v1"
	patientv1 "github.com/azuvicenna/selaras-services/patient-service/gen/patient/v1"
	pharmacyv1 "github.com/azuvicenna/selaras-services/pharmacy-service/gen/pharmacy/v1"
	practitionerv1 "github.com/azuvicenna/selaras-services/practitioner-service/gen/practitioner/v1"
	radiologyv1 "github.com/azuvicenna/selaras-services/radiology-service/gen/radiology/v1"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/config"
)

type ClientManager struct {
	conns []*grpc.ClientConn

	Auth         authv1.AuthServiceClient
	Patient      patientv1.PatientServiceClient
	Practitioner practitionerv1.PractitionerServiceClient
	Admission    admissionv1.AdmissionServiceClient
	Bed          bedv1.BedServiceClient
	Master       masterv1.MasterServiceClient
	EMR          emrv1.EmrServiceClient
	Billing      billingv1.BillingServiceClient
	Pharmacy     pharmacyv1.PharmacyServiceClient
	Laboratory   laboratoryv1.LaboratoryServiceClient
	Radiology    radiologyv1.RadiologyServiceClient
	Document     documentv1.DocumentServiceClient
}

func NewClientManager(cfg *config.Config) (*ClientManager, error) {
	cm := &ClientManager{}

	dial := func(name, addr string) (*grpc.ClientConn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = ctx

		conn, err := grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(25*1024*1024),
				grpc.MaxCallSendMsgSize(25*1024*1024),
			),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to dial %s at %s: %w", name, addr, err)
		}
		cm.conns = append(cm.conns, conn)
		return conn, nil
	}

	var conn *grpc.ClientConn
	var err error

	if conn, err = dial("auth-service", cfg.AuthServiceAddr); err == nil {
		cm.Auth = authv1.NewAuthServiceClient(conn)
	}
	if conn, err = dial("patient-service", cfg.PatientServiceAddr); err == nil {
		cm.Patient = patientv1.NewPatientServiceClient(conn)
	}
	if conn, err = dial("practitioner-service", cfg.PractitionerServiceAddr); err == nil {
		cm.Practitioner = practitionerv1.NewPractitionerServiceClient(conn)
	}
	if conn, err = dial("admission-service", cfg.AdmissionServiceAddr); err == nil {
		cm.Admission = admissionv1.NewAdmissionServiceClient(conn)
	}
	if conn, err = dial("bed-service", cfg.BedServiceAddr); err == nil {
		cm.Bed = bedv1.NewBedServiceClient(conn)
	}
	if conn, err = dial("master-service", cfg.MasterServiceAddr); err == nil {
		cm.Master = masterv1.NewMasterServiceClient(conn)
	}
	if conn, err = dial("emr-service", cfg.EmrServiceAddr); err == nil {
		cm.EMR = emrv1.NewEmrServiceClient(conn)
	}
	if conn, err = dial("billing-service", cfg.BillingServiceAddr); err == nil {
		cm.Billing = billingv1.NewBillingServiceClient(conn)
	}
	if conn, err = dial("pharmacy-service", cfg.PharmacyServiceAddr); err == nil {
		cm.Pharmacy = pharmacyv1.NewPharmacyServiceClient(conn)
	}
	if conn, err = dial("laboratory-service", cfg.LaboratoryServiceAddr); err == nil {
		cm.Laboratory = laboratoryv1.NewLaboratoryServiceClient(conn)
	}
	if conn, err = dial("radiology-service", cfg.RadiologyServiceAddr); err == nil {
		cm.Radiology = radiologyv1.NewRadiologyServiceClient(conn)
	}
	if conn, err = dial("document-service", cfg.DocumentServiceAddr); err == nil {
		cm.Document = documentv1.NewDocumentServiceClient(conn)
	}

	return cm, nil
}

func (cm *ClientManager) Close() {
	for _, conn := range cm.conns {
		if conn != nil {
			if err := conn.Close(); err != nil {
				slog.Warn("failed to close gRPC client connection", slog.String("error", err.Error()))
			}
		}
	}
}
