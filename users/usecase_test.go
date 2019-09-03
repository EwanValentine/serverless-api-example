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
	expected := &User{
		Name:  "testing",
		Email: "test@test.com",
		Age:   30,
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)
	repo.EXPECT().Create(context.Background(), expected).Return(nil)

	uc := Usecase{repo}
	err := uc.Create(context.Background(), expected)

	assert.NoError(t, err)
}

func TestCanValidateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)

	uc := Usecase{repo}

	users := []*User{
		&User{},                      // No required fields
		&User{Name: "", Age: 0},      // Blank name
		&User{Name: "123", Age: 200}, // Integers as name, age too high
		&User{Email: "nope"},
	}
	for _, val := range users {
		err := uc.Create(context.Background(), val)
		assert.Error(t, err)
	}
}

func TestCanUpdateUser(t *testing.T) {
	user := &UpdateUser{
		Name:  "new name",
		Email: "test@test.com",
		Age:   20,
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)
	repo.EXPECT().Update(context.Background(), "abc123", user).Return(nil)
	uc := Usecase{repo}
	err := uc.Update(context.Background(), "abc123", user)
	assert.NoError(t, err)
}

func TestCanDeleteUser(t *testing.T) {
	user := &User{
		ID: "abc123",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockrepository(ctrl)
	repo.EXPECT().Delete(context.Background(), user.ID).Return(nil)
	uc := Usecase{repo}
	err := uc.Delete(context.Background(), user.ID)
	assert.NoError(t, err)
}
