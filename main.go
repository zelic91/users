package main

import (
	"log"
	"zelic91/users/auth"
	"zelic91/users/gen/restapi"
	"zelic91/users/gen/restapi/operations"
	"zelic91/users/leaderboard"
	"zelic91/users/shared"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func authFunc(token string) (*shared.UserClaims, error) {
	return shared.ParseToken(token)
}

func loadConfig() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(err)
		}
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
	leaderboardRepo := leaderboard.NewRepo(db)

	authService := auth.Service{
		Repo: authRepo,
	}

	leaderboardService := leaderboard.NewService(leaderboardRepo)

	auth.SetupProfile(api)
	auth.SetupAuth(api, authService)
	leaderboard.SetupLeaderboard(api, leaderboardService)

	server := restapi.NewServer(api)
	server.Port = 9000

	defer server.Shutdown()

	// TODO: Setup middleware
	// server.SetHandler(api.Serve(nil))

	if err = server.Serve(); err != nil {
		log.Fatal(err)
	}
}
