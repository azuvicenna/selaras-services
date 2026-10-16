package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bpjs_request_total",
			Help: "Total number of BPJS Web Service requests",
		},
		[]string{"action", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bpjs_request_duration_seconds",
			Help:    "Duration of BPJS Web Service requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"action"},
	)

	QueueMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bpjs_queue_messages_total",
			Help: "Total number of BPJS tasks processed from RabbitMQ queue",
		},
		[]string{"status"},
	)

	WorkerUp = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bpjs_worker_up",
			Help: "Indicates whether the BPJS worker is running (1 = up, 0 = down)",
		},
	)
)

func RecordRequest(action, status string, duration float64) {
	RequestTotal.WithLabelValues(action, status).Inc()
	RequestDuration.WithLabelValues(action).Observe(duration)
}

func RecordQueueMessage(status string) {
	QueueMessagesTotal.WithLabelValues(status).Inc()
}

func SetWorkerUp(up bool) {
	if up {
		WorkerUp.Set(1)
	} else {
		WorkerUp.Set(0)
	}
}
