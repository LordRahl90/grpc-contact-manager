package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// latency use this to manage the latency of the service
	latency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "grpc_contact",
			Name:       "latency_seconds",
			Help:       "The latency value in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"server", "method", "path"},
	)
)

func RegisterPrometheusMetrics() {
	prometheus.MustRegister(latency)
}

// RecordRequestLatency records the latency of the given request per server
func RecordRequestLatency() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		elapsed := time.Since(start).Seconds()
		fmt.Printf("\nElapsed is: %.5f\n", elapsed)
		latency.WithLabelValues(c.Request.Host, c.Request.Method, c.Request.URL.Path).Observe(elapsed)
	}
}
