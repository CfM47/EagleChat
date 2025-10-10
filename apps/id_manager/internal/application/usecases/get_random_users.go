package usecases

import (
	"math/rand"

	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
)

type GetRandomUsersUseCase struct {
	repo user.UserRepository
}

func NewGetRandomUsersUseCase(repo user.UserRepository) *GetRandomUsersUseCase {
	return &GetRandomUsersUseCase{repo: repo}
}

type GetRandomUsersRequest struct {
	Amount int
}

type GetRandomUsersResponse struct {
	Users []UserData `json:"users"`
}

func (uc *GetRandomUsersUseCase) Execute(req *GetRandomUsersRequest) (*GetRandomUsersResponse, error) {
	allUsers, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var connectedUsers []*entities.User
	for _, u := range allUsers {
		if u.IP != nil {
			connectedUsers = append(connectedUsers, u)
		}
	}

	amount := req.Amount
	if len(connectedUsers) <= amount {
		users := make([]UserData, len(connectedUsers))
		for i, u := range connectedUsers {
			users[i] = UserData{
				Username:  u.Username,
				PublicKey: u.PublicKeyPEM,
				IP:        u.IP,
			}
		}
		return &GetRandomUsersResponse{Users: users}, nil
	}

	rand.Shuffle(len(connectedUsers), func(i, j int) {
		connectedUsers[i], connectedUsers[j] = connectedUsers[j], connectedUsers[i]
	})

	randomUsers := connectedUsers[:amount]

	users := make([]UserData, len(randomUsers))
	for i, u := range randomUsers {
		users[i] = UserData{
			Username:  u.Username,
			PublicKey: u.PublicKeyPEM,
			IP:        u.IP,
		}
	}

	return &GetRandomUsersResponse{Users: users}, nil
}
