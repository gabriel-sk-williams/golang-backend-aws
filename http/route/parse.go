package route

import (
	"strconv"
	"strings"
)

const (
	postModelQuery = `
		MATCH (player:Player {uuid: $puuid})-->(c:Circle)-->(space:Space {uuid: $suuid})
		WITH player, space
		MERGE (player)-[:SETS]->(model:Model)-[:FOR]->(space) SET model = {block}
		RETURN model
	`

	postPayoutQuery = `
		MATCH (player:Player {name: $name})-->(c:Circle)-->(space:Space {uuid: $suuid})
		WITH player, space
		MERGE (space)-[:SETS]->(payout:Payout)-[:FOR]->(player) SET payout = {block}
		RETURN payout
	`
)

// create string to add to cypher query
func formatProps(model map[string]float64, query string) string {

	spread := ``

	for str, val := range model {
		float := strconv.FormatFloat(val, 'f', 1, 64)
		line := str + ": " + float + ", "
		spread += line
	}

	props := strings.TrimSuffix(spread, ", ")
	final := strings.Replace(query, "block", props, 1)

	return final
}

func assertArray(list []interface{}) []string {

	array := make([]string, len(list))
	for i, v := range list {
		array[i] = v.(string)
	}
	return array
}

func assertModel(json map[string]interface{}) map[string]float64 {

	model := make(map[string]float64)
	for str, val := range json {
		model[str] = val.(float64)
	}
	return model //, error
}
