package idmanagerconn

import (
	"context"
	"fmt"
	"net/http"

	"eaglechat/common/ezlog"
)

// AnnouncePresence implements entities.IDManagerConnection.
func (c *idManagerConnectionImpl) AnnouncePresence(ctx context.Context) error {
	ezlog.Log(ctx).Infof("Announcing presence to ID manager at %s", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/announce", nil)
	if err != nil {
		msg := "Failed to create announce presence request"
		ezlog.Log(ctx).Errorf(msg+": %v", err)
		return fmt.Errorf(msg+": %w", err)
	}

	req.Header.Set("X-Client-ID", string(c.ownID))

	resp, err := c.client.Do(req)
	if err != nil {
		msg := "Failed to announce presence"
		ezlog.Log(ctx).Warnf(msg+": %v", err)
		return fmt.Errorf(msg+": %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		msg := "Received non-OK response when announcing presence: " + resp.Status
		ezlog.Log(ctx).Warnf(msg+": %v", err)
		return fmt.Errorf(msg+": %w", err)
	}

	return nil
}
