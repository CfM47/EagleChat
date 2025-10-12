package middleware

import (
	"context"
	"eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/ezlog"
	"sync"
	"time"
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

	var wg sync.WaitGroup

	idManagerConnections, err := m.iDManagerPool.GetAll()
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to get ID manager connections: %v", err)
		return
	}

	pendingMessageTargets := m.messageCache.GetTargets()

	allTargets := make([]entities.MessageTarget, 0, len(pendingMessageTargets.Immune)+len(pendingMessageTargets.NonImmune))
	allTargets = append(allTargets, pendingMessageTargets.Immune...)
	allTargets = append(allTargets, pendingMessageTargets.NonImmune...)

	ezlog.Log(ctx).Infof("Notifying %d ID managers of pending messages for %d targets", len(idManagerConnections), len(allTargets))

	for i, conn := range idManagerConnections {
		wg.Add(1)
		go func() {
			ezlog.Log(ctx).Infof("Notifying ID manager %d of pending messages", i)

			err := conn.NotifyOfPendingMessages(allTargets)
			if err != nil {
				ezlog.Log(ctx).Errorf("Failed to notify ID manager %d of pending messages: %v", i, err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
