package users

import (
	"context"

	"go.uber.org/zap"
)

// LoggerAdapter wraps the usecase interface
// with a logging adapter which can be swapped out
type LoggerAdapter struct {
	Logger  *zap.Logger
	Usecase UserService
}

func (a *LoggerAdapter) logErr(err error) {
	if err != nil {
		a.Logger.Error(err.Error())
	}
}

// Get a single user
func (a *LoggerAdapter) Get(ctx context.Context, id string) (*User, error) {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("id", id))
	a.Logger.Info("getting a single user")
	user, err := a.Usecase.Get(ctx, id)
	a.logErr(err)
	return user, err
}

// GetAll gets all users
func (a *LoggerAdapter) GetAll(ctx context.Context) ([]*User, error) {
	defer a.Logger.Sync()
	a.Logger.Info("getting all users")
	users, err := a.Usecase.GetAll(ctx)
	a.logErr(err)
	return users, err
}

// Update a single user
func (a *LoggerAdapter) Update(ctx context.Context, id string, user *UpdateUser) error {
	defer a.Logger.Sync()
	a.Logger.With(zap.String("id", id))
	a.Logger.Info("updating a single user")
	err := a.Usecase.Update(ctx, id, user)
	a.logErr(err)
	return err
}

// Create a single user
func (a *LoggerAdapter) Create(ctx context.Context, user *User) error {
	defer a.Logger.Sync()
	a.Logger.Info("creating a single user")
	err := a.Usecase.Create(ctx, user)
	a.logErr(err)
	return err
}

// Delete a single user
func (a *LoggerAdapter) Delete(ctx context.Context, id string) error {
	defer a.Logger.Sync()
	a.Logger.Info("deleting a single user")
	err := a.Usecase.Delete(ctx, id)
	a.logErr(err)
	return err
}
