package util

import "fmt"

const (
	PostWBM = `
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
    `
)

func carry() {
	fmt.Println("format")
}
