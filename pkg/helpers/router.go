package helpers

import "errors"

type handler interface {
	GetAll() (Response, error)
	Get(id string) (Response, error)
	Create(body []byte) (Response, error)
	Update(id string, body []byte) (Response, error)
	Delete(id string) (Response, error)
}

func Router(handler handler) func(Request) (Response, error) {
	return func(req Request) (Response, error) {
		switch req.HTTPMethod {
		case "GET":
			id, ok := req.PathParameters["id"]
			if !ok {
				return handler.GetAll()
			}
			return handler.Get(id)

		case "POST":
			return handler.Create([]byte(req.Body))

		case "PUT":
			id, ok := req.PathParameters["id"]
			if !ok {
				return Response{}, errors.New("id parameter missing")
			}
			return handler.Update(id, []byte(req.Body))

		case "DELETE":
			id, ok := req.PathParameters["id"]
			if !ok {
				return Response{}, errors.New("id parameter missing")
			}
			return handler.Delete(id)

		default:
			return Response{}, errors.New("invalid method")
		}
	}
}
