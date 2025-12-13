package idmanagerpool

import (
	"time"

	"eaglechat/apps/client/internal/domain/entities"
	managerpool_entities "eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/entities"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool/repositories"
	"eaglechat/common/simplecrypto/rsa"

	"eaglechat/apps/client/internal/middleware/domain/services"
)

const (
	// How often to poll DNS for new ID Managers.
	DNSPOllInterval = 15 * time.Second
	// How long to consider an ID Manager valid without a successful poll.
	ExpirationTime = 30 * time.Second
)

type idManagerPoolImpl struct {
	repository repositories.IDManagerRepository
	ownProfile entities.OwnProfile
	connector  managerpool_entities.IDManagerConnector
	CAPubkey   rsa.PublicKey
	quitChan   chan struct{}
	doneChan   chan struct{}
}

var _ services.IDManagerPool = (*idManagerPoolImpl)(nil)

func (p *idManagerPoolImpl) Close() error {
	close(p.quitChan)
	<-p.doneChan
	return nil
}

func (p *idManagerPoolImpl) Done() <-chan struct{} {
	return p.doneChan
}
