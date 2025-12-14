package controller

import (
	"context"
	"log"

	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"
)

const defaultListenPort = 8081

// connectToMiddleware brings the middleware online with a given user profile.
func (c *Controller) connectToMiddleware(ctx context.Context, profile entities.OwnProfile) {
	ezlog.Log(ctx).Infof("Connecting to middleware for user %s (%s)", profile.User.Name, profile.User.ID)

	connectedMiddleware, msgChan, err := c.connector.Connect(ctx, defaultListenPort, profile, c.CAPubkey)
	if err != nil {
		// This is a critical failure, should probably render an error and quit.
		log.Fatalf("Failed to connect to middleware: %v", err)
	}

	c.connectedMiddleware = connectedMiddleware
	c.messageChannel = msgChan
	c.clock.SetClock(ctx, connectedMiddleware)
}
