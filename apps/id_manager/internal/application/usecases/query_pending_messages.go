package usecases

import (
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/pendingmessage"
)

type QueryPendingMessagesUseCase struct {
	repo pendingmessage.PendingMessageRepository
}

func NewQueryPendingMessagesUseCase(repo pendingmessage.PendingMessageRepository) *QueryPendingMessagesUseCase {
	return &QueryPendingMessagesUseCase{repo: repo}
}

type QueryPendingMessagesRequest struct {
	TargetID   *string `json:"target_id"`
	GetCachers bool    `json:"get_cachers"`
}

type PendingMessageTarget struct {
	TargetID   string   `json:"target_id"`
	MessageID  string   `json:"message_id"`
	CachersIDs []string `json:"cachers_ids,omitempty"`
}

type QueryPendingMessagesResponse struct {
	MessageTargets []PendingMessageTarget `json:"message_targets"`
}

func (uc *QueryPendingMessagesUseCase) Execute(req *QueryPendingMessagesRequest) (*QueryPendingMessagesResponse, error) {
	var messages []*entities.PendingMessage
	var err error

	if req.TargetID != nil {
		messages, err = uc.repo.FindByTargetID(*req.TargetID)
	} else {
		messages, err = uc.repo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	var targets []PendingMessageTarget
	for _, msg := range messages {
		target := PendingMessageTarget{
			TargetID:  msg.TargetId,
			MessageID: msg.MessageId,
		}
		if req.GetCachers {
			target.CachersIDs = msg.CachersId
		}
		targets = append(targets, target)
	}

	return &QueryPendingMessagesResponse{MessageTargets: targets}, nil
}
