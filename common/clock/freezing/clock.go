package freezing

import (
	"sync"
	"time"
)

// FreezingClock implements the clock.Clock interface using a hybrid strategy.
// - If the clock is ahead of a target time, it freezes until real time catches up.
// - If the clock is behind a target time, it jumps forward by applying an offset.
type FreezingClock struct {
	mu                 sync.RWMutex
	offset             time.Duration // Applied to real time when not frozen.
	unfreezeSystemTime time.Time     // The real-world time at which the clock unfreezes.
	freezeValue        time.Time     // The virtual time value to return while frozen.
}

// NewFreezingClock creates and returns a new FreezingClock, synchronized
// with the current real time.
func NewFreezingClock() *FreezingClock {
	return &FreezingClock{} // Zero values are correct initial state.
}

