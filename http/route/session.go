package route

import (
	"riverboat/model"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/*
The methods ExecuteRead and ExecuteWrite have replaced ReadTransaction and WriteTransaction, which are deprecated in version 5.x and will be removed in version 6.0.
*/

func (env Env) getStatus() error {
	return env.Driver.VerifyConnectivity()
}

func (env Env) listCircles() ([]model.Circle, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (circle:Circle) 
			RETURN circle
		`, map[string]interface{}{})

		if err != nil {
			return nil, err
		}

		var circles []model.Circle
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("circle"); ok {
				node := value.(neo4j.Node)
				props := node.Props
				circle := model.Circle{
					Name: props["name"].(string),
					Uuid: props["uuid"].(string),
					// all_joined	false
					// all_modeled	false
					// all_paid	false
					// name	The Lab
					// spawned_by	Yakub
				}
				circles = append(circles, circle)
			}
		}

		return circles, err
	})

	return records.([]model.Circle), err
}

func (env Env) listSpaces(cuuid string) ([]model.Space, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (space:Space)<--(c:Circle {uuid: $cuuid})
			RETURN space
		`, map[string]interface{}{"cuuid": cuuid})

		if err != nil {
			return nil, err
		}

		var spaces []model.Space
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("space"); ok {
				node := value.(neo4j.Node)
				props := node.Props
				fields := props["fields"].([]interface{})
				space := model.Space{
					Fields:      assertArray(fields),
					Name:        props["name"].(string),
					Pattern:     props["pattern"].(string),
					Stake:       props["stake"].(float64),
					Uuid:        props["uuid"].(string),
					Description: props["description"].(string),
				}
				spaces = append(spaces, space)
			}
		}

		return spaces, err
	})

	return records.([]model.Space), err
}

func (env Env) getSpace(suuid string) (model.Space, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (space:Space {uuid: $suuid}) 
			RETURN space
		`, map[string]interface{}{"suuid": suuid})

		if err != nil {
			return nil, err
		}

		var space model.Space
		if result.Next() {
			record := result.Record()
			if value, ok := record.Get("space"); ok {
				node := value.(neo4j.Node)
				props := node.Props
				fields := props["fields"].([]interface{})
				object := model.Space{
					Fields:  assertArray(fields),
					Pattern: props["pattern"].(string),
					Stake:   props["stake"].(float64),
					Uuid:    props["uuid"].(string),
				}
				space = object
			}
		}
		return space, err
	})

	return records.(model.Space), err
}

// returns array of players
func (env Env) listJoined(cuuid string) ([]model.Player, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	people, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (player:Player)-[:JOINED]->(c:Circle {uuid: $cuuid})
			RETURN player
			`, map[string]interface{}{"cuuid": cuuid})

		if err != nil {
			return nil, err
		}

		var joined []model.Player
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("player"); ok {
				node := value.(neo4j.Node)
				props := node.Props
				player := model.Player{
					Name:  props["name"].(string),
					Uuid:  props["uuid"].(string),
					Money: props["money"].(float64),
					Risk:  props["risk"].(int64),
				}
				joined = append(joined, player)
			}
		}

		if err = result.Err(); err != nil {
			return nil, err
		}

		return joined, nil
	})

	if err != nil {
		return nil, err
	}

	return people.([]model.Player), nil
}

// func (env Env) listModels(suuid string) (map[string]model.WinByMethod, error) {
func (env Env) listModels(suuid string) (map[string]map[string]float64, error) {

	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	people, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (player:Player)-->(model:Model)-->(s:Space {uuid: $suuid})
			RETURN player, model
			`, map[string]interface{}{"suuid": suuid})

		if err != nil {
			return nil, err
		}

		var modelMap = make(map[string]map[string]float64)
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("player"); ok {
				node := value.(neo4j.Node)
				name := node.Props["name"]
				if val, err := record.Get("model"); err {
					modelNode := val.(neo4j.Node)
					props := modelNode.Props
					model := make(map[string]float64)

					for str, val := range props {
						model[str] = val.(float64)
					}

					modelMap[name.(string)] = model
				}
			}
		}

		if err = result.Err(); err != nil {
			return nil, err
		}

		return modelMap, nil
	})

	if err != nil {
		return nil, err
	}

	return people.(map[string]map[string]float64), nil
}

func (env Env) listPayouts(suuid string) (map[string]map[string]float64, error) {

	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	payouts, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (player:Player)<--(payout:Payout)<--(s:Space {uuid: $suuid})
			RETURN player, payout
			`, map[string]interface{}{"suuid": suuid})

		if err != nil {
			return nil, err
		}

		var payoutMap = make(map[string]map[string]float64)
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("player"); ok {
				node := value.(neo4j.Node)
				name := node.Props["name"]
				if val, err := record.Get("payout"); err {
					modelNode := val.(neo4j.Node)
					props := modelNode.Props
					payout := make(map[string]float64)

					for str, val := range props {
						payout[str] = val.(float64)
					}

					payoutMap[name.(string)] = payout
				}
			}
		}

		if err = result.Err(); err != nil {
			return nil, err
		}

		return payoutMap, nil
	})

	if err != nil {
		return nil, err
	}

	return payouts.(map[string]map[string]float64), nil
}

func (env Env) submitModel(puuid string, suuid string, json map[string]float64) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := formatProps(json, postModelQuery)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(query, map[string]interface{}{
			"puuid": puuid,
			"suuid": suuid,
		})

		if err != nil {
			return nil, err
		}

		return result.Collect() // Collects and commits
	})

	if err != nil {
		panic(err)
	}

	return "Model submitted.", nil
}

func (env Env) postPayouts(
	suuid string,
	payouts map[string]map[string]float64) (string, error) {

	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	for name, payout := range payouts {

		query := formatProps(payout, postPayoutQuery)

		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run(query, map[string]interface{}{
				"name":  name,
				"suuid": suuid,
			})

			if err != nil {
				return nil, err
			}

			return result.Collect() // Collects and commits
		})

		if err != nil {
			panic(err)
		}
	}

	return "Payouts posted.", nil
}

func (env Env) addRandom(cuuid string) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (c:Circle {uuid: $cuuid}) WITH c
			MATCH (p:Player) WHERE NOT (p)-[:JOINED]->(c) 
			WITH c, p, rand() as r ORDER BY r LIMIT 1
			MERGE (p)-[:JOINED]->(c)
			RETURN p
		`, map[string]interface{}{"cuuid": cuuid})

		if err != nil {
			return nil, err
		}

		return result.Collect() // Collects and commits
	})

	if err != nil {
		panic(err)
	}

	return "Player joined Circle.", nil
}

func (env Env) join(puuid string, cuuid string) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (p:Player {uuid: $puuid})
			WITH p MATCH (c:Circle {uuid: $cuuid})
			MERGE (p)-[:JOINED]->(c)
			RETURN p
		`, map[string]interface{}{"puuid": puuid, "cuuid": cuuid})

		if err != nil {
			return nil, err
		}

		return result.Collect() // Collects and commits
	})

	if err != nil {
		panic(err)
	}

	return "Player joined Circle.", nil
}

func (env Env) leave(puuid string, cuuid string) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (p:Player {uuid: $puuid})-[r:JOINED]->(c:Circle {uuid: $cuuid})
			DELETE r
		`, map[string]interface{}{"puuid": puuid, "cuuid": cuuid})

		if err != nil {
			return nil, err
		}

		return result.Collect() // Collects and commits
	})

	if err != nil {
		panic(err)
	}

	return "Player left Circle.", nil
}

func (env Env) deleteModel(puuid string, suuid string) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (p:Player {uuid: $puuid})--(n)--(s:Space {uuid: $suuid})
			WHERE n:Model OR n:Payout
			DETACH DELETE n
		`, map[string]interface{}{"puuid": puuid, "suuid": suuid})

		if err != nil {
			return nil, err
		}

		return result.Collect() // Collects and commits
	})

	if err != nil {
		panic(err)
	}

	return "Model deleted.", nil
}

func (env Env) mapModels(suuid string) (map[string]map[string]float64, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	people, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (player:Player)-->(model:Model)-->(s:Space {uuid: $suuid})
			RETURN player, model
			`, map[string]interface{}{"suuid": suuid})

		if err != nil {
			return nil, err
		}

		var modelMap = make(map[string]map[string]float64)
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("player"); ok {
				node := value.(neo4j.Node)
				name := node.Props["name"]
				if val, err := record.Get("model"); err {
					modelNode := val.(neo4j.Node)
					props := modelNode.Props
					model := make(map[string]float64)

					for str, val := range props {
						model[str] = val.(float64)
					}

					modelMap[name.(string)] = model
				}
			}
		}

		if err = result.Err(); err != nil {
			return nil, err
		}

		return modelMap, nil
	})

	if err != nil {
		return nil, err
	}

	return people.(map[string]map[string]float64), nil
}
