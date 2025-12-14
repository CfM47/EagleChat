package idmanagerpool

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/ezlog"
)

func (p *idManagerPoolImpl) Now(ctx context.Context) (time.Time, error) {
	ezlog.Log(ctx).Info("Querying ID manager pool for current time")

	idManagers := p.repository.GetAll()
	if len(idManagers) == 0 {
		msg := "No ID managers found in the pool"
		ezlog.Log(ctx).Warn(msg)
		return time.Time{}, errors.New(strings.ToLower(msg))
	}

	queryCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultChan := make(chan time.Time, 1)
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
				// Connector logs errors
				return
			}

			t, err := connection.Time(queryCtx)
			if err != nil {
				// The connection method logs errors
				return
			}

			select {
			case resultChan <- t:
				cancel() // First one to respond wins
			case <-queryCtx.Done():
				return
			}
		}(data)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	now, ok := <-resultChan
	if !ok {
		if ctx.Err() != nil {
			ezlog.Log(ctx).Errorf(ctx.Err().Error())
			return time.Time{}, ctx.Err()
		}
		msg := "All ID Managers failed to respond"
		ezlog.Log(ctx).Error(msg)
		return time.Time{}, errors.New(strings.ToLower(msg))
	}
	return now, nil
}

