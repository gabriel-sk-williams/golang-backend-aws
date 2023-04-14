<h2 align="center">Production Golang backend for AWS Elastic Beanstalk</h2>

- Uses [Goyave](https://goyave.dev/guide/installation.html) and the [offical Neo4j Go Driver](github.com/neo4j/neo4j-go-driver/v5/neo4j) to handle calls from [riverboat.zone](https://riverboat.zone) to Neo4j Aura.
- Handles typical users actions: login, joining and leaving spaces on the website

``` go
// start registration route
	if err := goyave.Start(func(router *goyave.Router) {
		router.CORS(cors.Default())
		router.Get("/", handler.Test)
		router.Get("/joined/{cuuid}", handler.ListJoined)
		router.Get("/models/{suuid}", handler.ListModels)
		router.Get("/space/{suuid}", handler.GetSpace)
		router.Get("/payouts/{suuid}", handler.ListPayouts)
		router.Post("/greeting", handler.Greeting)
		router.Post("/join", handler.Join).Validate(PlayerCircle)
		router.Post("/leave", handler.Leave).Validate(PlayerCircle)
		router.Post("/add_random", handler.AddRandom).Validate(Circle)
		router.Post("/submit", handler.SubmitModel).Validate(Submission)
		router.Post("/delete_model", handler.DeleteModel).Validate(PlayerSpace)
		router.Post("/calc", handler.CalculatePayouts).Validate(Space)

	}); err != nil {
		os.Exit(err.(*goyave.Error).ExitCode)
	}
```

- Also handles "payout" calculations whenever users edit or randomize their certainty for a given spread of outcomes
``` go
// receives Space
func (h Handler) CalculatePayouts(response *goyave.Response, r *goyave.Request) {
	fields := r.Data["fields"].([]string)
	pattern := r.String("pattern")
	stake := r.Numeric("stake")
	suuid := r.String("uuid")

	fmt.Println("calculating:", pattern)
	query := util.PostWBM // use pattern to match query

	models, _ := h.DB.mapModels(suuid)
	payouts, _ := calc.Payouts(models, fields, stake)
	result, err := h.DB.postPayouts(suuid, query, payouts)

	if err == nil {
		response.String(http.StatusOK, result)
	} else {
		response.String(http.StatusBadRequest, "Error: Could not calculate payouts.") // 400
	}
}
```
