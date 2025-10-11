package usecases

import (
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/pendingmessage"
	"slices"
)

type AddPendingMessagesUseCase struct {
	repo pendingmessage.PendingMessageRepository
}

func NewAddPendingMessagesUseCase(repo pendingmessage.PendingMessageRepository) *AddPendingMessagesUseCase {
	return &AddPendingMessagesUseCase{repo: repo}
}

type MessageTarget struct {
	TargetID  string `json:"target_id"`
	MessageID string `json:"message_id"`
}

type AddPendingMessagesRequest struct {
	MessageTargets []MessageTarget `json:"message_targets"`
	CacherID       string          `json:"cacher_id"`
}

func (uc *AddPendingMessagesUseCase) Execute(req *AddPendingMessagesRequest) error {
	for _, mt := range req.MessageTargets {
		pm, err := uc.repo.FindByID(mt.MessageID, mt.TargetID)
		if err != nil {
			if err == pendingmessage.ErrPendingMessageNotFound {
				newPm := entities.NewPendingMessage(mt.MessageID, mt.TargetID, []string{req.CacherID})
				if err := uc.repo.Save(newPm); err != nil {
					return err
				}
				continue
			}
			return err
		}

		cacherExists := slices.Contains(pm.CachersId, req.CacherID)

		if !cacherExists {
			pm.CachersId = append(pm.CachersId, req.CacherID)
			if err := uc.repo.Save(pm); err != nil {
				return err
			}
		}
	}
	return nil
}
