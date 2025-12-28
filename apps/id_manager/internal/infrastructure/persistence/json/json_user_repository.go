package json

import (
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories"
	"eaglechat/common/clock"
	"encoding/json"
	"errors"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type JSONUserRepository struct {
	filePath string
	mu       sync.RWMutex
	clock    clock.Clock
}

type UserRecord struct {
	ID           string    `json:"id"`         // user unique id
	Username     string    `json:"username"`   // user readable alias
	PublicKeyPEM []byte    `json:"public_key"` // PEM codified RSA Public Key
	IP           *string   `json:"ip,omitempty"`
	LastSeen     time.Time `json:"last_seen"` // last moment of connection/disconnection
}

func UserFromRecord(r *UserRecord) (*entities.User, error) {
	var ipPtr *net.IP

	ipPtr = nil
	if r.IP != nil {
		parsed := net.ParseIP(*r.IP)
		if parsed != nil {
			ipPtr = &parsed
		}
	}

	return &entities.User{
		ID:           r.ID,
		Username:     r.Username,
		PublicKeyPEM: r.PublicKeyPEM,
		IP:           ipPtr,
		LastSeen:     r.LastSeen,
	}, nil
}

func RecordFromUser(u *entities.User) *UserRecord {
	var ip *string
	if u.IP != nil {
		s := u.IP.String()
		ip = &s
	}

	return &UserRecord{
		ID:           u.ID,
		Username:     u.Username,
		PublicKeyPEM: u.PublicKeyPEM,
		IP:           ip,
		LastSeen:     u.LastSeen.UTC().Round(0),
	}
}

var _ repositories.UserRepository = (*JSONUserRepository)(nil)

func NewJSONUserRepository(filePath string, clock clock.Clock) *JSONUserRepository {
	return &JSONUserRepository{filePath: filePath, clock: clock}
}

func (r *JSONUserRepository) load() ([]*entities.User, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*entities.User{}, nil
		}
		return nil, err
	}

	var records []*UserRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	users := make([]*entities.User, 0, len(records))
	for _, rec := range records {
		u, err := UserFromRecord(rec)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *JSONUserRepository) save(users []*entities.User) error {
	records := make([]*UserRecord, 0, len(users))
	for _, u := range users {
		records = append(records, RecordFromUser(u))
	}

	data, err := json.MarshalIndent(records, "", "  ")
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
	return nil, repositories.ErrUserNotFound
}

func (r *JSONUserRepository) FindAll() ([]*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users, err := r.load()
	if err != nil {
		return nil, err
	}

	return users, nil
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
	return repositories.ErrUserNotFound
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
			ipCopy := make(net.IP, len(ip))
			copy(ipCopy, ip)

			users[i].IP = &ipCopy
			// Update LastSeen when IP is updated, rounded to zero like in NewUser
			users[i].LastSeen = r.clock.Now().UTC().Round(0)
			return r.save(users)
		}
	}

	return repositories.ErrUserNotFound
}

func (r *JSONUserRepository) Update(user *entities.User) error {
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

	return repositories.ErrUserNotFound
}

func (r *JSONUserRepository) Create(user *entities.User) (*entities.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	users, err := r.load()
	if err != nil {
		return nil, err
	}

	// Generate a new UUID for the user
	user.ID = uuid.New().String()
	user.LastSeen = r.clock.Now().UTC().Round(0)

	users = append(users, user)

	if err := r.save(users); err != nil {
		return nil, err
	}

	return user, nil
}
