package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ActivityMonitor = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gin_monitoring",
		Help: "Counting monitoring of gin apps",
	}, []string{"code", "path", "method"})

	EventTransactionMonitMonitor = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "event_transaction_monitoring",
		Help: "Counting monitoring of monit transaction event",
	}, []string{"event"})
)

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func RegisterPrometheus() {
	if err := prometheus.Register(ActivityMonitor); err != nil {
		return
	}

	if err := prometheus.Register(EventTransactionMonitMonitor); err != nil {
		return
	}
}
