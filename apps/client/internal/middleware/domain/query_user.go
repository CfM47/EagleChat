package middleware

import (
	"context"
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	usercache "eaglechat/apps/client/internal/middleware/domain/repositories/usercache"
	"eaglechat/common/ezlog"
	"errors"
	"log"
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

	idManagerConnections, err := m.iDManagerPool.GetAll()
	if err != nil {
		ezlog.Log(ctx).Errorf("ID manager pool error: %v", err)
		return nil, err
	}

	for _, conn := range idManagerConnections {
		answ, err := conn.QueryUsers(missingUsers, ensureConnected)
		if err != nil {
			log.Printf("error querying users: %v", err)
			continue
		}

		for id, user := range answ {
			if err := m.knownUsers.Save(user); err != nil {
				log.Printf("error caching user data: %v", err)
			}
			foundUsers[id] = user
		}

		// If we found all missing users, we can stop
		if len(foundUsers) == len(userIDs) {
			return foundUsers, nil
		}
	}

	return foundUsers, nil
}
