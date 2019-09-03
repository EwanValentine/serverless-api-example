package main

import (
	"context"
	"encoding/json"
	"github.com/EwanValentine/serverless-api-example/pkg/helpers"
	"github.com/EwanValentine/serverless-api-example/users"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"os"
	"testing"
)

var (
	id          = ""
	validUser   = `{ "name": "Test User", "email": "test@test.com", "age": 30 }`
	updatedUser = `{ "name": "Updated User", "email": "test@test.com", "age": 30 }`
)

func setup() *handler {
	os.Setenv("TABLE_NAME", "example-users-integration")
	usecase, err := users.Init(true)
	if err != nil {
		log.Panic(err)
	}

	h := &handler{usecase}
	return h
}

func clear() {
	os.Setenv("TABLE_NAME", "example-users-integration")
	usecase, err := users.Init(true)
	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()
	users, _ := usecase.GetAll(ctx)
	for _, user := range users {
		go usecase.Delete(ctx, user.ID)
	}
}

func TestCanCreate(t *testing.T) {
	ctx := context.Background()
	user := &users.User{}
	clear()
	h := setup()
	req := helpers.Request{
		HTTPMethod: "POST",
		Body:       validUser,
	}
	res, err := helpers.Router(h)(ctx, req)
	assert.NoError(t, err)

	err = json.Unmarshal([]byte(res.Body), &user)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotNil(t, user.ID)
	id = user.ID
}

func TestCanGetAllUsers(t *testing.T) {
	ctx := context.Background()
	u := []*users.User{}
	h := setup()
	req := helpers.Request{
		HTTPMethod: "GET",
	}
	res, err := helpers.Router(h)(ctx, req)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(res.Body), &u)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, "Test User", u[0].Name)
}

func TestCanGetUser(t *testing.T) {
	ctx := context.Background()
	u := &users.User{}
	h := setup()
	req := helpers.Request{
		HTTPMethod: "GET",
		PathParameters: map[string]string{
			"id": id,
		},
	}
	res, err := helpers.Router(h)(ctx, req)
	err = json.Unmarshal([]byte(res.Body), &u)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", u.Name)
}

func TestCanUpdateUser(t *testing.T) {
	ctx := context.Background()
	r := map[string]interface{}{}
	h := setup()
	req := helpers.Request{
		HTTPMethod: "PUT",
		PathParameters: map[string]string{
			"id": id,
		},
		Body: updatedUser,
	}
	res, err := helpers.Router(h)(ctx, req)
	err = json.Unmarshal([]byte(res.Body), &r)
	assert.NoError(t, err)
	assert.Equal(t, true, r["success"])
}

func TestCanDeleteUser(t *testing.T) {
	ctx := context.Background()
	r := map[string]interface{}{}
	h := setup()
	req := helpers.Request{
		HTTPMethod: "DELETE",
		PathParameters: map[string]string{
			"id": id,
		},
	}
	res, err := helpers.Router(h)(ctx, req)
	err = json.Unmarshal([]byte(res.Body), &r)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	expected := true
	assert.Equal(t, expected, r["success"])
}
