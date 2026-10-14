package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	documentv1 "github.com/azuvicenna/selaras-services/document-service/gen/document/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50062", "alamat gRPC document-service")
	file := flag.String("file", "db/seeders/templates.json", "file JSON data seed")
	flag.Parse()

	if err := run(*addr, *file); err != nil {
		slog.Error("seeding failed", "err", err)
		os.Exit(1)
	}
}

func run(addr, file string) error {
	data, err := loadSeedData(file)
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to create gRPC client: %w", err)
	}
	defer func() { _ = conn.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client := documentv1.NewDocumentServiceClient(conn)

	for _, item := range data.Templates {
		res, err := client.CreateDocumentTemplate(ctx, &documentv1.CreateDocumentTemplateRequest{
			Code:         item.Code,
			Name:         item.Name,
			Type:         documentv1.DocumentType(item.Type),
			TypstContent: item.TypstContent,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
				slog.Info("template already exists, skipping", "code", item.Code)
				continue
			}
			return fmt.Errorf("failed to create template %s: %w", item.Code, err)
		}
		slog.Info("template seeded", "id", res.GetTemplate().GetId(), "code", res.GetTemplate().GetCode(), "name", res.GetTemplate().GetName())
	}

	slog.Info("seeding completed successfully", "count", len(data.Templates))
	return nil
}
