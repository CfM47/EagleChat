package persistence

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"sync"
	"time"

	"eaglechat/apps/id_manager/internal/domain"
)

type JSONUserRepository struct {
	filePath string
	mu       sync.Mutex
}

func NewJSONUserRepository(filePath string) *JSONUserRepository {
	return &JSONUserRepository{filePath: filePath}
}

func (r *JSONUserRepository) load() ([]*domain.User, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*domain.User{}, nil
		}
		return nil, err
	}
	var users []*domain.User
	err = json.Unmarshal(data, &users)
	return users, err
}

func (r *JSONUserRepository) save(users []*domain.User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *JSONUserRepository) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	users, err := r.load()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.Username == user.Username {
			users[i] = user
			return r.save(users)
		}
	}

	user.LastSeen = time.Now()
	users = append(users, user)
	return r.save(users)
}

func (r *JSONUserRepository) FindByUsername(username string) (*domain.User, error) {
	users, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *JSONUserRepository) FindAll() ([]*domain.User, error) {
	return r.load()
}

func (r *JSONUserRepository) Delete(username string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	users, err := r.load()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.Username == username {
			users = append(users[:i], users[i+1:]...)
			return r.save(users)
		}
	}
	return errors.New("user not found")
}

func (r *JSONUserRepository) UpdateConnectionStatus(username string, connected bool, ip string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	users, err := r.load()
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Username == username {
			u.Connected = connected
			u.IP = ip
			u.LastSeen = time.Now()
			return r.save(users)
		}
	}
	return errors.New("user not found")
}

func (r *JSONUserRepository) AddPendingMessage(username string, messageID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	users, err := r.load()
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Username == username {
			if slices.Contains(u.PendingMsgs, messageID) {
				return nil
			}
			u.PendingMsgs = append(u.PendingMsgs, messageID)
			return r.save(users)
		}
	}
	return errors.New("user not found")
}

func (r *JSONUserRepository) RemovePendingMessage(username string, messageID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	users, err := r.load()
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.Username == username {
			for i, mid := range u.PendingMsgs {
				if mid == messageID {
					u.PendingMsgs = append(u.PendingMsgs[:i], u.PendingMsgs[i+1:]...)
					return r.save(users)
				}
			}
			return nil
		}
	}
	return errors.New("user not found")
}
