package freezing

import (
	"time"
)

// Sync adjusts the clock based on the provided targetTime.
// If the clock is ahead of the target, it freezes for the duration of the skew.
// If the clock is behind the target, it jumps forward to the target time.
func (c *FreezingClock) Sync(targetTime time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	systemNow := time.Now()
	currentTime := c.now()

	skew := currentTime.Sub(targetTime)

	if skew > 0 {
		// --- CLOCK IS AHEAD: FREEZE ---
		freezeDuration := skew

		// Pre-calculate and set the offset that will be used AFTER the freeze.
		c.offset = targetTime.Sub(systemNow)

		// Set the freeze state.
		c.freezeValue = currentTime
		c.unfreezeSystemTime = systemNow.Add(freezeDuration)

	} else {
		// --- CLOCK IS BEHIND: JUMP ---
		// Jump forward by setting the offset.
		c.offset = targetTime.Sub(time.Now())
	}
}

