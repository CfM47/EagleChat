package json

import (
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type JSONPendingMessageRepository struct {
	filePath string
	mu       sync.RWMutex
}

var _ repositories.PendingMessageRepository = (*JSONPendingMessageRepository)(nil)

func NewJSONPendingMessageRepository(filepath string) *JSONPendingMessageRepository {
	return &JSONPendingMessageRepository{filePath: filepath}
}

func (r *JSONPendingMessageRepository) load() ([]*entities.PendingMessage, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*entities.PendingMessage{}, nil
		}
		return nil, err
	}
	var collection []*entities.PendingMessage
	err = json.Unmarshal(data, &collection)
	return collection, err
}

func (r *JSONPendingMessageRepository) save(collection []*entities.PendingMessage) error {
	data, err := json.MarshalIndent(collection, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *JSONPendingMessageRepository) Save(item *entities.PendingMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	collection, err := r.load()
	if err != nil {
		return err
	}

	for i, p := range collection {
		if p.MessageId == item.MessageId && p.TargetId == item.TargetId {
			collection[i] = item
			return r.save(collection)
		}
	}
	collection = append(collection, item)
	return r.save(collection)
}

// FindByID implements pendingmessage.PendingMessageRepository.
func (r *JSONPendingMessageRepository) FindByID(message_id string, target_id string) (*entities.PendingMessage, error) {
	r.mu.RLock()
	defer r.mu.RLock()

	collection, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, item := range collection {
		if item.MessageId == message_id && item.TargetId == target_id {
			return item, nil
		}
	}
	return nil, repositories.ErrPendingMessageNotFound
}

func (r *JSONPendingMessageRepository) FindByTargetID(target_id string) ([]*entities.PendingMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	collection, err := r.load()
	if err != nil {
		return nil, err
	}

	var results []*entities.PendingMessage
	for _, item := range collection {
		if item.TargetId == target_id {
			results = append(results, item)
		}
	}
	return results, nil
}

// FindAll implements pendingmessage.PendingMessageRepository.
func (r *JSONPendingMessageRepository) FindAll() ([]*entities.PendingMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.load()
}

// Delete implements pendingmessage.PendingMessageRepository.
func (r *JSONPendingMessageRepository) Delete(message_id string, target_id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	collection, err := r.load()
	if err != nil {
		return err
	}

	for i, item := range collection {
		if item.MessageId == message_id && item.TargetId == target_id {
			collection = append(collection[:i], collection[i+1:]...)
			return r.save(collection)
		}
	}
	return repositories.ErrPendingMessageNotFound
}

// RemoveCacher implements pendingmessage.PendingMessageRepository.
func (r *JSONPendingMessageRepository) RemoveCacher(message_id string, target_id string, cacher_id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	collection, err := r.load()
	if err != nil {
		return err
	}

	for i := range collection {
		item := collection[i]
		if item.MessageId == message_id && item.TargetId == target_id {
			for j := range item.CachersId {
				if item.CachersId[j] == cacher_id {
					item.CachersId = append(item.CachersId[:j], item.CachersId[j+1:]...)
					break
				}
			}
			return r.save(collection)
		}
	}
	return repositories.ErrPendingMessageNotFound
}

// RemoveCachersMany implements pendingmessage.PendingMessageRepository.
func (r *JSONPendingMessageRepository) RemoveCachersMany(message_id string, target_id string, cachers_id []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	collection, err := r.load()
	if err != nil {
		return err
	}

	removeSet := make(map[string]struct{}, len(cachers_id))
	for _, id := range cachers_id {
		removeSet[id] = struct{}{}
	}

	for i := range collection {
		item := collection[i]
		if item.MessageId == message_id && item.TargetId == target_id {
			filtered := item.CachersId[:0]
			for _, id := range item.CachersId {
				if _, toRemove := removeSet[id]; !toRemove {
					filtered = append(filtered, id)
				}
			}
			item.CachersId = filtered
			return r.save(collection)
		}
	}

	return repositories.ErrPendingMessageNotFound
}
