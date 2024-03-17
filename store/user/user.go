package user

import (
	"context"

	"github.com/billymosis/marketplace-app/model"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type userStore struct {
	db       *pgxpool.Pool
	validate *validator.Validate
}

func NewUserStore(db *pgxpool.Pool, validate *validator.Validate) model.UserStore {
	return &userStore{
		db:       db,
		validate: validate,
	}
}

func (us *userStore) GetValidator() *validator.Validate {
	return us.validate
}

func (us *userStore) GetById(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	query := "SELECT * FROM users WHERE id = $1 LIMIT 1"
	err := us.db.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Name,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by ID")
	}
	return &user, nil
}

func (us *userStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	query := "SELECT * FROM users WHERE username = $1 LIMIT 1"
	err := us.db.QueryRow(ctx, query, username).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Name,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by username")
	}
	return &user, nil
}

func (us *userStore) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := "INSERT INTO users (password, username, name) VALUES($1,$2,$3) RETURNING id"
	err := us.db.QueryRow(ctx, query, user.Password, user.Username, user.Name).Scan(&user.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}
	return user, nil
}
