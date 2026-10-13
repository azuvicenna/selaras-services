package middleware

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	grpcRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "practitioner_grpc_requests_total",
			Help: "Total number of gRPC requests completed in practitioner-service",
		},
		[]string{"method", "code"},
	)

	grpcRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "practitioner_grpc_request_duration_seconds",
			Help:    "Histogram of gRPC request durations in seconds in practitioner-service",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "code"},
	)
)

func Metrics(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start).Seconds()

	code := status.Code(err).String()
	grpcRequestsTotal.WithLabelValues(info.FullMethod, code).Inc()
	grpcRequestDuration.WithLabelValues(info.FullMethod, code).Observe(duration)

	return resp, err
}

func RegisterDBMetrics(pool *pgxpool.Pool) {
	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "practitioner_db_pool_total_conns",
			Help: "Total number of database connections in the practitioner-service pool",
		},
		func() float64 {
			return float64(pool.Stat().TotalConns())
		},
	)

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "practitioner_db_pool_idle_conns",
			Help: "Number of idle database connections in the practitioner-service pool",
		},
		func() float64 {
			return float64(pool.Stat().IdleConns())
		},
	)

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "practitioner_db_pool_acquired_conns",
			Help: "Number of currently acquired database connections in practitioner-service",
		},
		func() float64 {
			return float64(pool.Stat().AcquiredConns())
		},
	)
}
