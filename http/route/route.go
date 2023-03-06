package route

import (
	"fmt"
	"net/http"
	"riverboat/http/calc"
	"riverboat/model"
	"riverboat/util"

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
	listJoined(cuuid string) ([]model.Player, error)
	listModels(suuid string) (map[string]model.WinByMethod, error)
	listPayouts(suuid string) (map[string]model.WinByMethod, error)
	submitModel(puuid string, suuid string, json model.WinByMethod) (string, error)
	deleteModel(puuid string, suuid string) (string, error)
	addRandom(cuuid string) (string, error)
	join(puuid string, cuuid string) (string, error)
	leave(puuid string, cuuid string) (string, error)
	mapModels(suuid string) (map[string]map[string]interface{}, error)
	postPayouts(suuid string, query string, payouts map[string]map[string]float64) (string, error)
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

func (h Handler) CalculatePayouts(response *goyave.Response, r *goyave.Request) {
	fields := r.Data["fields"].([]string)
	pattern := r.String("pattern")
	stake := r.Numeric("stake")
	suuid := r.String("uuid")

	fmt.Println(pattern) // use pattern to match
	query := util.PostWBM

	models, _ := h.DB.mapModels(suuid)
	payouts, _ := calc.Payouts(models, fields, stake)
	result, err := h.DB.postPayouts(suuid, query, payouts)

	if err == nil {
		response.String(http.StatusOK, result)
	}
}

// receives Submission
func (h Handler) SubmitModel(response *goyave.Response, r *goyave.Request) {
	puuid := r.String("puuid")
	suuid := r.String("suuid")
	spread := r.Object("model")

	model := &model.WinByMethod{
		AByDec: spread["a_by_dec"].(float64),
		AByKO:  spread["a_by_ko"].(float64),
		BByDec: spread["b_by_dec"].(float64),
		BByKO:  spread["b_by_ko"].(float64),
		DrawNC: spread["draw_nc"].(float64),
	}

	text, err := h.DB.submitModel(puuid, suuid, *model)

	if err == nil {
		response.String(http.StatusOK, text)
	} else {
		fmt.Println(err)
	}
}

func (h Handler) Join(response *goyave.Response, r *goyave.Request) {
	puuid := r.String("puuid")
	cuuid := r.String("cuuid")

	res, err := h.DB.join(puuid, cuuid)

	if err == nil {
		response.String(http.StatusOK, res)
	} else {
		fmt.Println(err)
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
		fmt.Println(err)
	}
}

// receives Circle
func (h Handler) AddRandom(response *goyave.Response, r *goyave.Request) {
	cuuid := r.String("cuuid")

	res, err := h.DB.addRandom(cuuid)
	fmt.Println("adding random", cuuid)
	if err == nil {
		response.String(http.StatusOK, res)
	} else {
		fmt.Println(err)
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
		fmt.Println(err)
	}

}

func Test(response *goyave.Response, r *goyave.Request) {
	response.String(http.StatusOK, "d e e p   n u m b e r s")
}

func (h Handler) Greeting(response *goyave.Response, r *goyave.Request) {
	response.String(http.StatusOK, "Welcome!")
}
