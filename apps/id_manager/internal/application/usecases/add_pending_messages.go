package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/pendingmessage"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
	"log"
	"net"
	"slices"
)

type AddPendingMessagesUseCase struct {
	pendingMessageRepo pendingmessage.PendingMessageRepository
	userRepo           user.UserRepository
}

func NewAddPendingMessagesUseCase(pmr pendingmessage.PendingMessageRepository, ur user.UserRepository) *AddPendingMessagesUseCase {
	return &AddPendingMessagesUseCase{pendingMessageRepo: pmr, userRepo: ur}
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

	log.Printf("starting lop")
	for _, mt := range req.MessageTargets {
		pm, err := uc.pendingMessageRepo.FindByID(mt.MessageID, mt.TargetID)
		log.Printf("i was able to find by id")
		if err != nil {
			if err == pendingmessage.ErrPendingMessageNotFound {
				newPm := entities.NewPendingMessage(mt.MessageID, mt.TargetID, []string{req.CacherID})
				if err := uc.pendingMessageRepo.Save(newPm); err != nil {
					log.Printf("i had an error new pending message")
					return err
				}
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
		}
	}
	return nil
}
