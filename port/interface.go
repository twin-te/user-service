package port

import (
	"context"

	"github.com/twin-te/user-service/entity"
)

// IRepository is the interface that groups the mothods requested by usecase.UseCase.
type IRepository interface {
	// CreateUser creates a new user.
	// If the user is already registered, return ErrAlreadyExists.
	CreateUser(ctx context.Context, u *entity.User) error

	// GetUserByID returns the user specified by id.
	// If the corresponding user does not exist, return ErrNotFound.
	GetUserByID(ctx context.Context, id string) (*entity.User, error)

	// GetUserByAuthentication returns the user associated with the specified authentication.
	// If the corresponding user does not exist, return ErrNotFound.
	GetUserByAuthentication(ctx context.Context, a *entity.Authentication) (*entity.User, error)

	// DeleteUserByID deletes the user specified by id.
	// If the corresponding user does not exist, return ErrNotFound.
	DeleteUserByID(ctx context.Context, id string) error

	// AddAuthentication associates an authentication with the user specified by id.
	// If the user does not exist, return ErrNotFound.
	// If the authentication is already registered, return ErrAlreadyExists.
	AddAuthentication(ctx context.Context, id string, a *entity.Authentication) error
}
