package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
)

type RegisterUserUseCase struct {
	repo user.UserRepository
}

func NewRegisterUserUseCase(repo user.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{repo: repo}
}

type RegisterUserRequest struct {
	Username  string `json:"username"`
	PublicKey []byte `json:"public_key"`
}

type RegisterUserResponse struct {
	Id string `json:"id"`
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, req *RegisterUserRequest) (*RegisterUserResponse, error) {
	// The ID is left empty because the repository is responsible for generating it.
	newUser := entities.NewUser("", req.Username, string(req.PublicKey))

	createdUser, err := uc.repo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return &RegisterUserResponse{Id: createdUser.ID}, nil
}
