package freezing

import "time"

// Now returns the current virtual time. It returns a frozen value if a
// freeze is active, otherwise it returns the real time plus any active offset.
func (c *FreezingClock) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.now()
}

// now is the internal, non-locking implementation of the Now() logic.
// It must be called with at least a read lock held.
func (c *FreezingClock) now() time.Time {
	systemNow := time.Now()

	// If we are in a freeze period (unfreezeTime is in the future), return the frozen value.
	// The IsZero check handles the initial state.
	if !c.unfreezeSystemTime.IsZero() && systemNow.Before(c.unfreezeSystemTime) {
		return c.freezeValue
	}

	// If not frozen (or the freeze is over), apply the offset to the real time.
	return systemNow.Add(c.offset).UTC().Round(0)
}
