package main

import (
	"fmt"
	"os"
	"riverboat/http/route"
	"riverboat/model"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"goyave.dev/goyave/v4"
	"goyave.dev/goyave/v4/cors"
)

func main() {
	fmt.Println("Goyave server active")
	godotenv.Load(".env")

	var uri, user, pw string

	if os.Getenv("APP_ENV") == "production" {
		fmt.Printf("Production environment: production \n")
		uri = "AURA_URI"
		user = "AURA_USERNAME"
		pw = "AURA_PASSWORD"
	} else {
		fmt.Printf("Production environment: development \n")
		uri = "DB_URI"
		user = "DB_USERNAME"
		pw = "DB_PASSWORD"
	}

	dbUri, found := os.LookupEnv(uri)
	//fmt.Println(dbUri)
	if !found {
		panic("DB_URI not set")
	}
	dbUser, found := os.LookupEnv(user)
	if !found {
		panic("DB_USERNAME not set")
	}
	dbPass, found := os.LookupEnv(pw)
	if !found {
		panic("DB_PASSWORD not set")
	}

	neoDriver := route.Env{
		Driver: driver(dbUri, dbUser, dbPass), // driver is thread-safe
	}

	handler := &route.Handler{
		DB: &neoDriver,
	}

	// start registration route
	if err := goyave.Start(func(router *goyave.Router) {
		router.CORS(cors.Default())
		router.Get("/", handler.GetStatus)
		router.Get("/circles", handler.ListCircles)       // public circles
		router.Get("/spaces/{cuuid}", handler.ListSpaces) // spawned by circle
		router.Get("/joined/{cuuid}", handler.ListJoined)
		router.Get("/models/{suuid}", handler.ListModels)
		router.Get("/space/{suuid}", handler.GetSpace)
		router.Get("/payouts/{suuid}", handler.ListPayouts)
		router.Post("/greeting", handler.Greeting)
		router.Post("/join", handler.Join).Validate(model.PlayerCircleProps)
		router.Post("/leave", handler.Leave).Validate(model.PlayerCircleProps)
		router.Post("/add_random", handler.AddRandom).Validate(model.CircleProps)
		router.Post("/submit", handler.SubmitModel).Validate(model.SubmissionProps)
		router.Post("/delete_model", handler.DeleteModel).Validate(model.PlayerSpaceProps)
		router.Post("/calc", handler.CalculatePayouts).Validate(model.SpaceProps)

	}); err != nil {
		os.Exit(err.(*goyave.Error).ExitCode)
	}
}

func driver(uri string, user string, pw string) neo4j.Driver {
	token := neo4j.BasicAuth(user, pw, "")
	result, err := neo4j.NewDriver(uri, token)
	if err != nil {
		panic(err)
	}
	return result
}
