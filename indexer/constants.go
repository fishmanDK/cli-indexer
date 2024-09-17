package indexer

import "time"

const (
	maxRetries     = 5
	retryDelay     = 3000 * time.Millisecond
	parceBlockSleep = 8000 * time.Millisecond
)