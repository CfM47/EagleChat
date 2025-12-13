package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/lib"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/apps/client/internal/middleware/domain/repositories/messagecache"
)

const (
	messageSenderInterval           = 10 * time.Second
	maxUsersToSendMessagesToPerTick = 0
	usersToSpreadMessagesTo         = 10
)

func (m *Middleware) messageSender(ctx context.Context) {
	ezlog.Log(ctx).Info("Starting message sender...")

	defer ezlog.Log(ctx).Info("Stopped message sender.")

	for {
		select {
		case <-m.Done():
			return
		//   FIXME: timeouts
		case <-m.messageSenderTicker.C:

			targetedPendingMessageAttemp := ezlog.NewLoggerContext("targeted-send-pending-messages")
			m.sendPendingMessages(targetedPendingMessageAttemp)

			spreadPendingMessagesAttemp := ezlog.NewLoggerContext("spread-pending-messages")
			m.spreadOwnMessages(spreadPendingMessagesAttemp)
		}
	}
}

func (m *Middleware) sendPendingMessages(ctx context.Context) {
	ezlog.Log(ctx).Info("Attempting to send pending messages...")

	cacheTargets := m.messageCache.GetTargets()
	onlineUsers, err := m.getPrioritizedOnlineUsers(ctx, cacheTargets)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to get prioritized online users: %v", err)
		return
	}

	if maxUsersToSendMessagesToPerTick == 0 || len(onlineUsers) < maxUsersToSendMessagesToPerTick {
		ezlog.Log(ctx).Info("No online users with pending messages found.")
		return
	}

	// Limit the number of users to process in this tick
	limit := min(len(onlineUsers), maxUsersToSendMessagesToPerTick)

	ezlog.Log(ctx).Infof("found %d online users with pending messages. processing %d this tick.", len(onlineUsers), limit)

	var wg sync.WaitGroup

	for i := range limit {
		userData := onlineUsers[i]
		wg.Add(1)
		go m.sendAllMessagesToTarget(ctx, &wg, userData)
	}

	wg.Wait()
}

func (m *Middleware) spreadOwnMessages(ctx context.Context) {
	immuneTargets := m.messageCache.GetTargets().Immune
	messages := m.messageCache.GetByTargets(immuneTargets)

	if len(immuneTargets) == 0 {
		return
	}

	usersData, err := m.iDManagerPool.GetRandomConnectedUsers(ctx, usersToSpreadMessagesTo)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to get random connected users %v", err)
		return
	}

	if len(usersData) == 0 {
		ezlog.Log(ctx).Warn("No connected users to spread own messages to")
		return
	}

	ezlog.Log(ctx).Infof("Spreading %d own messages to %d users", len(immuneTargets), len(usersData))

	var wg sync.WaitGroup

	for _, data := range usersData {
		wg.Add(1)
		go func() {
			err := m.clientConnPool.Message(ctx, messages, data)
			if err != nil {
				ezlog.Log(ctx).Errorf("Failed to spread messages to user with ID: '%s', %s", data.ID, data.Name)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

// getPrioritizedOnlineUsers fetches all users that have pending messages and are currently online.
// The returned list is prioritized, with users associated with immune messages appearing first.
func (m *Middleware) getPrioritizedOnlineUsers(ctx context.Context, cacheTargets messagecache.PendingMessageTargetLists) ([]middleware_entities.UserData, error) {
	immuneFirstAllCacheTargets := append(cacheTargets.Immune, cacheTargets.NonImmune...)
	if len(immuneFirstAllCacheTargets) == 0 {
		return nil, nil
	}

	uniqueUserIDs := make(map[entities.UserID]struct{})
	for _, target := range immuneFirstAllCacheTargets {
		uniqueUserIDs[target.TargetID] = struct{}{}
	}

	userIDs := lib.Keys(uniqueUserIDs)

	// Find which of those users are online
	onlineUsersData, err := m.getUserData(ctx, userIDs, true)
	if err != nil {
		msg := "Could not query users"
		ezlog.Log(ctx).Infof(msg+": %v", err)
		return nil, fmt.Errorf(msg+": %w", err)
	}

	if len(onlineUsersData) == 0 {
		ezlog.Log(ctx).Info("Found no online target users")
		return nil, nil
	}

	// Create a list of online users to process, maintaining immune-first priority
	usersToProcess := make(map[entities.UserID]middleware_entities.UserData, 0)

	for _, target := range immuneFirstAllCacheTargets {
		userData, isOnline := onlineUsersData[target.TargetID]
		if !isOnline {
			continue
		}

		usersToProcess[target.TargetID] = userData
	}

	return lib.Values(usersToProcess), nil
}

// sendAllMessagesToTarget finds all messages for a given user and attempts to send them.
func (m *Middleware) sendAllMessagesToTarget(ctx context.Context, wg *sync.WaitGroup, userData middleware_entities.UserData) {
	defer wg.Done()
	userID := userData.ID

	pendingMessages := m.messageCache.GetByTargetId(userID)
	if len(pendingMessages) == 0 {
		return
	}

	ezlog.Log(ctx).Infof("User %s: attempting to send %d pending messages", userID, len(pendingMessages))

	err := m.clientConnPool.Message(ctx, pendingMessages, userData)
	if err != nil {
		ezlog.Log(ctx).Warnf("User %s: failed to send messages: %v", userID, err)
		return
	}

	ezlog.Log(ctx).Infof("User %s: successfully sent %d messages", userID, len(pendingMessages))

	err = m.messageCache.DeleteImmune(getMessageTargets(pendingMessages))
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to delete messages from cache: %v", err)
	}
}

func getMessageTargets(messages []middleware_entities.PendingMessage) []middleware_entities.MessageTarget {
	answ := make([]middleware_entities.MessageTarget, len(messages))

	for i, msg := range messages {
		answ[i] = msg.Target
	}

	return answ
}
