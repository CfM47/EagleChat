package middleware

import (
	"eaglechat/apps/client/internal/domain/entities"
	"fmt"
	"log"
	"sync"
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

const (
	messageSenderInterval = 5 * time.Second
	maxUsersPerTick       = 10
)

func (m *Middleware) messageSender() {
	log.Println("starting message sender...")
	defer log.Println("stopped message sender.")

	for {
		select {
		case <-m.quit:
			return
		case <-m.messageSenderTicker.C:
			m.trySendPendingMessages()
		}
	}
}

func (m *Middleware) trySendPendingMessages() {
	usersToProcess, err := m.getPrioritizedOnlineUsers()
	if err != nil {
		log.Println(err)
		return
	}
	if len(usersToProcess) == 0 {
		return
	}

	// Limit the number of users to process in this tick
	limit := min(len(usersToProcess), maxUsersPerTick)

	log.Printf("found %d online users with pending messages. processing %d this tick.", len(usersToProcess), limit)

	var wg sync.WaitGroup

	for i := range limit {
		userData := usersToProcess[i]
		wg.Add(1)
		go m.sendAllMessagesToUser(&wg, userData)
	}

	wg.Wait()
}

// getPrioritizedOnlineUsers fetches all users that have pending messages and are currently online.
// The returned list is prioritized, with users associated with immune messages appearing first.
func (m *Middleware) getPrioritizedOnlineUsers() ([]middleware_entities.UserData, error) {
	cacheTargets := m.messageCache.GetTargets()
	allCacheTargets := append(cacheTargets.Immune, cacheTargets.NonImmune...)
	if len(allCacheTargets) == 0 {
		return nil, nil
	}

	// Get unique user IDs from all targets
	userIDs := make([]entities.UserID, 0)
	uniqueUserIDs := make(map[entities.UserID]struct{})
	for _, target := range allCacheTargets {
		if _, exists := uniqueUserIDs[target.TargetID]; !exists {
			userIDs = append(userIDs, target.TargetID)
			uniqueUserIDs[target.TargetID] = struct{}{}
		}
	}

	// Find which of those users are online
	onlineUsersData, err := m.getUserData(userIDs, true)
	if err != nil {
		return nil, fmt.Errorf("could not query users: %w", err)
	}

	if len(onlineUsersData) == 0 {
		return nil, nil
	}

	// Create a list of online users to process, maintaining immune-first priority
	usersToProcess := make([]middleware_entities.UserData, 0)
	processedUsers := make(map[entities.UserID]struct{})

	// Add users from immune targets first
	for _, target := range cacheTargets.Immune {
		if userData, isOnline := onlineUsersData[target.TargetID]; isOnline {
			if _, alreadyAdded := processedUsers[target.TargetID]; !alreadyAdded {
				usersToProcess = append(usersToProcess, userData)
				processedUsers[target.TargetID] = struct{}{}
			}
		}
	}
	// Then add users from non-immune targets
	for _, target := range cacheTargets.NonImmune {
		if userData, isOnline := onlineUsersData[target.TargetID]; isOnline {
			if _, alreadyAdded := processedUsers[target.TargetID]; !alreadyAdded {
				usersToProcess = append(usersToProcess, userData)
				processedUsers[target.TargetID] = struct{}{}
			}
		}
	}
	return usersToProcess, nil
}

// sendAllMessagesToUser finds all messages for a given user and attempts to send them.
func (m *Middleware) sendAllMessagesToUser(wg *sync.WaitGroup, userData middleware_entities.UserData) {
	defer wg.Done()
	userID := userData.ID

	// Per user instruction, GetByTargetId expects a UserID and returns all messages for them.
	pendingMessages := m.messageCache.GetByTargetId(userID)
	if len(pendingMessages) == 0 {
		return
	}

	log.Printf("user %s: attempting to send %d pending messages", userID, len(pendingMessages))

	successfullySent := make([]middleware_entities.MessageTarget, 0)
	for _, msg := range pendingMessages {
		err := m.p2pConnPool.Message(userData.IP.String(), fmt.Sprint(m.ownPort), msg.Content)
		if err == nil {
			log.Printf("successfully sent message %s to user %s", msg.Target.ID, msg.Target.TargetID)
			successfullySent = append(successfullySent, msg.Target)
		} else {
			log.Printf("failed to send message %s to user %s: %s", msg.Target.ID, msg.Target.TargetID, err)
		}
	}

	if len(successfullySent) > 0 {
		log.Printf("user %s: deleting %d sent messages from cache", userID, len(successfullySent))
		m.messageCache.DeleteImmune(successfullySent)
	}
}
