package model

type WinByMethod struct {
	AByDec float64 `json:"a_by_dec"`
	AByKO  float64 `json:"a_by_ko"`
	BByDec float64 `json:"b_by_dec"`
	BByKO  float64 `json:"b_by_ko"`
	DrawNC float64 `json:"draw_nc"`
}

type Space struct {
	Fields  []string `json:"fields"`
	Pattern string   `json:"pattern"`
	Stake   float64  `json:"stake"`
	Uuid    string   `json:"uuid"`
}

type Player struct {
	Name  string  `json:"name"`
	Uuid  string  `json:"uuid"`
	Money float64 `json:"money"`
	Risk  int64   `json:"risk"`
}
