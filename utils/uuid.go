package utils

import (
	uuid "github.com/satori/go.uuid"
	"github.com/zheng-ji/goSnowFlake"
)

//GenUUID create uuid based on random numbers
func GenUUID() string {
	// or error handling
	u := uuid.NewV4()
	return u.String()
}

//GenUID uidGenerater workerID between 1-255
func GenUID(workID int64) (int64, error) {
	iw, err := goSnowFlake.NewIdWorker(workID)
	if err != nil {
		return 0, err
	}
	return iw.NextId()
}
