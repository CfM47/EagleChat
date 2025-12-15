package controller

import (
	"context"
	"time"

	"eaglechat/apps/client/internal/domain/services"
	"eaglechat/common/ezlog"
)

// controllerClock is a clock implementation for the controller. It handles cases
// where the middleware has not been constructed yet, and its clock is not available
// this shouldn't happen, so we log it as an error
type controllerClock struct {
	clock services.ClientClock
}

func newControllerClock() controllerClock {
	return controllerClock{}
}

func (c *controllerClock) Now(ctx context.Context) time.Time {
	if c.clock == nil {
		ezlog.Log(ctx).Error("Controller querying clock, but middleware clock has not been set yet")
		return time.Now()
	}

	ezlog.Log(ctx).Info("Controller clock queried")

	return c.clock.Now()
}

func (c *controllerClock) SetClock(ctx context.Context, clock services.ClientClock) {
	ezlog.Log(ctx).Info("Setting controller clock")
	c.clock = clock
}
