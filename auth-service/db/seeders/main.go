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

	authv1 "github.com/azuvicenna/selaras-services/auth-service/gen/auth/v1"
)

func main() {
	addr := flag.String("addr", "localhost:50057", "alamat gRPC auth-service")
	file := flag.String("file", "db/seeders/users.json", "file JSON data seed")
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

	client := authv1.NewAuthServiceClient(conn)

	for _, item := range data.Users {
		res, err := client.CreateUser(ctx, &authv1.CreateUserRequest{
			Username:       item.Username,
			Email:          item.Email,
			Password:       item.Password,
			FullName:       item.FullName,
			Role:           authv1.Role(item.Role),
			PractitionerId: item.PractitionerID,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
				slog.Info("user already exists, skipping", "username", item.Username)
				continue
			}
			return fmt.Errorf("failed to seed user %s: %w", item.Username, err)
		}
		slog.Info("successfully seeded user", "id", res.GetUser().GetId(), "username", res.GetUser().GetUsername())
	}

	slog.Info("seeding completed successfully")
	return nil
}
