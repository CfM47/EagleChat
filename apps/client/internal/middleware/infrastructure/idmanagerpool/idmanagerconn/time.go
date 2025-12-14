package idmanagerconn

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"eaglechat/common/ezlog"
)

// Time implements entities.idManagerConnection.
func (c *idManagerConnectionImpl) Time(ctx context.Context) (time.Time, error) {
	ezlog.Log(ctx).Infof("Querying time from ID manager at %s", c.baseURL)

	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/time", nil)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to create time request", c.baseURL)
		return time.Time{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to query time from ID manager at %s", c.baseURL)
		return time.Time{}, err
	}

	var timeResponse time.Time
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ezlog.Log(ctx).Errorf("Failed to read ID manager at %s time response", c.baseURL)
		return time.Time{}, err
	}
	if err = json.Unmarshal(bodyBytes, &timeResponse); err != nil {
		ezlog.Log(ctx).Errorf("Failed to unmarshall ID manager at %s time response '%s'", c.baseURL, string(bodyBytes))
		return time.Time{}, err
	}

	ezlog.Log(ctx).Infof("Successfully queried time from ID manager at %s", c.baseURL)

	return timeResponse, nil
}
