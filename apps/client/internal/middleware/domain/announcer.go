package middleware

import (
	"context"
	"time"

	"eaglechat/common/ezlog"
)

const presenceAnnouncerInterval = time.Second * 10

func (m *Middleware) presenceAnnouncer(ctx context.Context) {
	ezlog.Log(ctx).Info("Starting announcer...")
	defer ezlog.Log(ctx).Info("Announcer stopped")

	announcementTicker := time.NewTicker(presenceAnnouncerInterval)

	for {
		select {
		case <-m.Done():
			return
		case <-announcementTicker.C:
			announcementCtx := ezlog.NewLoggerContext("presence-announcement")
			m.announcePresence(announcementCtx)
		}
	}
}

func (m *Middleware) announcePresence(ctx context.Context) {
	ezlog.Log(ctx).Info("Announcing presence to ID managers...")
	defer ezlog.Log(ctx).Info("Presence announcement done")

	err := m.iDManagerPool.AnnouncePresence(ctx)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to notify ID managers of presence: %v", err)
		return
	}

	ezlog.Log(ctx).Info("Successfully notified ID managers of presence")
}
