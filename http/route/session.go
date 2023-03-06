package route

import (
	"fmt"
	"riverboat/model"
	"riverboat/util"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/*
The methods ExecuteRead and ExecuteWrite have replaced ReadTransaction and WriteTransaction, which are deprecated in version 5.x and will be removed in version 6.0.
*/

func (env Env) getSpace(suuid string) (model.Space, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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
					Fields:  util.AssertArray(fields),
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
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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

func (env Env) listModels(suuid string) (map[string]model.WinByMethod, error) {
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

		var modelMap = make(map[string]model.WinByMethod)
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("player"); ok {
				node := value.(neo4j.Node)
				name := node.Props["name"]
				if val, err := record.Get("model"); err {
					modelNode := val.(neo4j.Node)
					props := modelNode.Props
					wbm := model.WinByMethod{
						AByDec: props["a_by_dec"].(float64),
						AByKO:  props["a_by_ko"].(float64),
						BByDec: props["b_by_dec"].(float64),
						BByKO:  props["b_by_ko"].(float64),
						DrawNC: props["draw_nc"].(float64),
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

	return people.(map[string]model.WinByMethod), nil
}

func (env Env) listPayouts(suuid string) (map[string]model.WinByMethod, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	payouts, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (player:Player)<--(payout:Payout)<--(s:Space {uuid: $suuid})
			RETURN player, payout
			`, map[string]interface{}{"suuid": suuid})

		if err != nil {
			return nil, err
		}

		var payoutMap = make(map[string]model.WinByMethod)
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("player"); ok {
				node := value.(neo4j.Node)
				name := node.Props["name"]
				if val, err := record.Get("payout"); err {
					modelNode := val.(neo4j.Node)
					props := modelNode.Props
					wbm := model.WinByMethod{
						AByDec: props["a_by_dec"].(float64),
						AByKO:  props["a_by_ko"].(float64),
						BByDec: props["b_by_dec"].(float64),
						BByKO:  props["b_by_ko"].(float64),
						DrawNC: props["draw_nc"].(float64),
					}
					payoutMap[name.(string)] = wbm
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

	return payouts.(map[string]model.WinByMethod), nil
}

func (env Env) submitModel(puuid string, suuid string, json model.WinByMethod) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	records, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (player:Player {uuid: $puuid})-->(circle)-->(space:Space {uuid: $suuid})
			WITH player, space
			MERGE (player)-[:SETS]->(model:Model)-[:FOR]->(space) SET model = {
				a_by_dec: $adec,
				a_by_ko: $ako,
				b_by_dec: $bdec,
				b_by_ko: $bko,
				draw_nc: $draw
			}
			RETURN model
			`, map[string]interface{}{
			"puuid": puuid,
			"suuid": suuid,
			"adec":  json.AByDec,
			"ako":   json.AByKO,
			"bdec":  json.BByDec,
			"bko":   json.BByKO,
			"draw":  json.DrawNC,
		})

		if err != nil {
			return nil, err
		}

		return result.Collect() // Collects and commits
	})

	if err != nil {
		panic(err)
	}

	// gives use for records
	for _, record := range records.([]*neo4j.Record) {
		temp := record.Values[0].(neo4j.Node)
		fmt.Println("submit:", temp)
	}

	return "model submitted", nil
}

// json mode.WinByMethod
// runs calcs with puuid instead of name
func (env Env) postPayouts(
	suuid string,
	query string,
	payouts map[string]map[string]float64) (string, error) {

	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	for name, payout := range payouts {
		newp := map[string]interface{}{
			"name":  name,
			"suuid": suuid,
		}
		for val, float := range payout {
			newp[val] = float
		}

		records, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run(`
				MATCH (player:Player {name: $name})-->(circle)-->(space:Space {uuid: $suuid})
				WITH player, space
				MERGE (space)-[:SETS]->(payout:Payout)-[:FOR]->(player) SET payout = {
					a_by_dec: $a_by_dec,
					a_by_ko: $a_by_ko,
					b_by_dec: $b_by_dec,
					b_by_ko: $b_by_ko,
					draw_nc: $draw_nc
				}
				RETURN payout
				`, newp)

			if err != nil {
				return nil, err
			}

			return result.Collect() // Collects and commits
		})

		if err != nil {
			panic(err)
		}

		fmt.Println(records)
	}

	return "payouts posted", nil
}

/*
#[tokio::main]
pub async fn post_payout(name: String, suuid: String,
                         payout: WinByMethod, graph: &Graph)
                         -> tide::Result<()> {

    let WinByMethod { a_by_dec, a_by_ko, b_by_dec, b_by_ko, draw_nc } = payout;

    graph.run(
        Query::new("
            MATCH (player:Player {name: $name})-->(circle)-->(space:Space {uuid: $suuid})
            WITH player, space
            MERGE (space)-[:SETS]->(payout:Payout)-[:FOR]->(player) SET payout = {
                a_by_dec: $adec,
                a_by_ko: $ako,
                b_by_dec: $bdec,
                b_by_ko: $bko,
                draw_nc: $draw
            }
            RETURN payout
        ")
        .param("name", name)
        .param("suuid", suuid)
        .param("adec", a_by_dec)
        .param("ako", a_by_ko)
        .param("bdec", b_by_dec)
        .param("bko", b_by_ko)
        .param("draw", draw_nc)
        ).await.unwrap();

    Ok(())
}
*/

func (env Env) addRandom(cuuid string) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	records, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
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

	// gives use for records
	for _, record := range records.([]*neo4j.Record) {
		temp := record.Values[0].(neo4j.Node)
		fmt.Println("random:", temp)
	}

	return "random added", nil
}

func (env Env) join(puuid string, cuuid string) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	records, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
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

	// gives use for records
	for _, record := range records.([]*neo4j.Record) {
		temp := record.Values[0].(neo4j.Node)
		fmt.Println("join", temp)
	}

	return "joined", nil
}

func (env Env) leave(puuid string, cuuid string) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	records, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
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

	// gives use for records
	for _, record := range records.([]*neo4j.Record) {
		temp := record.Values[0].(neo4j.Node)
		fmt.Println("leave", temp)
	}

	return "left", nil
}

func (env Env) deleteModel(puuid string, suuid string) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	records, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
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

	// gives use for records
	for _, record := range records.([]*neo4j.Record) {
		temp := record.Values[0].(neo4j.Node)
		fmt.Println("delete:", temp)
	}

	return "model deleted", nil
}
