package idmanagerpool

import (
	"context"
	"errors"
	"strings"

	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/lib/iter"
)

// AnnouncePresence implements services.IDManagerPool.
func (p *idManagerPoolImpl) AnnouncePresence(ctx context.Context) error {
	ezlog.Log(ctx).Info("Announcing presence to ID Managers")

	responses := broadcast(p, ctx, func(ctx context.Context, connection entities.IDManagerConnection) *struct{} {
		err := connection.AnnouncePresence(ctx)
		if err != nil {
			return nil
		}
		return &struct{}{}
	})

	_, _, ok := iter.Find(iter.FromSlice(responses), func(r *struct{}) bool {
		return r != nil
	})

	if !ok {
		msg := "Failed to announce presence to any ID Manager"
		ezlog.Log(ctx).Error(msg)
		return errors.New(strings.ToLower(msg))
	}

	return nil
}
