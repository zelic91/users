package main

import (
	"log"
	"zelic91/users/auth"
	"zelic91/users/gen/restapi"
	"zelic91/users/gen/restapi/operations"
	"zelic91/users/swagger"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func authFunc(token string) (*auth.UserClaims, error) {
	return &auth.UserClaims{
		ID:       1234,
		Username: "tester",
	}, nil
}

func loadConfig() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	loadConfig()
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

	dbUrl := viper.GetString("DATABASE_URL")

	log.Println(dbUrl)
	db, err := sqlx.Connect("postgres", dbUrl)

	if err != nil {
		log.Fatal(err)
	}

	authRepo := auth.NewRepo(db)

	authService := auth.AuthService{
		Repo: authRepo,
	}

	swagger.SetupProfile(api)
	swagger.SetupAuth(api, authService)

	server := restapi.NewServer(api)
	server.Port = 9000

	defer server.Shutdown()

	// TODO: Setup middleware
	// server.SetHandler(api.Serve(nil))

	if err = server.Serve(); err != nil {
		log.Fatal(err)
	}
}
