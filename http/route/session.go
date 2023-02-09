package route

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Player struct {
	Name  string  `json:"name"`
	Uuid  string  `json:"uuid"`
	Money float64 `json:"money"`
	Risk  int64   `json:"risk"`
}

type WinByMethod struct {
	AByDec float64 `json:"a_by_dec"`
	AByKO  float64 `json:"a_by_ko"`
	BByDec float64 `json:"b_by_dec"`
	BByKO  float64 `json:"b_by_ko"`
	DrawNC float64 `json:"draw_nc"`
}

// returns array of players
func (env Env) listJoined(circle string) ([]Player, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	people, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		result, err := tx.Run(`
				MATCH (player:Player)-[:JOINED]->(c:Circle {uuid: $uuid})
				RETURN player
				`, map[string]interface{}{"uuid": circle})

		if err != nil {
			return nil, err
		}

		var joined []Player
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("player"); ok {
				node := value.(neo4j.Node)
				props := node.Props
				player := Player{
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

	return people.([]Player), nil
}

func (env Env) listModels(space string) (map[string]WinByMethod, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	people, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		result, err := tx.Run(`
			MATCH (player:Player)-->(model:Model)-->(s:Space {uuid: $uuid})
			RETURN player, model
			`, map[string]interface{}{"uuid": space})

		if err != nil {
			return nil, err
		}

		var modelMap = make(map[string]WinByMethod)
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("player"); ok {
				node := value.(neo4j.Node)
				name := node.Props["name"]
				if val, err := record.Get("model"); err {
					modelNode := val.(neo4j.Node)
					props := modelNode.Props
					wbm := WinByMethod{
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

	return people.(map[string]WinByMethod), nil
}

func (env Env) listPayouts(space string) (map[string]WinByMethod, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	payouts, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		result, err := tx.Run(`
			MATCH (player:Player)<--(payout:Payout)<--(s:Space {uuid: $uuid})
			RETURN player, payout
			`, map[string]interface{}{"uuid": space})

		if err != nil {
			return nil, err
		}

		var payoutMap = make(map[string]WinByMethod)
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("player"); ok {
				node := value.(neo4j.Node)
				name := node.Props["name"]
				if val, err := record.Get("payout"); err {
					modelNode := val.(neo4j.Node)
					props := modelNode.Props
					wbm := WinByMethod{
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

	return payouts.(map[string]WinByMethod), nil
}

func (env Env) submitModel(name string, suuid string, json Model) (string, error) {
	session := env.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	records, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		result, err := tx.Run(`
			MATCH (player:Player {name: $name})-->(circle)-->(space:Space {uuid: $suuid})
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
			"name":  name,
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

	for _, record := range records.([]*neo4j.Record) {
		temp := record.Values[0].(neo4j.Node)
		fmt.Println(temp)
	}

	//fmt.Println("session", records)
	return "test text", nil
}

/*
a_by_dec: $adec,
a_by_ko: $ako,
b_by_dec: $bdec,
b_by_ko: $bko,
draw_nc: $draw
*/

/*
pub async fn post_model(name: String, suuid: String,
                        model: WinByMethod, graph: &Graph)
                        -> tide::Result<()> {

    let WinByMethod { a_by_dec, a_by_ko, b_by_dec, b_by_ko, draw_nc } = model;

    graph.run(
        Query::new("
            MATCH (player:Player {name: $name})-->(circle)-->(space:Space {uuid: $suuid})
            WITH player, space
            MERGE (player)-[:SETS]->(model:Model)-[:FOR]->(space) SET model = {
                a_by_dec: $adec,
                a_by_ko: $ako,
                b_by_dec: $bdec,
                b_by_ko: $bko,
                draw_nc: $draw
            }
            RETURN model
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

func Write() {
	/*
		records, err := session.WriteTransaction(
			func(tx neo4j.Transaction) (interface{}, error) {
				createRelationshipBetweenPeopleQuery := `
					MERGE (p1:Person { name: $person1_name })
					MERGE (p2:Person { name: $person2_name })
					MERGE (p1)-[:KNOWS]->(p2)
					RETURN p1, p2`
				result, err := tx.Run(createRelationshipBetweenPeopleQuery, map[string]interface{}{
					"person1_name": "Alice",
					"person2_name": "David",
				})
				if err != nil {
					return nil, err
				}

				return result.Collect() // Collects and commits
			})

		if err != nil {
			panic(err)
		}

		for _, record := range records.([]*neo4j.Record) {
			firstPerson := record.Values[0].(neo4j.Node)
			fmt.Printf("First: '%s'\n", firstPerson.Props["name"].(string))
			secondPerson := record.Values[1].(neo4j.Node)
			fmt.Printf("Second: '%s'\n", secondPerson.Props["name"].(string))
		}
	*/
}

func Read() {
	/*
		// Now read the created persons. By using ReadTransaction method a connection
		// to a read replica can be used which reduces load on writer nodes in cluster.
		_, err = session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			// Code within this function might be invoked more than once in case of
			// transient errors.
			readPersonByName := `
				MATCH (p:Person)
				WHERE p.name = $person_name
				RETURN p.name AS name`
			result, err := tx.Run(readPersonByName, map[string]interface{}{
				"person_name": "Alice",
			})
			if err != nil {
				return nil, err
			}
			// Iterate over the result within the transaction instead of using
			// Collect (just to show how it looks...). Result.Next returns true
			// while a record could be retrieved, in case of error result.Err()
			// will return the error.
			for result.Next() {
				fmt.Printf("Person name: '%s' \n", result.Record().Values[0].(string))
			}
			// Again, return any error back to driver to indicate rollback and
			// retry in case of transient error.
			return nil, result.Err()
		})
		if err != nil {
			panic(err)
		}
	*/
}
