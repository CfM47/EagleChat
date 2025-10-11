package repositories

import (
	"eaglechat/apps/client/internal/middleware/domain/entities"
	"sync"
	"time"
)

const pruneInterval = 10 * time.Second

// expiringIDManagerData wraps IDManagerData with a timestamp to track its freshness.
type expiringIDManagerData struct {
	data     entities.IDManagerData
	lastSeen time.Time
}

// inMemoryIdManagerRepository is a thread-safe, in-memory implementation of the
// IDManagerRepository that automatically prunes stale entries.
type inMemoryIdManagerRepository struct {
	mu       sync.RWMutex
	managers map[string]expiringIDManagerData
}

// NewInMemoryIDManagerRepository creates a new in-memory repository and starts a
// background goroutine to prune stale entries.
func NewInMemoryIDManagerRepository(expirationTime time.Duration) IDManagerRepository {
	repo := &inMemoryIdManagerRepository{
		managers: make(map[string]expiringIDManagerData),
	}

	go repo.startPruning(expirationTime)

	return repo
}

func (r *inMemoryIdManagerRepository) Add(id string, data entities.IDManagerData) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.managers[id] = expiringIDManagerData{
		data:     data,
		lastSeen: time.Now(),
	}
}

func (r *inMemoryIdManagerRepository) Get(id string) (entities.IDManagerData, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	expiringData, ok := r.managers[id]
	if !ok {
		return entities.IDManagerData{}, false
	}
	return expiringData.data, true
}

func (r *inMemoryIdManagerRepository) GetAll() []entities.IDManagerData {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allData []entities.IDManagerData
	for _, expiringData := range r.managers {
		allData = append(allData, expiringData.data)
	}
	return allData
}

// startPruning runs a loop that periodically removes stale entries from the repository.
func (r *inMemoryIdManagerRepository) startPruning(expirationTime time.Duration) {
	ticker := time.NewTicker(pruneInterval)
	defer ticker.Stop()

	for range ticker.C {
		r.prune(expirationTime)
	}
}

// prune removes entries that have not been seen for longer than the expiration time.
func (r *inMemoryIdManagerRepository) prune(expirationTime time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pruneCutoff := time.Now().Add(-expirationTime)

	for id, expiringData := range r.managers {
		if expiringData.lastSeen.Before(pruneCutoff) {
			delete(r.managers, id)
		}
	}
}
