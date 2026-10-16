package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	SyncTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "satusehat_sync_total",
			Help: "Total number of SATUSEHAT synchronization attempts",
		},
		[]string{"resource_type", "status"},
	)

	SyncDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "satusehat_sync_duration_seconds",
			Help:    "Duration of SATUSEHAT synchronization in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"resource_type"},
	)

	QueueMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "satusehat_queue_messages_total",
			Help: "Total number of messages processed from the RabbitMQ queue",
		},
		[]string{"status"},
	)

	WorkerUp = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "satusehat_worker_up",
			Help: "Indicates whether the SATUSEHAT worker is running (1 = up, 0 = down)",
		},
	)
)

func RecordSync(resourceType, status string, duration float64) {
	SyncTotal.WithLabelValues(resourceType, status).Inc()
	SyncDuration.WithLabelValues(resourceType).Observe(duration)
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
