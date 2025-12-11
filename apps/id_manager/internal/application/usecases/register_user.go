package usecases

import (
	"context"
	"eaglechat/apps/id_manager/internal/application/ports"
	"eaglechat/apps/id_manager/internal/domain/entities"
	"eaglechat/apps/id_manager/internal/domain/repositories/user"
	"eaglechat/common/ezlog"
	"net"
)

type RegisterUserUseCase struct {
	repo     user.UserRepository
	notifier ports.Notifier

	// Logger context
	logCtx context.Context
}

func NewRegisterUserUseCase(repo user.UserRepository, notifier ports.Notifier) *RegisterUserUseCase {
	return &RegisterUserUseCase{repo: repo, notifier: notifier, logCtx: ezlog.NewLoggerContext("register user usecase")}
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
	ezlog.Log(uc.logCtx).Infof("RegisterUserUseCase: Executing with username: %s", req.Username)
	// The ID is left empty because the repository is responsible for generating it.
	newUser := entities.NewUser("", req.Username, req.PublicKey)

	createdUser, err := uc.repo.Create(newUser)
	if err != nil {
		ezlog.Log(uc.logCtx).Errorf("Error creating user: %v", err)
		return nil, err
	}

	// After creating the user, update their IP address.
	if req.IP != "" {
		parsedIP := net.ParseIP(req.IP)
		ezlog.Log(uc.logCtx).Infof("Parsed IP: %s", parsedIP.String())
		if parsedIP != nil {
			// Error handling for UpdateIP can be added here if necessary,
			// but for now we can proceed even if it fails.
			err = uc.repo.UpdateIP(createdUser.ID, parsedIP)
			if err != nil {
				ezlog.Log(uc.logCtx).Error("An error happened while updating ip")
			}
		}
	}

	if err := uc.notifier.NotifyPeersOfUpdate(ctx); err != nil {
		// Log the error but don't fail the operation.
		// The periodic sync will eventually catch up.
		ezlog.Log(uc.logCtx).Errorf("RegisterUserUseCase: failed to notify peers of update: %v", err)
	}

	// log.Printf("New user registered with Ip: %s", createdUser.IP.String())

	return &RegisterUserResponse{Id: createdUser.ID}, nil
}
