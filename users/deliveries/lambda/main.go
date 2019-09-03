package main

import (
	"context"
	"encoding/json"
	"github.com/EwanValentine/serverless-api-example/pkg/helpers"
	"github.com/EwanValentine/serverless-api-example/users"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
)

type handler struct {
	usecase users.UserService
}

// Get a single user
func (h *handler) Get(ctx context.Context, id string) (helpers.Response, error) {
	user, err := h.usecase.Get(ctx, id)
	if err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(user, http.StatusOK)
}

// GetAll users
func (h *handler) GetAll(ctx context.Context) (helpers.Response, error) {
	users, err := h.usecase.GetAll(ctx)
	if err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(users, http.StatusOK)
}

// Update a single user
func (h *handler) Update(ctx context.Context, id string, body []byte) (helpers.Response, error) {
	updateUser := &users.UpdateUser{}
	if err := json.Unmarshal(body, &updateUser); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	if err := h.usecase.Update(ctx, id, updateUser); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(map[string]interface{}{
		"success": true,
	}, http.StatusNoContent)
}

// Create a user
func (h *handler) Create(ctx context.Context, body []byte) (helpers.Response, error) {
	user := &users.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	if err := h.usecase.Create(ctx, user); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(user, http.StatusCreated)
}

// Delete a user
func (h *handler) Delete(ctx context.Context, id string) (helpers.Response, error) {
	if err := h.usecase.Delete(ctx, id); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(map[string]interface{}{
		"success": true,
	}, http.StatusNoContent)
}

func main() {
	usecase, err := users.Init(false)
	if err != nil {
		log.Panic(err)
	}

	h := &handler{usecase}
	lambda.Start(helpers.Router(h))
}
