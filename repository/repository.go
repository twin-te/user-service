package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/twin-te/user-service/entity"
	"github.com/twin-te/user-service/port"
)

// Repository is used to implement port.IRepository.
type Repository struct {
	db *sql.DB
}

func (r *Repository) CreateUser(ctx context.Context, u *entity.User) error {
	failed := func(err error) error {
		return fmt.Errorf("Repository.CreateUser %+v -> %w", u, err)
	}

	// Check if the given user id is not registered.
	query := "SELECT FROM users WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, u.ID).Scan()
	if err == nil {
		return failed(entity.ErrUserAlreadyExists)
	}
	if err != sql.ErrNoRows {
		return failed(err)
	}

	// Check if the given authentication is not registered.
	for _, a := range u.Authentications {
		_, err := r.getUserIDByAuthentication(ctx, a)
		if err == nil {
			return failed(entity.ErrAuthenticationAlreadyExists)
		}
		if !errors.Is(err, entity.ErrUserNotFound) {
			return failed(err)
		}
	}

	// Start Transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return failed(err)
	}
	defer tx.Rollback()

	// Insert into users
	query = "INSERT INTO users (id) VALUES ($1)"
	_, err = tx.ExecContext(ctx, query, u.ID)
	if err != nil {
		return failed(err)
	}

	// Insert into user_authentications
	query = "INSERT INTO user_authentications (user_id,provider,social_id) VALUES ($1,$2,$3)"
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return failed(err)
	}
	defer stmt.Close()
	for _, a := range u.Authentications {
		_, err = stmt.ExecContext(ctx, u.ID, a.Provider, a.SocialID)
		if err != nil {
			return failed(err)
		}
	}

	// Commit Transaction
	if err = tx.Commit(); err != nil {
		return failed(err)
	}

	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	failed := func(err error) (*entity.User, error) {
		return nil, fmt.Errorf("Repository.GetUserByID %q -> %w", id, err)
	}

	u := &entity.User{}
	query := "SELECT id FROM users WHERE id = $1 AND \"deletedAt\" IS NULL"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return failed(entity.ErrUserNotFound)
		}
		return failed(err)
	}

	u.Authentications, err = r.getAuthenticationsByUserID(ctx, id)
	if err != nil {
		return failed(err)
	}

	return u, nil
}

func (r *Repository) GetUserByAuthentication(ctx context.Context, a *entity.Authentication) (*entity.User, error) {
	failed := func(err error) (*entity.User, error) {
		return nil, fmt.Errorf("Repository.GetUserByAuthentication %+v -> %w", a, err)
	}

	id, err := r.getUserIDByAuthentication(ctx, a)
	if err != nil {
		return failed(err)
	}

	u, err := r.GetUserByID(ctx, id)
	if err != nil {
		return failed(err)
	}

	return u, nil
}

func (r *Repository) DeleteUserByID(ctx context.Context, id string) error {
	failed := func(err error) error {
		return fmt.Errorf("Repository.DeleteUserByID %q -> %w", id, err)
	}

	// Check if the given user id is existed.
	query := "SELECT FROM users WHERE id = $1 AND \"deletedAt\" IS NULL"
	err := r.db.QueryRowContext(ctx, query, id).Scan()
	if err == sql.ErrNoRows {
		return failed(entity.ErrUserNotFound)
	}
	if err != nil {
		return failed(err)
	}

	// Start Transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return failed(err)
	}
	defer tx.Rollback()

	query = "UPDATE users SET \"deletedAt\" = $1 WHERE id = $2"
	_, err = tx.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return failed(err)
	}

	query = "DELETE FROM user_authentications WHERE user_id = $1"
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return failed(err)
	}

	if err = tx.Commit(); err != nil {
		return failed(err)
	}

	return nil
}

func (r *Repository) AddAuthentication(ctx context.Context, id string, a *entity.Authentication) error {
	failed := func(err error) error {
		return fmt.Errorf("Repository.AddAuthentication %q %+v -> %w", id, a, err)
	}

	query := "SELECT FROM users WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return failed(entity.ErrUserNotFound)
		}
		return failed(err)
	}

	_, err = r.getUserIDByAuthentication(ctx, a)
	if err == nil {
		return failed(entity.ErrAuthenticationAlreadyExists)
	}
	if !errors.Is(err, entity.ErrUserNotFound) {
		return failed(err)
	}

	query = "INSERT INTO user_authentications (user_id,provider,social_id) VALUES ($1,$2,$3)"
	_, err = r.db.ExecContext(ctx, query, id, a.Provider, a.SocialID)
	if err != nil {
		return failed(err)
	}

	return nil
}

// getUserIDByAuthentication returns the id of the user authorized by the given authentication.
// If the corresponding user does not exist, return ErrUserNotFound.
func (r *Repository) getUserIDByAuthentication(ctx context.Context, a *entity.Authentication) (id string, err error) {
	failed := func(err error) (string, error) {
		return id, fmt.Errorf("Repository.getUserIDByAuthentication %+v -> %w", a, err)
	}

	query := "SELECT user_id FROM user_authentications WHERE provider = $1 AND social_id = $2"
	err = r.db.QueryRowContext(ctx, query, a.Provider, a.SocialID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return failed(entity.ErrUserNotFound)
		}
		return failed(err)
	}

	return id, nil
}

// getAuthenticationsByUserID returns the authentications associated with the user specified by the given id.
func (r *Repository) getAuthenticationsByUserID(ctx context.Context, id string) ([]*entity.Authentication, error) {
	failed := func(err error) ([]*entity.Authentication, error) {
		return nil, fmt.Errorf("Repository.getAuthenticationsByUserID %q -> %w", id, err)
	}

	var sa []*entity.Authentication
	query := "SELECT provider, social_id FROM user_authentications WHERE user_id = $1"
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return failed(err)
	}
	defer rows.Close()
	for rows.Next() {
		a := &entity.Authentication{}
		if err := rows.Scan(&a.Provider, &a.SocialID); err != nil {
			return failed(err)
		}
		sa = append(sa, a)
	}
	if err := rows.Err(); err != nil {
		return failed(err)
	}
	return sa, nil
}

func NewRepository() (port.IRepository, func() error, error) {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	database := os.Getenv("PG_DATABASE")

	// https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, err
	}
	return &Repository{db: db}, db.Close, nil
}
