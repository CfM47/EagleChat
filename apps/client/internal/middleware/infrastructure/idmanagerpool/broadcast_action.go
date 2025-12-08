package idmanagerpool

import (
	"context"
	"sync"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	managerpool_entities "eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/entities"
	"eaglechat/common/ezlog"
)

type broadcastAction[R any] func(ctx context.Context, connection managerpool_entities.IDManagerConnection) *R

func broadcast[R any](
	p *idManagerPoolImpl,
	ctx context.Context,
	action broadcastAction[R],
) []*R {
	idManagers := p.repository.GetAll()
	if len(idManagers) == 0 {
		return []*R{}
	}

	results := make(chan *R, len(idManagers))
	var wg sync.WaitGroup
	wg.Add(len(idManagers))

	for _, data := range idManagers {
		go singleManagerAction(p, ctx, action, data, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	answ := make([]*R, 0)
	for result := range results {
		if result != nil {
			answ = append(answ, result)
		}
	}

	return answ
}

func singleManagerAction[R any](
	p *idManagerPoolImpl,
	ctx context.Context,
	action broadcastAction[R],
	data middleware_entities.IDManagerData,
	results chan<- *R,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	connection, err := p.connector.Connect(ctx, p.ownProfile, data)
	if err != nil {
		ezlog.Log(ctx).Errorf("Error connecting to id manager at %s: %v", data.IP, err)
		return
	}

	results <- action(ctx, connection)
}
