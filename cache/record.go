package cache

import (
	"time"
)

type Record struct {
	Frequency     int
	LastAccessed  time.Time
	InsertionTime time.Time
}
