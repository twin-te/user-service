package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/twin-te/user-service/entity"
	"github.com/twin-te/user-service/port"
)

type UseCase struct {
	r port.IRepository
}

// GetUserByID returns the user specified by id.
// If the corresponding user does not exist, return ErrUserNotFound.
func (uc *UseCase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	u, err := uc.r.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UseCase.GetUserByID %q -> %w", id, err)
	}
	return u, nil
}

// GetOrCreateUser returns the authorized user with the specified authentication.
// If the corresponding user does not exist, create a new user.
func (uc *UseCase) GetOrCreateUser(ctx context.Context, a *entity.Authentication) (*entity.User, error) {
	u, err := uc.r.GetUserByAuthentication(ctx, a)
	if err == nil {
		return u, nil
	}
	if !errors.Is(err, entity.ErrUserNotFound) {
		return nil, fmt.Errorf("UseCase.GetOrCreateUser %+v -> %w", a, err)
	}
	u = &entity.User{ID: uuid.NewString(), Authentications: []*entity.Authentication{a}}
	err = uc.r.CreateUser(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("UseCase.GetOrCreateUser %+v -> %w", a, err)
	}
	return u, nil
}

// DeleteUserByID deletes the user specified by id.
// If the corresponding user does not exist, return ErrUserNotFound.
func (uc *UseCase) DeleteUserByID(ctx context.Context, id string) error {
	err := uc.r.DeleteUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("UseCase.DeleteUserByID %q -> %w", id, err)
	}
	return nil
}

// AddAuthentication associates an authentication with the user specified by id.
// If the user does not exist, return ErrUserNotFound.
// If the authentication is already registered, return ErrAlreadyExists.
func (uc *UseCase) AddAuthentication(ctx context.Context, id string, a *entity.Authentication) error {
	err := uc.r.AddAuthentication(ctx, id, a)
	if err != nil {
		return fmt.Errorf("UseCase.AddAuthentication %q %+v -> %w", id, a, err)
	}
	return nil
}

func NewUseCase(r port.IRepository) UseCase {
	return UseCase{r: r}
}
