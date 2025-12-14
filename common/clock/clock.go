package clock

import "time"

// Clock represents a source of time that can be queried and adjusted.
// Implementations of this interface provide a controlled and virtualized
// view of time, useful for scenarios requiring custom time progression
// or synchronization.
type Clock interface {
	// Now returns the current time according to this clock.
	// The returned time is guaranteed to be monotonic; it will never
	// return a time earlier than a previously returned value.
	Now() time.Time

	// Sync adjusts the clock to align with the provided targetTime.
	// This method initiates a process to smoothly bring the clock's current
	// time closer to the targetTime, preserving monotonicity.
	// If a prior synchronization process is active, it will be superseded
	// by this new request.
	Sync(targetTime time.Time)
}
