package users

import (
	"context"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type repository interface {
	Get(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, id string, user *User) error
	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// Usecase for interacting with users
type Usecase struct {
	repository repository
}

// Get a single user
func (u *Usecase) Get(ctx context.Context, id string) (*User, error) {
	user, err := u.repository.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching a single user")
	}
	return user, nil
}

// GetAll gets all users
func (u *Usecase) GetAll(ctx context.Context) ([]*User, error) {
	users, err := u.repository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching all users")
	}
	return users, nil
}

// Update a single user
func (u *Usecase) Update(ctx context.Context, id string, user *User) error {
	validate = validator.New()
	if err := validate.Struct(user); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	if err := u.repository.Update(ctx, id, user); err != nil {
		return errors.Wrap(err, "error updating user")
	}
	return nil
}

// Create a single user
func (u *Usecase) Create(ctx context.Context, user *User) error {
	validate = validator.New()
	if err := validate.Struct(user); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	if err := u.repository.Create(ctx, user); err != nil {
		return errors.Wrap(err, "error creating new user")
	}
	return nil
}

// Delete a single user
func (u *Usecase) Delete(ctx context.Context, id string) error {
	if err := u.repository.Delete(ctx, id); err != nil {
		return errors.Wrap(err, "error deleting user")
	}
	return nil
}
