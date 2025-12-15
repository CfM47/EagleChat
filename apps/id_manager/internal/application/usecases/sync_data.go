package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
	"eaglechat/common/clock"
	"eaglechat/common/ezlog"
	"time"
)

// SyncData represents the aggregate data exchanged between ID Managers for synchronization.
type SyncData struct {
	Users []*entities.User `json:"users"`
}

// SyncDataUseCase handles the aggregation and merging of ID Manager data.
type SyncDataUseCase struct {
	userRepo           user.UserRepository
	clock              clock.Clock
	expirationDuration time.Duration

	// Logger context
	logCtx context.Context
}

// NewSyncDataUseCase creates a new SyncDataUseCase.
func NewSyncDataUseCase(
	userRepo user.UserRepository,
	clock clock.Clock,
	expirationDuration time.Duration,
) *SyncDataUseCase {
	return &SyncDataUseCase{
		userRepo:           userRepo,
		clock:              clock,
		expirationDuration: expirationDuration,
		logCtx:             ezlog.NewLoggerContext("sync data usecase"),
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

		expirationThreshold := uc.clock.Now().Add(-uc.expirationDuration)

		// If timestamps are equal or local is newer, check for ip expiration on the existing user.
		if existingUser.LastSeen.Before(expirationThreshold) && existingUser.IP != nil {
			ezlog.Log(uc.logCtx).Debugf("User %s ip has expired", existingUser.Username)
			// Expire user by setting IP to nil
			err := uc.userRepo.UpdateIP(existingUser.ID, nil)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
