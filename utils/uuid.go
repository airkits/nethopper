package utils

import (
	uuid "github.com/satori/go.uuid"
	"github.com/yitter/idgenerator-go/idgen"
)

// GenUUID create uuid based on random numbers
func GenUUID() string {
	// or error handling
	u := uuid.NewV4()
	return u.String()
}

// InitUID uidGenerater workerID between 1-255
func InitUID(workID uint16) {
	options := idgen.NewIdGeneratorOptions(workID)
	//	options.SeqBitLength = 10
	idgen.SetIdGenerator(options)
}

// GenUID
func GenUID() uint64 {
	return idgen.NextId()
}
