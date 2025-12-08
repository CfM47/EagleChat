package idmanagerpool

import (
	"context"
	"sync"

	"eaglechat/apps/client/internal/domain/entities"
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/ezlog"
)

// QueryUsers implements services.IDManagerPool.
func (p *idManagerPoolImpl) QueryUsers(ctx context.Context, IDs []entities.UserID, omitDisconnected bool) (map[entities.UserID]middleware_entities.UserData, error) {
	idManagers := p.repository.GetAll()
	if len(idManagers) == 0 {
		return make(map[entities.UserID]middleware_entities.UserData), nil
	}

	queryCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultsChan := make(chan map[entities.UserID]middleware_entities.UserData, len(idManagers))
	var wg sync.WaitGroup

	for _, data := range idManagers {
		wg.Add(1)
		go func(managerData middleware_entities.IDManagerData) {
			defer wg.Done()

			if queryCtx.Err() != nil {
				return
			}

			connection, err := p.connector.Connect(queryCtx, p.ownProfile, managerData)
			if err != nil {
				if err != context.Canceled {
					ezlog.Log(queryCtx).Warnf("Error connecting to id manager at %s: %v", managerData.IP, err)
				}
				return
			}

			users, err := connection.QueryUsers(queryCtx, IDs, omitDisconnected)
			if err != nil {
				if err != context.Canceled {
					ezlog.Log(queryCtx).Warnf("Error querying users from id manager at %s: %v", managerData.IP, err)
				}
				return
			}

			select {
			case resultsChan <- users:
			case <-queryCtx.Done():
			}
		}(data)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	finalResults := make(map[entities.UserID]middleware_entities.UserData)
	neededIDs := make(map[entities.UserID]struct{}, len(IDs))
	for _, id := range IDs {
		neededIDs[id] = struct{}{}
	}

	for users := range resultsChan {
		for id, userData := range users {
			if _, found := neededIDs[id]; found {
				finalResults[id] = userData
				delete(neededIDs, id)
			}
		}
		if len(neededIDs) == 0 {
			cancel()
		}
	}

	return finalResults, nil
}
