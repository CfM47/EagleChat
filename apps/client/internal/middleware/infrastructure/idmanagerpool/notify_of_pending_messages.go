package idmanagerpool

import (
	"context"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	managerpool_entities "eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/entities"
	"eaglechat/common/ezlog"
)

// NotifyOfPendingMessages implements services.IDManagerPool.
func (p *idManagerPoolImpl) NotifyOfPendingMessages(ctx context.Context, targets []middleware_entities.MessageTarget) error {
	broadcast(
		p,
		ctx,
		func(
			ctx context.Context,
			connection managerpool_entities.IDManagerConnection,
		) *struct{} {
			err := connection.NotifyOfPendingMessages(ctx, p.ownProfile.User.ID, targets)
			if err != nil {
				ezlog.Log(ctx).Errorf("Failed to notify ID manager at %s of pending messages: %v", connection.BaseURL(), err)
			}

			return &struct{}{}
		})

	return nil
}
