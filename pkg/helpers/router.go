package helpers

import (
	"context"
	"errors"
	"time"
)

type handler interface {
	GetAll(ctx context.Context) (Response, error)
	Get(ctx context.Context, id string) (Response, error)
	Create(ctx context.Context, body []byte) (Response, error)
	Update(ctx context.Context, id string, body []byte) (Response, error)
	Delete(ctx context.Context, id string) (Response, error)
}

const fiveSecondsTimeout = time.Second * 5

// Router takes a lambda handler and returns a higher order function
// which can route each request using the HTTP verb and path params.
func Router(handler handler) func(context.Context, Request) (Response, error) {
	return func(ctx context.Context, req Request) (Response, error) {

		// Add cancellation deadline to context
		ctx, cancel := context.WithTimeout(ctx, fiveSecondsTimeout)
		defer cancel()

		switch req.HTTPMethod {
		case "GET":
			id, ok := req.PathParameters["id"]
			if !ok {
				return handler.GetAll(ctx)
			}
			return handler.Get(ctx, id)

		case "POST":
			return handler.Create(ctx, []byte(req.Body))

		case "PUT":
			id, ok := req.PathParameters["id"]
			if !ok {
				return Response{}, errors.New("id parameter missing")
			}
			return handler.Update(ctx, id, []byte(req.Body))

		case "DELETE":
			id, ok := req.PathParameters["id"]
			if !ok {
				return Response{}, errors.New("id parameter missing")
			}
			return handler.Delete(ctx, id)

		default:
			return Response{}, errors.New("invalid method")
		}
	}
}
