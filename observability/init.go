package observability

import (
	"log"
)

// Init initializes all observability components
// This is the ONLY function that business logic calls
func Init() {
	log.Println("🔍 Initializing observability...")

	// Initialize tracing
	if err := initTracing(); err != nil {
		log.Printf("Tracing init failed: %v", err)
	}

	// Initialize profiling
	if err := initProfiling(); err != nil {
		log.Printf("Profiling init failed: %v", err)
	}

	// Initialize metrics
	initMetrics()

	log.Println("✅ Observability initialized")
}
