package main

import (
	"context"
	"encoding/json"
	"github.com/EwanValentine/serverless-api-example/pkg/helpers"
	"github.com/EwanValentine/serverless-api-example/users"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type usecase interface {
	Get(ctx context.Context, id string) (*users.User, error)
	GetAll(ctx context.Context) ([]*users.User, error)
	Update(ctx context.Context, id string, user *users.UpdateUser) error
	Create(ctx context.Context, user *users.User) error
	Delete(ctx context.Context, id string) error
}

type handler struct {
	usecase usecase
}

const fiveSecondsTimeout = time.Second*5

// Get a single user
func (h *handler) Get(id string) (helpers.Response, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.With(zap.String("id", id))
	logger.Info("fetching user")

	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()
	user, err := h.usecase.Get(ctx, id)
	if err != nil {
		logger.Error(err.Error())
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(user, http.StatusOK)
}

// GetAll users
func (h *handler) GetAll() (helpers.Response, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("fetching all users")

	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()
	users, err := h.usecase.GetAll(ctx)
	if err != nil {
		logger.Error(err.Error())
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(users, http.StatusOK)
}

// Update a single user
func (h *handler) Update(id string, body []byte) (helpers.Response, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	logger.With(zap.String("id", id))
	logger.Info("updating user")

	updateUser := &users.UpdateUser{}
	if err := json.Unmarshal(body, &updateUser); err != nil {
		logger.Error(err.Error())
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	if err := h.usecase.Update(ctx, id, updateUser); err != nil {
		logger.Error(err.Error())
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(map[string]interface{}{
		"success": true,
	}, http.StatusNoContent)
}

// Create a user
func (h *handler) Create(body []byte) (helpers.Response, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	logger.Info("creating user")

	user := &users.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		logger.Error(err.Error())
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	if err := h.usecase.Create(ctx, user); err != nil {
		logger.Error(err.Error())
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(user, http.StatusCreated)
}

// Delete a user
func (h *handler) Delete(id string) (helpers.Response, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	logger.With(zap.String("id", id))
	if err := h.usecase.Delete(ctx, id); err != nil {
		logger.Error(err.Error())
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(map[string]interface{}{
		"success": true,
	}, http.StatusNoContent)
}

func main() {
	usecase, err := users.Init()
	if err != nil {
		log.Panic(err)
	}

	h := &handler{usecase}
	lambda.Start(helpers.Router(h))
}
