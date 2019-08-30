package main

import (
	"context"
	"encoding/json"
	"github.com/EwanValentine/serverless-api-example/pkg/helpers"
	"github.com/EwanValentine/serverless-api-example/users"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"net/http"
	"os"
	"time"
)

type usecase interface {
	Get(ctx context.Context, id string) (*users.User, error)
	GetAll(ctx context.Context) ([]*users.User, error)
	Update(ctx context.Context, id string, user *users.User) error
	Create(ctx context.Context, user *users.User) error
	Delete(ctx context.Context, id string) error
}

type handler struct {
	usecase usecase
}

const fiveSecondsTimeout = time.Second*5

// Get a single user
func (h *handler) Get(id string) (helpers.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()
	user, err := h.usecase.Get(ctx, id)
	if err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(user, http.StatusOK)
}

// GetAll users
func (h *handler) GetAll() (helpers.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()
	users, err := h.usecase.GetAll(ctx)
	if err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(users, http.StatusOK)
}

// Update a single user
func (h *handler) Update(id string, body []byte) (helpers.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()
	user := &users.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	if err := h.usecase.Update(ctx, id, user); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(map[string]interface{}{
		"success": true,
	}, http.StatusNoContent)
}

// Create a user
func (h *handler) Create(body []byte) (helpers.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()
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
func (h *handler) Delete(id string) (helpers.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()
	if err := h.usecase.Delete(ctx, id); err != nil {
		return helpers.Fail(err, http.StatusInternalServerError)
	}

	return helpers.Success(map[string]interface{}{
		"success": true,
	}, http.StatusNoContent)
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	if err != nil {
		log.Fatal(err)
	}

	tableName := os.Getenv("TABLE_NAME")
	repository := users.NewDynamoDBRepository(dynamodb.New(sess), tableName)
	usecase := users.Usecase{repository}

	h := &handler{usecase}
	lambda.Start(helpers.Router(h))
}
