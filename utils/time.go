package utils

import (
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
