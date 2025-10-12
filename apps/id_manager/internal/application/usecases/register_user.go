package usecases

import (
	"context"
	"log"
	"net"

	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
)

// Broadcaster defines the interface for broadcasting user updates.
type Broadcaster interface {
	BroadcastUserUpdate(user *entities.User) error
}

type RegisterUserUseCase struct {
	repo        user.UserRepository
	broadcaster Broadcaster
}

func NewRegisterUserUseCase(repo user.UserRepository, broadcaster Broadcaster) *RegisterUserUseCase {
	return &RegisterUserUseCase{repo: repo, broadcaster: broadcaster}
}

type RegisterUserRequest struct {
	Username  string `json:"username"`
	PublicKey []byte `json:"public_key"`
	Port      int    `json:"port"`
}

type RegisterUserResponse struct {
	Id string `json:"id"`
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, req *RegisterUserRequest, ip net.IP) (*RegisterUserResponse, error) {
	// The ID is left empty because the repository is responsible for generating it.
	newUser := entities.NewUser("", req.Username, req.PublicKey, ip)

	createdUser, err := uc.repo.Create(newUser)
	if err != nil {
		return nil, err
	}

	if err := uc.broadcaster.BroadcastUserUpdate(createdUser); err != nil {
		// Log the error but don't fail the main operation.
		log.Printf("error broadcasting user update: %v", err)
	}

	return &RegisterUserResponse{Id: createdUser.ID}, nil
}
