package middleware

import (
	"context"
	"errors"
	"maps"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	usercache "eaglechat/apps/client/internal/middleware/domain/repositories/usercache"
	"eaglechat/common/ezlog"
)

// QueryUser implements domain.Middleware.
func (m *Middleware) QueryUser(userID entities.UserID) (entities.User, error) {
	ctx := ezlog.NewLoggerContext("user-query")
	ezlog.Log(ctx).Infof("Querying user %s", userID)

	users, err := m.getUserData(ctx, []entities.UserID{userID}, false)
	if err != nil {
		ezlog.Log(ctx).Errorf("error querying user %s: %v", userID, err)
		return entities.User{}, err
	}

	data, ok := users[userID]
	if !ok {
		ezlog.Log(ctx).Warnf("user %s not found", userID)
		return entities.User{}, errors.New("user not found")
	}

	return data.GetUser(), nil
}

func (m *Middleware) getUserData(ctx context.Context, userIDs []entities.UserID, ensureConnected bool) (map[entities.UserID]middleware_entities.UserData, error) {
	foundUsers := make(map[entities.UserID]middleware_entities.UserData)
	missingUsers := make([]entities.UserID, 0)

	for _, userID := range userIDs {
		data, err := m.knownUsers.Get(userID)
		ezlog.Log(ctx).Debugf("User cache lookup for %s returned: %v, %v", userID, data, err)

		if err == nil {
			if !ensureConnected || data.IP != nil {
				foundUsers[userID] = data
				ezlog.Log(ctx).Debugf("User %s found in cache", userID)
				continue
			}
		} else if err != usercache.ErrUserNotFound {
			ezlog.Log(ctx).Errorf("User cache error for %s: %v", userID, err)
		}
		missingUsers = append(missingUsers, userID)
	}

	if len(missingUsers) == 0 {
		ezlog.Log(ctx).Info("All users found in cache")
		return foundUsers, nil
	}

	queriedUsers, err := m.iDManagerPool.QueryUsers(ctx, missingUsers, ensureConnected)
	if err != nil {
		ezlog.Log(ctx).Errorf("ID manager pool query error: %v", err)
	}

	go func() {
		for _, user := range queriedUsers {
			err := m.knownUsers.Save(user)
			if err != nil {
				ezlog.Log(ctx).Errorf("Failed to store user data in cache: %v", err)
			}
		}
	}()

	maps.Copy(foundUsers, queriedUsers)

	return foundUsers, nil
}
