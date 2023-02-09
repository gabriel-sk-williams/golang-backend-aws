package route

import (
	"fmt"
	"net/http"

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
	listJoined(circle string) ([]Player, error)
	listModels(space string) (map[string]WinByMethod, error)
	listPayouts(space string) (map[string]WinByMethod, error)
	submitModel(player string, space string, json Model) (string, error)
}

type Model struct {
	AByDec float64 `json:"a_by_dec"`
	AByKO  float64 `json:"a_by_ko"`
	BByDec float64 `json:"b_by_dec"`
	BByKO  float64 `json:"b_by_ko"`
	DrawNC float64 `json:"draw_nc"`
}

type Submission struct {
	Player float64 `json:"player"`
	Space  float64 `json:"space"`
	Model  Model   `json:"model"`
}

/*
	AByDec float64 `json:"a_by_dec"`
	AByKO  float64 `json:"a_by_ko"`
	BByDec float64 `json:"b_by_dec"`
	BByKO  float64 `json:"b_by_ko"`
	DrawNC float64 `json:"draw_nc"`

	"a_by_dec": validation.List{"required", "numeric"},
	"a_by_ko":  validation.List{"required", "numeric"},
	"b_by_dec": validation.List{"required", "numeric"},
	"b_by_ko":  validation.List{"required", "numeric"},
	"draw_nc":  validation.List{"required", "numeric"},
*/

func (h Handler) ListJoined(response *goyave.Response, r *goyave.Request) {
	joined, err := h.DB.listJoined(r.Params["cuuid"])
	if err == nil {
		response.JSON(http.StatusOK, joined)
	}
}

func (h Handler) ListModels(response *goyave.Response, r *goyave.Request) {
	models, err := h.DB.listModels(r.Params["suuid"])
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

// need suuid and model
func (h Handler) SubmitModel(response *goyave.Response, r *goyave.Request) {
	player := r.String("player")
	space := r.String("space")
	certs := r.Object("model")

	model := &Model{
		AByDec: certs["a_by_dec"].(float64),
		AByKO:  certs["a_by_ko"].(float64),
		BByDec: certs["b_by_dec"].(float64),
		BByKO:  certs["b_by_ko"].(float64),
		DrawNC: certs["draw_nc"].(float64),
	}

	text, err := h.DB.submitModel(player, space, *model)
	fmt.Println("submission return:", text) //http response here
	if err == nil {
		response.JSON(http.StatusOK, "positive response")
	} else {
		fmt.Println(err)
	}
}

// axios post with NameCircle { name, cuuid }
// axios.post(urlToggle, NameCircle, { headers: defaultHeader })
func (h Handler) Join(response *goyave.Response, r *goyave.Request) {

	// joined, err := h.DB.getAll("1251094a-b643-4ccb-b12e-081c38ddb700")
	fmt.Println("can i join pwease??")
	//if err == nil {
	response.JSON(http.StatusOK, "ur boy joined")
	//}
}

func Test(response *goyave.Response, r *goyave.Request) {
	response.String(http.StatusOK, "d e e p   n u m b e r s")
}
