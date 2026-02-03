package observability

import (
	"github.com/grafana/pyroscope-go"
)

func initProfiling() error {
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: GetEnv("OTEL_SERVICE_NAME", "ecommerce-backend"),
		ServerAddress:   GetEnv("PYROSCOPE_ENDPOINT", "http://pyroscope-distributor.monitoring.svc.cluster.local:4040"),
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
		},
	})
	return err
}

func initMetrics() {
	// Metrics are registered in handlers package
	// This function can be extended for custom metric initialization
}
