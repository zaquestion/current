package types

import (
	"time"
)

const (
	// citrusleaf epoc: Jan 01 2010 00:00:00 GMT
	CITRUSLEAF_EPOCH = 1262304000
)

// TTL converts an Expiration time from citrusleaf epoc to TTL in seconds.
func TTL(secsFromCitrusLeafEpoc uint32) uint32 {
	return uint32(int64(CITRUSLEAF_EPOCH+secsFromCitrusLeafEpoc) - time.Now().Unix())
}
