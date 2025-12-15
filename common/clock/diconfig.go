package clock

import "eaglechat/common/clock/freezing"

func NewClock() Clock {
	return freezing.NewFreezingClock()
}
