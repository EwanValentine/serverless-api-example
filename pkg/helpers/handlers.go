package helpers

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

type Request events.APIGatewayProxyRequest

type Response events.APIGatewayProxyResponse

// Fail returns an internal server error with the error message
func Fail(err error, status int) (Response, error) {
	e := make(map[string]string, 0)
	e["message"] = err.Error()

	// We don't need to worry about this error,
	// as we're controlling the input.
	body, _ := json.Marshal(e)

	return Response{
		Body:       string(body),
		StatusCode: status,
	}, nil
}

// Success returns a valid response
func Success(data interface{}, status int) (Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return Fail(err, http.StatusInternalServerError)
	}

	return Response{
		Body:       string(body),
		StatusCode: status,
	}, nil
}
