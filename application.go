package main

import (
	"fmt"
	"os"
	"riverboat/http/route"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"goyave.dev/goyave/v4"
	"goyave.dev/goyave/v4/cors"
	"goyave.dev/goyave/v4/validation"
)

func main() {
	fmt.Println("goyave server activated")
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
	fmt.Println(dbUri)
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

	// SubmitModel()
	var (
		SubmitRequest = validation.RuleSet{
			"player": validation.List{"required", "string"},
			"space":  validation.List{"required", "string"},
			"model":  validation.List{"required", "object"},
		}
	)

	// GetTable()
	var (
		TableRequest = validation.RuleSet{
			"cuuid": validation.List{"required", "string"},
			"suuid": validation.List{"required", "string"},
		}
	)

	fmt.Println("env:", dbUri, dbUser, dbPass)

	// start registration route
	if err := goyave.Start(func(router *goyave.Router) {
		router.CORS(cors.Default())
		router.Get("/", route.Test)
		router.Get("/joined/{cuuid}", handler.ListJoined)
		router.Get("/models/{suuid}", handler.ListModels)
		router.Get("/payouts/{suuid}", handler.ListPayouts)
		router.Get("/table", handler.GetTable).Validate(TableRequest)
		router.Post("/join", handler.Join)
		router.Post("/submit", handler.SubmitModel).Validate(SubmitRequest)
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
