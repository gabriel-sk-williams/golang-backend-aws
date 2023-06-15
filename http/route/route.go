package route

import (
	"fmt"
	"net/http"
	"riverboat/http/calc"
	"riverboat/model"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"goyave.dev/goyave/v4"
)

type Env struct {
	Driver neo4j.Driver
}

type Handler struct {
	DB Controls
}

// implements functions for structs
type Controls interface {
	getSpace(suuid string) (model.Space, error)
	listCircles() ([]model.Circle, error)
	listSpaces(cuuid string) ([]model.Space, error)
	listJoined(cuuid string) ([]model.Player, error)
	listModels(suuid string) (map[string]map[string]float64, error)
	listPayouts(suuid string) (map[string]map[string]float64, error)
	deleteModel(puuid string, suuid string) (string, error)
	addRandom(cuuid string) (string, error)
	join(puuid string, cuuid string) (string, error)
	leave(puuid string, cuuid string) (string, error)
	mapModels(suuid string) (map[string]map[string]float64, error)
	submitModel(puuid string, suuid string, json map[string]float64) (string, error)
	postPayouts(suuid string, payouts map[string]map[string]float64) (string, error)
	getStatus() error
}

//
// Initialization & Connection Functions
//

func (h Handler) GetStatus(response *goyave.Response, r *goyave.Request) {
	err := h.DB.getStatus()
	if err == nil {
		response.String(http.StatusOK, "online")
	} else {
		response.String(http.StatusOK, "offline")
	}
}

func (h Handler) Greeting(response *goyave.Response, r *goyave.Request) {
	response.String(http.StatusOK, "Welcome!")
}

//
// GET Functions
//

func (h Handler) ListCircles(response *goyave.Response, r *goyave.Request) {
	circles, err := h.DB.listCircles() // public circles
	if err == nil {
		response.JSON(http.StatusOK, circles)
	}
}

func (h Handler) ListSpaces(response *goyave.Response, r *goyave.Request) {
	spaces, err := h.DB.listSpaces(r.Params["cuuid"]) // spawned by circle
	if err == nil {
		response.JSON(http.StatusOK, spaces)
	}
}

func (h Handler) GetSpace(response *goyave.Response, r *goyave.Request) {
	space, err := h.DB.getSpace(r.Params["suuid"])
	if err == nil {
		response.JSON(http.StatusOK, space)
	}
}

func (h Handler) ListJoined(response *goyave.Response, r *goyave.Request) {
	joined, err := h.DB.listJoined(r.Params["cuuid"])
	if err == nil {
		response.JSON(http.StatusOK, joined)
	}
}

func (h Handler) ListModels(response *goyave.Response, r *goyave.Request) {
	models, err := h.DB.listModels(r.Params["suuid"]) // joined/suuid
	if err == nil {
		response.JSON(http.StatusOK, models)
	}
}

func (h Handler) ListPayouts(response *goyave.Response, r *goyave.Request) {
	payouts, err := h.DB.listPayouts(r.Params["suuid"])
	if err == nil {
		response.JSON(http.StatusOK, payouts)
	}
}

//
// POST Functions
//

// receives Space
func (h Handler) CalculatePayouts(response *goyave.Response, r *goyave.Request) {
	fields := r.Data["fields"].([]string)
	pattern := r.String("pattern")
	stake := r.Numeric("stake")
	suuid := r.String("uuid")

	fmt.Println("calculating:", pattern)

	models, _ := h.DB.mapModels(suuid)
	payouts, _ := calc.Payouts(models, fields, stake)
	result, err := h.DB.postPayouts(suuid, payouts) // change query

	if err == nil {
		response.String(http.StatusOK, result)
	} else {
		response.String(http.StatusBadRequest, "Error: Could not calculate payouts.") // 400
	}
}

// receives Submission
func (h Handler) SubmitModel(response *goyave.Response, r *goyave.Request) {
	puuid := r.String("puuid")
	suuid := r.String("suuid")
	spread := r.Object("model")

	model := assertModel(spread)

	res, err := h.DB.submitModel(puuid, suuid, model)

	if err == nil {
		response.String(http.StatusOK, res)
	} else {
		response.String(http.StatusBadRequest, "Error: Bad submission.") // 400
	}
}

func Read(spread map[string]float64) {
	fmt.Println("spread", spread)
}

// receives PlayerCircle
func (h Handler) Join(response *goyave.Response, r *goyave.Request) {
	puuid := r.String("puuid")
	cuuid := r.String("cuuid")

	res, err := h.DB.join(puuid, cuuid)

	if err == nil {
		response.String(http.StatusOK, res)
	} else {
		response.String(http.StatusBadRequest, "Error: Could not join Circle.") // 400
	}
}

// receives PlayerCircle
func (h Handler) Leave(response *goyave.Response, r *goyave.Request) {
	puuid := r.String("puuid")
	cuuid := r.String("cuuid")

	res, err := h.DB.leave(puuid, cuuid)

	if err == nil {
		response.String(http.StatusOK, res)
	} else {
		response.String(http.StatusBadRequest, "Error: Could not leave Circle.") // 400
	}
}

// receives Circle
func (h Handler) AddRandom(response *goyave.Response, r *goyave.Request) {
	cuuid := r.String("cuuid")

	res, err := h.DB.addRandom(cuuid)

	if err == nil {
		response.String(http.StatusOK, res)
	} else {
		response.String(http.StatusBadRequest, "Error: Could not join Circle.") // 400
	}
}

// receives PlayerSpace
func (h Handler) DeleteModel(response *goyave.Response, r *goyave.Request) {
	puuid := r.String("puuid")
	suuid := r.String("suuid")

	res, err := h.DB.deleteModel(puuid, suuid)

	if err == nil {
		response.String(http.StatusOK, res)
	} else {
		response.String(http.StatusBadRequest, "Error: Could not join delete Model") // 400
	}
}
