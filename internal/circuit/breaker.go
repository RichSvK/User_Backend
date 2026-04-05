package circuit

import (
	"log"
	"time"

	"github.com/sony/gobreaker"
)

func NewCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	setting := gobreaker.Settings{
		Name:        name,
		MaxRequests: 5, // Max request in Half-Open state
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second, // Open state time
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio > 0.5
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("[Circuit breaker] %s : %s → %s", name, from.String(), to.String())
		},
	}

	return gobreaker.NewCircuitBreaker(setting)
}
