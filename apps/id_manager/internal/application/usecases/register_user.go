package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/application/ports"
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
	"log"
	"net"
)

type RegisterUserUseCase struct {
	repo     user.UserRepository
	notifier ports.Notifier
}

func NewRegisterUserUseCase(repo user.UserRepository, notifier ports.Notifier) *RegisterUserUseCase {
	return &RegisterUserUseCase{repo: repo, notifier: notifier}
}

type RegisterUserRequest struct {
	Username  string `json:"username"`
	PublicKey []byte `json:"public_key"`
	IP        string `json:"-"` // Ignored by JSON binder
}

type RegisterUserResponse struct {
	Id string `json:"id"`
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, req *RegisterUserRequest) (*RegisterUserResponse, error) {
	log.Printf("RegisterUserUseCase: Executing with username: %s", req.Username)
	// The ID is left empty because the repository is responsible for generating it.
	newUser := entities.NewUser("", req.Username, req.PublicKey)

	createdUser, err := uc.repo.Create(newUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	// After creating the user, update their IP address.
	if req.IP != "" {
		parsedIP := net.ParseIP(req.IP)
		log.Printf("Parsed IP: %s", parsedIP.String())
		if parsedIP != nil {
			// Error handling for UpdateIP can be added here if necessary,
			// but for now we can proceed even if it fails.
			err = uc.repo.UpdateIP(createdUser.ID, parsedIP)
			if err != nil {
				log.Printf("An error happened while updating ip")
			}
		}
	}

	if err := uc.notifier.NotifyPeersOfUpdate(ctx); err != nil {
		// Log the error but don't fail the operation.
		// The periodic sync will eventually catch up.
		log.Printf("RegisterUserUseCase: failed to notify peers of update: %v", err)
	}

	// log.Printf("New user registered with Ip: %s", createdUser.IP.String())

	return &RegisterUserResponse{Id: createdUser.ID}, nil
}
