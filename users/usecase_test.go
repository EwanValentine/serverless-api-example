package users

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanGetUser(t *testing.T) {
	expected := &User{Name: "Ewan"}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockrepository(ctrl)
	repo.EXPECT().Get(context.Background(), "abc123").Return(expected, nil)

	uc := Usecase{repo}

	user, err := uc.Get(context.Background(), "abc123")

	assert.NoError(t, err)
	assert.Equal(t, expected, user)
}

func TestCanGetAllUsers(t *testing.T) {
	expected := []*User{
		&User{Name: "test1", Age: 1},
		&User{Name: "test2", Age: 2},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockrepository(ctrl)
	repo.EXPECT().GetAll(context.Background()).Return(expected, nil)

	uc := Usecase{repo}

	users, err := uc.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, expected, users)
}

func TestCanCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)

	uc := Usecase{repo}
	err := uc.Create(context.Background(), &User{
		Name: "testing",
		Age: 30,
	})

	assert.NoError(t, err)
}

func TestCanValidateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)

	uc := Usecase{repo}

	users := map[string]*User{
		"missing name, missing age": &User{}, // No required fields
		"name must be greater than 1 characters long": &User{Name: "", Age: 0}, // Blank name
		"name must be a string, age must be lower than 150": &User{Name: "123", Age: 200}, // Integers as name, age too high
	}
	for message, val := range users {
		err := uc.Create(context.Background(), val)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), message)
	}
}

func TestCanUpdateUser(t *testing.T) {

}

func TestCanDeleteUser(t *testing.T) {

}