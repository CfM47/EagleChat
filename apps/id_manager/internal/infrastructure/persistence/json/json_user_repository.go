package json

import (
	"eaglechat/apps/id_manager/internal/domain/entities"
	repository "eaglechat/apps/id_manager/internal/domain/repositories/user"
	"encoding/json"
	"errors"
	"net"
	"os"
	"sync"
)

type JSONUserRepository struct {
	filePath string
	mu       sync.RWMutex
}

var _ repository.UserRepository = (*JSONUserRepository)(nil)

func NewJSONUserRepository(filePath string) *JSONUserRepository {
	return &JSONUserRepository{filePath: filePath}
}

func (r *JSONUserRepository) load() ([]*entities.User, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*entities.User{}, nil
		}
		return nil, err
	}
	var users []*entities.User
	err = json.Unmarshal(data, &users)
	return users, err
}

func (r *JSONUserRepository) save(users []*entities.User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *JSONUserRepository) Save(user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	users, err := r.load()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.ID == user.ID {
			users[i] = user
			return r.save(users)
		}
	}

	users = append(users, user)
	return r.save(users)
}

func (r *JSONUserRepository) FindByID(ID string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.ID == ID {
			return u, nil
		}
	}
	return nil, repository.ErrUserNotFound
}

func (r *JSONUserRepository) FindAll() ([]*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.load()
}

func (r *JSONUserRepository) Delete(ID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	users, err := r.load()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.ID == ID {
			users = append(users[:i], users[i+1:]...)
			return r.save(users)
		}
	}
	return repository.ErrUserNotFound
}

func (r *JSONUserRepository) UpdateIP(ID string, ip net.IP) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	users, err := r.load()
	if err != nil {
		return err
	}

	for i := range users {
		if users[i].ID == ID {
			users[i].IP = &ip
			return r.save(users)
		}
	}

	return repository.ErrUserNotFound
}

