package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/application/ports"
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/pendingmessage"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
	"eaglechat/common/ezlog"
	"net"
	"slices"
)

type AddPendingMessagesUseCase struct {
	pendingMessageRepo pendingmessage.PendingMessageRepository
	userRepo           user.UserRepository
	notifier           ports.Notifier

	// Logger context
	logCtx context.Context
}

func NewAddPendingMessagesUseCase(
	pmr pendingmessage.PendingMessageRepository,
	ur user.UserRepository,
	notifier ports.Notifier,
) *AddPendingMessagesUseCase {
	return &AddPendingMessagesUseCase{pendingMessageRepo: pmr, userRepo: ur, notifier: notifier, logCtx: ezlog.NewLoggerContext("pending messages usecase")}
}

type MessageTarget struct {
	TargetID  string `json:"target_id"`
	MessageID string `json:"message_id"`
}

type AddPendingMessagesRequest struct {
	MessageTargets []MessageTarget `json:"message_targets"`
	CacherID       string          `json:"cacher_id"`
	IP             net.IP          `json:"-"`
}

func (uc *AddPendingMessagesUseCase) Execute(ctx context.Context, req *AddPendingMessagesRequest) error {
	if req.IP != nil {
		_ = uc.userRepo.UpdateIP(req.CacherID, req.IP)
	}

	dataChanged := false
	for _, mt := range req.MessageTargets {
		pm, err := uc.pendingMessageRepo.FindByID(mt.MessageID, mt.TargetID)
		if err != nil {
			if err == pendingmessage.ErrPendingMessageNotFound {
				newPm := entities.NewPendingMessage(mt.MessageID, mt.TargetID, []string{req.CacherID})
				if err := uc.pendingMessageRepo.Save(newPm); err != nil {
					return err
				}
				dataChanged = true
				continue
			}
			return err
		}

		cacherExists := slices.Contains(pm.CachersId, req.CacherID)

		if !cacherExists {
			pm.CachersId = append(pm.CachersId, req.CacherID)
			if err := uc.pendingMessageRepo.Save(pm); err != nil {
				return err
			}
			dataChanged = true
		}
	}

	if dataChanged {
		if err := uc.notifier.NotifyPeersOfUpdate(ctx); err != nil {
			// Log the error but don't fail the operation.
			// The periodic sync will eventually catch up.
			ezlog.Log(uc.logCtx).Errorf("AddPendingMessagesUseCase: failed to notify peers of update: %v", err)
		}
	}

	return nil
}
