package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/pendingmessage"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
	"net"
)

type QueryPendingMessagesUseCase struct {
	pendingMessageRepo pendingmessage.PendingMessageRepository
	userRepo           user.UserRepository
}

func NewQueryPendingMessagesUseCase(pmr pendingmessage.PendingMessageRepository, ur user.UserRepository) *QueryPendingMessagesUseCase {
	return &QueryPendingMessagesUseCase{pendingMessageRepo: pmr, userRepo: ur}
}

type QueryPendingMessagesRequest struct {
	TargetID   *string `json:"target_id"`
	GetCachers bool    `json:"get_cachers"`
	QuerierID  string  `json:"-"`
	IP         net.IP  `json:"-"`
}

type PendingMessageTarget struct {
	TargetID   string   `json:"target_id"`
	MessageID  string   `json:"message_id"`
	CachersIDs []string `json:"cachers_ids,omitempty"`
}

type QueryPendingMessagesResponse struct {
	MessageTargets []PendingMessageTarget `json:"message_targets"`
}

func (uc *QueryPendingMessagesUseCase) Execute(ctx context.Context, req *QueryPendingMessagesRequest) (*QueryPendingMessagesResponse, error) {
	if req.IP != nil && req.QuerierID != "" {
		_ = uc.userRepo.UpdateIP(req.QuerierID, req.IP)
	}

	var messages []*entities.PendingMessage
	var err error

	if req.TargetID != nil {
		messages, err = uc.pendingMessageRepo.FindByTargetID(*req.TargetID)
	} else {
		messages, err = uc.pendingMessageRepo.FindAll()
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
