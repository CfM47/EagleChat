package controller

import (
	"eaglechat/apps/client/internal/domain/entities"
	"eaglechat/common/ezlog"
	"log"
)

const defaultListenPort = 8081

// connectToMiddleware brings the middleware online with a given user profile.
func (c *Controller) connectToMiddleware(profile entities.OwnProfile) {
	connectorCtx := ezlog.NewLoggerContext("connector")

	ezlog.Log(connectorCtx).Infof("Connecting to middleware for user %s (%s)", profile.User.Name, profile.User.ID)

	connectedMiddleware, msgChan, err := c.connector.Connect(connectorCtx, defaultListenPort, profile.User, profile.PrivateKey)
	if err != nil {
		// This is a critical failure, should probably render an error and quit.
		log.Fatalf("Failed to connect to middleware: %v", err)
	}

	c.connectedMiddleware = connectedMiddleware
	c.messageChannel = msgChan
}
