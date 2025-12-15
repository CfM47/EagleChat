package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/application/ports"
	"eaglechat/apps/id_manager/internal/domain/repositories"
	"eaglechat/common/clock"
	"eaglechat/common/ezlog"
	"net"
	"time"
)

type AnnounceUseCase struct {
	userRepo           repositories.UserRepository
	notifier           ports.Notifier
	clock              clock.Clock
	expirationDuration time.Duration

	// Logger context
	logCtx context.Context
}

func NewAnnounceUseCase(userRepo repositories.UserRepository, notifier ports.Notifier, clock clock.Clock, expirationDuration time.Duration) *AnnounceUseCase {
	return &AnnounceUseCase{
		userRepo:           userRepo,
		notifier:           notifier,
		clock:              clock,
		expirationDuration: expirationDuration,
		logCtx:             ezlog.NewLoggerContext("Announce UC"),
	}
}

func (uc *AnnounceUseCase) Execute(ctx context.Context, clientID string, ip net.IP) error {
	ezlog.Log(uc.logCtx).Infof("AnnounceUseCase: Executing for clientID: %s, IP: %s", clientID, ip.String())

	// Announce makes the user active by updating their last seen time and IP address.
	// it also expires the user after a certain period of inactivity by making the IP nil.

	// Find user and update IP and last seen time
	err := uc.userRepo.UpdateIP(clientID, ip)
	if err != nil {
		return err
	}

	// Expire every user that has been inactive for a certain period
	expirationThreshold := uc.clock.Now().Add(-uc.expirationDuration)
	users, err := uc.userRepo.FindAll()
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.LastSeen.Before(expirationThreshold) && user.IP != nil {
			ezlog.Log(uc.logCtx).Debugf("User %s ip has expired", user.Username)
			// Expire user by setting IP to nil
			err := uc.userRepo.UpdateIP(user.ID, nil)
			if err != nil {
				return err
			}
		}
	}

	if err := uc.notifier.NotifyPeersOfUpdate(ctx); err != nil {
		// Log the error but don't fail the operation.
		// The periodic sync will eventually catch up.
		ezlog.Log(uc.logCtx).Errorf("Failed to notify peers of update: %v", err)
	}

	return nil
}
