package route

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

//map[string]map[int]

func (env Env) mapModels(suuid string) (map[string]map[string]interface{}, error) {
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

		var modelMap = make(map[string]map[string]interface{})
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("player"); ok {
				node := value.(neo4j.Node)
				name := node.Props["name"]
				if val, err := record.Get("model"); err {
					modelNode := val.(neo4j.Node)
					props := modelNode.Props
					wbm := map[string]interface{}{
						"a_by_dec": props["a_by_dec"].(float64),
						"a_by_ko":  props["a_by_ko"].(float64),
						"b_by_dec": props["b_by_dec"].(float64),
						"b_by_ko":  props["b_by_ko"].(float64),
						"draw_nc":  props["draw_nc"].(float64),
					}
					modelMap[name.(string)] = wbm
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

	return people.(map[string]map[string]interface{}), nil
}
