package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
	"eaglechat/common/ezlog"
)

// SyncData represents the aggregate data exchanged between ID Managers for synchronization.
type SyncData struct {
	Users []*entities.User `json:"users"`
}

// SyncDataUseCase handles the aggregation and merging of ID Manager data.
type SyncDataUseCase struct {
	userRepo user.UserRepository

	// Logger context
	logCtx context.Context
}

// NewSyncDataUseCase creates a new SyncDataUseCase.
func NewSyncDataUseCase(
	userRepo user.UserRepository,
) *SyncDataUseCase {
	return &SyncDataUseCase{
		userRepo: userRepo,
		logCtx:   ezlog.NewLoggerContext("sync data usecase"),
	}
}

// GetAllDataForSync retrieves all relevant data from the current ID Manager for synchronization.
func (uc *SyncDataUseCase) GetAllDataForSync(ctx context.Context) (SyncData, error) {
	users, err := uc.userRepo.FindAll()
	if err != nil {
		return SyncData{}, err
	}

	return SyncData{
		Users: users,
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
					ezlog.Log(uc.logCtx).Errorf("SyncDataUseCase: failed to save new user %s: %v", incomingUser.ID, err)
				}
			} else {
				ezlog.Log(uc.logCtx).Errorf("SyncDataUseCase: failed to find user %s: %v", incomingUser.ID, err)
			}
			continue
		}

		// User exists locally, compare LastSeen to decide on update
		if incomingUser.LastSeen.After(existingUser.LastSeen) {
			// Incoming user data is more recent, update local user.
			if err := uc.userRepo.Update(incomingUser); err != nil {
				ezlog.Log(uc.logCtx).Errorf("SyncDataUseCase: failed to update user %s: %v", incomingUser.ID, err)
			}
		}
		// TODO: If timestamps are equal, we could do a more detailed merge of IPs, but for now,
		// we'll consider it up-to-date or handle conflicts by keeping existing.
		// For simplicity, if LastSeen is equal or older, we do nothing.
		// This implies the local data is authoritative or equally fresh.

	}

	return nil
}
