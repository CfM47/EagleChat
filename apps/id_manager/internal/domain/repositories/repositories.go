package repositories

import (
	"eaglechat/apps/id_manager/internal/domain/repositories/pendingmessage"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
)

type (
	PendingMessageRepository = pendingmessage.PendingMessageRepository
	UserRepository           = user.UserRepository
)

var (
	ErrPendingMessageNotFound = pendingmessage.ErrPendingMessageNotFound
	ErrUserNotFound           = user.ErrUserNotFound
)
