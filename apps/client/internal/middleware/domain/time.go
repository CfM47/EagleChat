package middleware

import (
	"context"
	"time"

	"eaglechat/common/ezlog"
)

const timeSyncInterval = time.Second * 5

// Now implements services.Middleware.
func (m *Middleware) Now() time.Time {
	return m.clock.Now()
}

func (m *Middleware) timeSynchronizer(ctx context.Context) {
	ezlog.Log(ctx).Info("Starting time synchronizer")
	defer ezlog.Log(ctx).Info("Shutting down time synchronizer")

	syncTicker := time.NewTicker(timeSyncInterval)

	for {
		select {
		case <-syncTicker.C:
			time, err := m.iDManagerPool.Now(ctx)
			if err != nil {
				ezlog.Log(ctx).Errorf("Failed to query id managers for time: %v", err)
			}
			m.clock.Sync(time)
		case <-m.Done():
			return
		}
	}
}
