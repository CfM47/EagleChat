package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/pendingmessage"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
	"log" // Temporarily for internal logging
)

// SyncData represents the aggregate data exchanged between ID Managers for synchronization.
type SyncData struct {
	Users           []*entities.User           `json:"users"`
	PendingMessages []*entities.PendingMessage `json:"pending_messages"`
}

// SyncDataUseCase handles the aggregation and merging of ID Manager data.
type SyncDataUseCase struct {
	userRepo           user.UserRepository
	pendingMessageRepo pendingmessage.PendingMessageRepository
}

// NewSyncDataUseCase creates a new SyncDataUseCase.
func NewSyncDataUseCase(
	userRepo user.UserRepository,
	pendingMessageRepo pendingmessage.PendingMessageRepository,
) *SyncDataUseCase {
	return &SyncDataUseCase{
		userRepo:           userRepo,
		pendingMessageRepo: pendingMessageRepo,
	}
}

// GetAllDataForSync retrieves all relevant data from the current ID Manager for synchronization.
func (uc *SyncDataUseCase) GetAllDataForSync(ctx context.Context) (SyncData, error) {
	users, err := uc.userRepo.FindAll()
	if err != nil {
		return SyncData{}, err
	}

	pendingMessages, err := uc.pendingMessageRepo.FindAll()
	if err != nil {
		return SyncData{}, err
	}

	return SyncData{
		Users:           users,
		PendingMessages: pendingMessages,
	}, nil
}

// MergeData merges incoming synchronization data into the local ID Manager's repositories.
func (uc *SyncDataUseCase) MergeData(ctx context.Context, incomingData SyncData) error {
	// Merge Users
	for _, incomingUser := range incomingData.Users {
		existingUser, err := uc.userRepo.FindByID(incomingUser.ID)
		if err != nil {
			if err == user.ErrUserNotFound {
				// User does not exist locally, create it.
				if err := uc.userRepo.Save(incomingUser); err != nil {
					log.Printf("SyncDataUseCase: failed to save new user %s: %v", incomingUser.ID, err)
				}
			} else {
				log.Printf("SyncDataUseCase: failed to find user %s: %v", incomingUser.ID, err)
			}
			continue
		}

		// User exists locally, compare LastSeen to decide on update
		if incomingUser.LastSeen.After(existingUser.LastSeen) {
			// Incoming user data is more recent, update local user.
			if err := uc.userRepo.Update(incomingUser); err != nil {
				log.Printf("SyncDataUseCase: failed to update user %s: %v", incomingUser.ID, err)
			}
		}
		// TODO: If timestamps are equal, we could do a more detailed merge of IPs, but for now,
		// we'll consider it up-to-date or handle conflicts by keeping existing.
		// For simplicity, if LastSeen is equal or older, we do nothing.
		// This implies the local data is authoritative or equally fresh.

	}

	// Merge PendingMessages
	for _, incomingPM := range incomingData.PendingMessages {
		existingPM, err := uc.pendingMessageRepo.FindByID(incomingPM.MessageId, incomingPM.TargetId)
		if err != nil {
			if err == pendingmessage.ErrPendingMessageNotFound {
				// Pending message does not exist locally, save it.
				if err := uc.pendingMessageRepo.Save(incomingPM); err != nil {
					log.Printf("SyncDataUseCase: failed to save new pending message %s/%s: %v", incomingPM.MessageId, incomingPM.TargetId, err)
				}
			} else {
				log.Printf("SyncDataUseCase: failed to find pending message %s/%s: %v", incomingPM.MessageId, incomingPM.TargetId, err)
			}
			continue
		}

		// Pending message exists locally, merge cachers_id lists.
		// Create a map for quick lookup of existing cachers
		existingCachersMap := make(map[string]struct{})
		for _, cacherID := range existingPM.CachersId {
			existingCachersMap[cacherID] = struct{}{}
		}

		// Add new cachers from incoming data
		updatedCachers := make([]string, len(existingPM.CachersId))
		copy(updatedCachers, existingPM.CachersId) // Start with existing cachers

		for _, incomingCacherID := range incomingPM.CachersId {
			if _, exists := existingCachersMap[incomingCacherID]; !exists {
				updatedCachers = append(updatedCachers, incomingCacherID)
				// Add to map so we don't add duplicates if it appears again in incomingPM.CachersId
				existingCachersMap[incomingCacherID] = struct{}{}
			}
		}

		// Only update if the cachers list has actually changed
		if len(updatedCachers) > len(existingPM.CachersId) {
			existingPM.CachersId = updatedCachers
			if err := uc.pendingMessageRepo.Save(existingPM); err != nil { // Save will overwrite, effectively updating
				log.Printf("SyncDataUseCase: failed to update cachers for pending message %s/%s: %v", incomingPM.MessageId, incomingPM.TargetId, err)
			}
		}
	}
	return nil
}
