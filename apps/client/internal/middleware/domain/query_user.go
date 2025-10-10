package middleware

import (
	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	usercache "eaglechat/apps/client/internal/middleware/domain/repositories/user_cache"
	"errors"
	"log"
)

// QueryUser implements domain.Middleware.
func (m *Middleware) QueryUser(userID entities.UserID) (entities.User, error) {
	log.SetPrefix("[QueryUser]")
	data, err := m.getUserData(userID)
	if err != nil {
		return entities.User{}, err
	}

	return data.GetUser(), nil
}

func (m *Middleware) getUserData(userID entities.UserID) (middleware_entities.UserData, error) {
	data, err := m.knownUsers.Get(userID)
	if err == nil {
		return data, nil
	}

	if err != usercache.ErrUserNotFound {
		log.Printf("unexpected user cache error: %v", err)
	}

	idManagerConnections, err := m.iDManagerPool.GetAll()
	if err != nil {
		log.Printf("unexpected id manager pool error: %v", err)
		return middleware_entities.UserData{}, err
	}

	for _, conn := range idManagerConnections {
		answ, err := conn.QueryUsers([]entities.UserID{userID})
		if err != nil {
			log.Printf("error querying users: %v", err)
		}

		user, ok := answ[userID]
		if !ok {
			continue
		}

		return user, nil
	}

	log.Printf("user with id '%s' not found in network", userID)

	return middleware_entities.UserData{}, errors.New("user not found in network")
}
