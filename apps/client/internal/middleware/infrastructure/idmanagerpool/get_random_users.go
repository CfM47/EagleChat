package idmanagerpool

import (
	"context"
	"sync"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/ezlog"
	"eaglechat/common/lib"
)

func (p *idManagerPoolImpl) GetRandomConnectedUsers(ctx context.Context, count int) ([]middleware_entities.UserData, error) {
	ezlog.Log(ctx).Infof("Getting %d random connected users from ID Manager Pool", count)

	if count <= 0 {
		ezlog.Log(ctx).Warnf("Requested non-positive count (%d) of random users", count)
		return []middleware_entities.UserData{}, nil
	}

	idManagers := p.repository.GetAll()
	if len(idManagers) == 0 {
		ezlog.Log(ctx).Warn("No ID Managers available in pool to get random users from")
		return []middleware_entities.UserData{}, nil
	}

	queryCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultsChan := make(chan middleware_entities.UserData, count)
	var wg sync.WaitGroup
	wg.Add(len(idManagers))

	for _, data := range idManagers {
		go func(managerData middleware_entities.IDManagerData) {
			defer wg.Done()

			if queryCtx.Err() != nil {
				return
			}

			connection, err := p.connector.Connect(queryCtx, managerData)
			if err != nil {
				// Connector logs errors, so we just exit the goroutine.
				return
			}

			users, err := connection.GetRandomConnectedUsers(queryCtx, count)
			if err != nil {
				// The connection method logs errors, so we just exit.
				return
			}

			for _, user := range users {
				select {
				case resultsChan <- user:
				case <-queryCtx.Done():
					return
				}
			}
		}(data)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	userMap := make(map[entities.UserID]middleware_entities.UserData)
	for user := range resultsChan {
		// This check ensures we don't add more than `count` users,
		// and that we only trigger cancel once.
		if len(userMap) < count {
			userMap[user.ID] = user
			if len(userMap) == count {
				cancel()
			}
		}
	}

	ezlog.Log(ctx).Infof("Retrieved %d unique random connected users from ID Manager Pool", len(userMap))

	return lib.Values(userMap), nil
}
