package utils

import (
	"fmt"
	"time"
)

// TimeYMDHIS get current time
// return format yearmouthday hour:minute:second
func TimeYMDHIS() string {
	return time.Now().Format("20060102 15:04:05")
}

// TimeYMDH get current time
// return format yearmouthday-hour
func TimeYMDH() string {
	return time.Now().Format("20060102-15")
}

// TimeYMD get current time
// return format yearmouthday
func TimeYMD() string {
	return time.Now().Format("20060102")
}

// Trace calc the request cost ms
func Trace(msg string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("[Trace] %s cast (%s)\n", msg, time.Since(start))
	}
}
