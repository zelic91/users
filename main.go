package main

import (
	"log"
	"zelic91/users/gen/models"
	"zelic91/users/gen/restapi"
	"zelic91/users/gen/restapi/operations"
	"zelic91/users/swagger"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/sirupsen/logrus"
)

func authFunc(token string) (*models.User, error) {
	return &models.User{
		Email:    "test@email.com",
		ExtIID:   "1234",
		Username: "tester",
	}, nil
}

func main() {
	// ctx := context.Background()

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")

	if err != nil {
		log.Fatal(err)
	}

	api := operations.NewUsersAPI(swaggerSpec)

	api.Logger = logrus.Printf

	api.AuthorizationAuth = authFunc
	api.ApplicationJSONConsumer = runtime.JSONConsumer()
	api.ApplicationJSONProducer = runtime.JSONProducer()

	swagger.Profile(api)

	server := restapi.NewServer(api)
	server.Port = 9999

	defer server.Shutdown()

	// TODO: Setup middleware
	// server.SetHandler(api.Serve(nil))

	if err = server.Serve(); err != nil {
		log.Fatal(err)
	}
}
