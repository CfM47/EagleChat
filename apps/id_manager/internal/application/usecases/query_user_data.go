package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
	"net"
)

type QueryUserDataUseCase struct {
	repo user.UserRepository
}

func NewQueryUserDataUseCase(repo user.UserRepository) *QueryUserDataUseCase {
	return &QueryUserDataUseCase{repo: repo}
}

type QueryUserRequest struct {
	Ids              []string `json:"Ids"`
	OmitDisconnected bool     `json:"OmitDisconnected"`
}

type UserData struct {
	Username  string  `json:"username"`
	PublicKey string  `json:"public_key"`
	IP        *net.IP `json:"ip,omitempty"`
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

		data := &UserData{
			Username:  user.Username,
			PublicKey: user.PublicKeyPEM,
			IP:        user.IP,
		}

		result[id] = data
	}

	return result, nil
}
