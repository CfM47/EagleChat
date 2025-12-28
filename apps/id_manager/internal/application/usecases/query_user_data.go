package usecases

import (
	"context"
	"net"
	"time"

	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
)

type QueryUserDataUseCase struct {
	repo user.UserRepository
}

func NewQueryUserDataUseCase(repo user.UserRepository) *QueryUserDataUseCase {
	return &QueryUserDataUseCase{repo: repo}
}

type QueryUserRequest struct {
	Ids              []string `json:"Ids"`
	OmitDisconnected bool     `json:"omit_disconnected"`
}

type UserData struct {
	Username  string    `json:"username"`
	PublicKey []byte    `json:"public_key"`
	IP        *net.IP   `json:"ip,omitempty"`
	ID        string    `json:"id"`
	LastSeen  time.Time `json:"last_seen"`
}

func NewUserData(user *entities.User) *UserData {
	return &UserData{
		Username:  user.Username,
		PublicKey: user.PublicKeyPEM,
		IP:        user.IP,
		ID:        user.ID,
		LastSeen:  user.LastSeen,
	}
}

type QueryUserResponse map[string]*UserData

func (uc *QueryUserDataUseCase) Execute(ctx context.Context, req *QueryUserRequest) (QueryUserResponse, error) {
	result := make(QueryUserResponse)

	for _, id := range req.Ids {
		user, err := uc.repo.FindByID(id)
		if err != nil {
			continue
		}

		if req.OmitDisconnected && user.IP == nil {
			continue
		}

		data := NewUserData(user)
		result[id] = data
	}

	return result, nil
}
