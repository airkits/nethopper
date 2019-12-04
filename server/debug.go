package server

import (
	"time"
)

// TraceCost calc the api cost time
// usage: defer TraceCose("func")()
func TraceCost(msg string) func() {
	start := time.Now()
	return func() {
		Debug("[TraceCost] %s cost (%s)\n", msg, time.Since(start))
	}
}
