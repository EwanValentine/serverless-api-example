package http

import (
	"context"
	"encoding/json"
	"github.com/EwanValentine/serverless-api-example/users"
	"github.com/gorilla/mux"
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

const fiveSecondsTimeout = time.Second*5

type delivery struct {
	usecase usecase
}

func writeErr(w http.ResponseWriter, err error, logger *zap.Logger) {
	logger.Error(
		"failed to get a user by id",
	)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func (d *delivery) Get(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]

	logger.With(zap.String("id", id))
	logger.Info(
		"get single user",
	)

	user, err := d.usecase.Get(ctx, id)
	if err != nil {
		writeErr(w, err, logger)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		writeErr(w, err, logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (d *delivery) GetAll(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	users, err := d.usecase.GetAll(ctx)
	if err != nil {
		writeErr(w, err, logger)
		return
	}

	data, err := json.Marshal(users)
	if err != nil {
		writeErr(w, err, logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (d *delivery) Update(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	user := &users.UpdateUser{}
	if err := decoder.Decode(&user); err != nil {
		writeErr(w, err, logger)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := d.usecase.Update(ctx, id, user); err != nil {
		writeErr(w, err, logger)
		return
	}
}

func (d *delivery) Create(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	user := &users.User{}
	if err := decoder.Decode(&user); err != nil {
		writeErr(w, err, logger)
		return
	}

	logger.Info("creating user")
	if err := d.usecase.Create(ctx, user); err != nil {
		writeErr(w, err, logger)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Ok"))
}

func (d *delivery) Delete(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]
	logger.With(zap.String("id", id))
	logger.Info("deleting user")

	if err := d.usecase.Delete(ctx, id); err != nil {
		writeErr(w, err, logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Deleted"))
}

func Routes() (*mux.Router, error) {
	usecase, err := users.Init()
	if err != nil {
		log.Panic(err)
	}

	delivery := &delivery{usecase}

	r := mux.NewRouter()
	r.HandleFunc("/users", delivery.Create).Methods("POST")
	r.HandleFunc("/users", delivery.GetAll).Methods("GET")
	r.HandleFunc("/users/{id}", delivery.Get).Methods("GET")
	r.HandleFunc("/users/{id}", delivery.Update).Methods("PUT")
	r.HandleFunc("/users/{id}", delivery.Delete).Methods("DELETE")

	return r, nil
}
