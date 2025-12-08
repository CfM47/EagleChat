package middleware

import (
	"context"
	"time"

	"eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/ezlog"
)

const presenceAnnouncerInterval = time.Second * 10

func (m *Middleware) presenceAnnouncer(ctx context.Context) {
	ezlog.Log(ctx).Info("Starting announcer...")
	defer ezlog.Log(ctx).Info("Announcer stopped")

	for {
		select {
		case <-m.Done():
			return
		case <-m.announcementTicker.C:
			announcementCtx := ezlog.NewLoggerContext("presence-announcement")
			m.announcePresence(announcementCtx)
		}
	}
}

func (m *Middleware) announcePresence(ctx context.Context) {
	ezlog.Log(ctx).Info("Announcing presence to ID managers...")
	defer ezlog.Log(ctx).Info("Presence announcement done")

	pendingMessageTargets := m.messageCache.GetTargets()

	allTargets := make([]entities.MessageTarget, 0, len(pendingMessageTargets.Immune)+len(pendingMessageTargets.NonImmune))
	allTargets = append(allTargets, pendingMessageTargets.Immune...)
	allTargets = append(allTargets, pendingMessageTargets.NonImmune...)

	err := m.iDManagerPool.NotifyOfPendingMessages(ctx, allTargets)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to notify ID managers of pending messages: %v", err)
		return
	}

	ezlog.Log(ctx).Info("Successfully notified ID managers of pending messages")
}
