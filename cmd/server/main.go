package main

import (
	delivery "github.com/EwanValentine/serverless-api-example/users/deliveries/http"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	router, err := delivery.Routes()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Running on port: ", port)
	log.Panic(http.ListenAndServe(":"+port, router))
}
