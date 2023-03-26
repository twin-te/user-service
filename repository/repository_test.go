package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/twin-te/user-service/entity"
	"github.com/twin-te/user-service/port"
)

func compareUser(t *testing.T, u1, u2 *entity.User) {
	t.Helper()

	if u1.ID != u2.ID {
		t.Fatalf("ID does not match")
	}

	if len(u1.Authentications) != len(u2.Authentications) {
		t.Fatalf("authentications do not match")
	}

	for i, auth := range u2.Authentications {
		if auth.Provider != u1.Authentications[i].Provider || auth.SocialID != u1.Authentications[i].SocialID {
			t.Fatalf("authentication at index %d does not match", i)
		}
	}
}

func createUser(t *testing.T, repo port.IRepository) entity.User {
	t.Helper()

	user := entity.User{
		ID: uuid.NewString(),
		Authentications: []*entity.Authentication{{
			Provider: entity.ProviderGoogle,
			SocialID: uuid.NewString(),
		}},
	}

	if err := repo.CreateUser(context.Background(), &user); err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	return user
}

func TestCreateUser(t *testing.T) {
	repo, close, err := NewRepository()
	if err != nil {
		t.Fatalf("Error creating repository: %v", err)
	}
	defer close()

	user := createUser(t, repo)

	retrievedUser, err := repo.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("Error retrieving user: %v", err)
	}

	compareUser(t, retrievedUser, &user)

	if err := repo.CreateUser(context.Background(), &entity.User{ID: user.ID}); !errors.Is(err, entity.ErrAlreadyExists) {
		t.Fatalf("Error creating user whose id have already registered: %v", err)
	}

	if err := repo.CreateUser(context.Background(), &entity.User{ID: uuid.NewString(), Authentications: user.Authentications}); !errors.Is(err, entity.ErrAlreadyExists) {
		t.Fatalf("Error creating user whose authentications have already registered: %v", err)
	}
}

func TestGetUserByID(t *testing.T) {
	repo, close, err := NewRepository()
	if err != nil {
		t.Fatalf("Error creating repository: %v", err)
	}
	defer close()

	if _, err = repo.GetUserByID(context.Background(), uuid.NewString()); !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("Error retrieving not created user: %v", err)
	}

	expectedUser := createUser(t, repo)

	retrievedUser, err := repo.GetUserByID(context.Background(), expectedUser.ID)
	if err != nil {
		t.Fatalf("Error retrieving user: %v", err)
	}

	compareUser(t, retrievedUser, &expectedUser)
}

func TestGetUserByAuthentication(t *testing.T) {
	repo, close, err := NewRepository()
	if err != nil {
		t.Fatalf("Error creating repository: %v", err)
	}
	defer close()

	_, err = repo.GetUserByAuthentication(context.Background(), &entity.Authentication{Provider: entity.ProviderGoogle, SocialID: uuid.NewString()})
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("Error retrieving not created user: %v", err)
	}

	expectedUser := createUser(t, repo)

	retrievedUser, err := repo.GetUserByAuthentication(context.Background(), expectedUser.Authentications[0])
	if err != nil {
		t.Fatalf("Error retrieving user: %v", err)
	}

	compareUser(t, retrievedUser, &expectedUser)
}

func TestDeleteUserByID(t *testing.T) {
	repo, close, err := NewRepository()
	if err != nil {
		t.Fatalf("Error creating repository: %v", err)
	}
	defer close()

	expectedUser := createUser(t, repo)

	if err := repo.DeleteUserByID(context.Background(), expectedUser.ID); err != nil {
		t.Fatalf("Error deleting user: %v", err)
	}

	_, err = repo.GetUserByID(context.Background(), expectedUser.ID)
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("Error retrieving deleted user: %v", err)
	}

	err = repo.DeleteUserByID(context.Background(), uuid.NewString())
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("Error retrieving not created user: %v", err)
	}

	err = repo.CreateUser(context.Background(), &entity.User{ID: uuid.NewString(), Authentications: expectedUser.Authentications})
	if err != nil {
		t.Fatalf("Error creating user who has deleted authentications: %v", err)
	}
}

func TestAddAuthentication(t *testing.T) {
	repo, close, err := NewRepository()
	if err != nil {
		t.Fatalf("Error creating repository: %v", err)
	}
	defer close()

	user := entity.User{
		ID:              uuid.NewString(),
		Authentications: []*entity.Authentication{},
	}

	if err := repo.CreateUser(context.Background(), &user); err != nil {
		t.Fatalf("Error creating expected user: %v", err)
	}

	auth := &entity.Authentication{
		Provider: entity.ProviderGoogle,
		SocialID: uuid.NewString(),
	}

	if err := repo.AddAuthentication(context.Background(), uuid.NewString(), auth); !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("Error adding authentication to not created user: %v", err)
	}

	if err := repo.AddAuthentication(context.Background(), user.ID, auth); err != nil {
		t.Fatalf("Error adding authentication: %v", err)
	}

	retrievedUser, err := repo.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("Error retrieving user: %v", err)
	}

	if len(retrievedUser.Authentications) != 1 {
		t.Fatalf("Retrieved user authentications do not match expected user authentications")
	}

	if retrievedUser.Authentications[0].Provider != auth.Provider || retrievedUser.Authentications[0].SocialID != auth.SocialID {
		t.Fatalf("Retrieved user authentication does not match expected user authentication")
	}

	if err := repo.AddAuthentication(context.Background(), user.ID, auth); !errors.Is(err, entity.ErrAlreadyExists) {
		t.Fatalf("Error adding authentication already registered: %v", err)
	}
}
