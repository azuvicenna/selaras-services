package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	NotificationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notification_sent_total",
			Help: "Total number of notification transmissions",
		},
		[]string{"channel", "type", "status"},
	)

	NotificationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "notification_duration_seconds",
			Help:    "Duration of notification transmissions in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"channel"},
	)

	QueueMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notification_queue_messages_total",
			Help: "Total number of notification messages processed from RabbitMQ",
		},
		[]string{"status"},
	)

	WorkerUp = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "notification_worker_up",
			Help: "Indicates whether the notification worker is running (1 = up, 0 = down)",
		},
	)
)

func RecordNotification(channel, notifType, status string, duration float64) {
	NotificationTotal.WithLabelValues(channel, notifType, status).Inc()
	NotificationDuration.WithLabelValues(channel).Observe(duration)
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
