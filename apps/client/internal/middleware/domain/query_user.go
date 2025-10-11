package middleware

import (
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	usercache "eaglechat/apps/client/internal/middleware/domain/repositories/usercache"
	"errors"
	"log"
)

// QueryUser implements domain.Middleware.
func (m *Middleware) QueryUser(userID entities.UserID) (entities.User, error) {
	log.SetPrefix("[QueryUser] ")

	users, err := m.getUserData([]entities.UserID{userID}, false)
	if err != nil {
		return entities.User{}, err
	}

	data, ok := users[userID]
	if !ok {
		return entities.User{}, errors.New("user not found")
	}

	return data.GetUser(), nil
}

func (m *Middleware) getUserData(userIDs []entities.UserID, ensureConnected bool) (map[entities.UserID]middleware_entities.UserData, error) {
	foundUsers := make(map[entities.UserID]middleware_entities.UserData)
	missingUsers := make([]entities.UserID, 0)

	for _, userID := range userIDs {
		data, err := m.knownUsers.Get(userID)
		if err == nil {
			if !ensureConnected || data.IP != nil {
				foundUsers[userID] = data
				continue
			}
		} else if err != usercache.ErrUserNotFound {
			log.Printf("unexpected user cache error: %v", err)
		}
		missingUsers = append(missingUsers, userID)
	}

	if len(missingUsers) == 0 {
		return foundUsers, nil
	}

	idManagerConnections, err := m.iDManagerPool.GetAll()
	if err != nil {
		log.Printf("unexpected id manager pool error: %v", err)
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
