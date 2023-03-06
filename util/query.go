package util

import "fmt"

const (
	PostWBM = `
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
    `
)

func carry() {
	fmt.Println("format")
}
