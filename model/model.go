package model

type Space struct {
	Fields      []string `json:"fields"`
	Name        string   `json:"name"`
	Pattern     string   `json:"pattern"`
	Stake       float64  `json:"stake"`
	Uuid        string   `json:"uuid"`
	Description string   `json:"description"`
}

type Circle struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
	// all_joined	false
	// all_modeled	false
	// all_paid	false
	// spawned_by	Yakub
	// uuid	1251094a-b643-4ccb-b12e-081c38ddb700
}

type Player struct {
	Name  string  `json:"name"`
	Uuid  string  `json:"uuid"`
	Money float64 `json:"money"`
	Risk  int64   `json:"risk"`
}
