package main

import (
	"encoding/json"
	"github.com/EwanValentine/serverless-api-example/pkg/helpers"
	"github.com/EwanValentine/serverless-api-example/users"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

type repository interface {
	Get(id string) (*users.User, error)
	GetAll() ([]*users.User, error)
	Update(id string, user *users.User) error
	Create(user *users.User) error
	Delete(id string) error
}

type handler struct {
	repository repository
}

func (h *handler) Get(id string) (helpers.Response, error) {
	user, err := h.repository.Get(id)
	if err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(user, http.StatusOK)
}

func (h *handler) GetAll() (helpers.Response, error) {
	users, err := h.repository.GetAll()
	if err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(users, http.StatusOK)
}

func (h *handler) Update(id string, body []byte) (helpers.Response, error) {
	user := &users.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	if err := h.repository.Update(id, user); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(map[string]interface{}{
		"success": true,
	}, http.StatusNoContent)
}

func (h *handler) Create(body []byte) (helpers.Response, error) {
	user := &users.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	if err := h.repository.Create(user); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(user, http.StatusCreated)
}

func (h *handler) Delete(id string) (helpers.Response, error) {
	if err := h.repository.Delete(id); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(map[string]interface{}{
		"success": true,
	}, http.StatusNoContent)
}

func main() {
	h := &handler{}
	lambda.Start(helpers.Router(h))
}
