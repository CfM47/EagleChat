package services

import (
	"context"
	"errors"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
)

var (
	ErrAuthenticationFailed = errors.New("clientconn: peer authentication failed")
	ErrInvalidIP            = errors.New("clientconn: invalid target IP for messaging")
)

type ClientConnPool interface {
	Message(ctx context.Context, messages []middleware_entities.PendingMessage, target middleware_entities.UserData) error
	Receive() <-chan []middleware_entities.PendingMessage
	Done() <-chan struct{}
	Close() error
}

type ClientConnPoolBuilder interface {
	Build(ctx context.Context, ownProfile entities.OwnProfile, listenPort uint16) (ClientConnPool, error)
}
