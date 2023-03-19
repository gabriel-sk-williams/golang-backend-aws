package calc

import (
	"math"
	"sort"
)

type Pair struct {
	Name string
	Cert float64
}

// PAYOUTS
// payouts with this method are calculated to each of person's certainty in the group
// using a "reverse waterfall" method

// receives hashmap of prediction models -> name: { outcome: certainty, ... }
// returns hashmap of payouts -> name: { outcome: payout, ... }
func Payouts(
	models map[string]map[string]interface{},
	fields []string,
	stake float64) (map[string]map[string]float64, error) {

	// outcome: { name: payout, ... }
	outcomeMap := make(map[string]map[string][]float64)
	for _, field := range fields {
		oca := outcomeArray(models, field)
		pomap := payoutMap(oca, stake)
		outcomeMap[field] = pomap //outcome_map.insert(String::from(oc), pomap);
	}

	// convert back to original form -> name: { outcome: payout, ... }
	payoutMap := make(map[string]map[string]float64)
	for name := range models {
		personalMap := make(map[string]float64)
		for _, oc := range fields {
			payout := aggregate(oc, name, outcomeMap)
			personalMap[oc] = payout
		}
		payoutMap[name] = personalMap
	}

	return payoutMap, nil
}

// generate a vec of tuples -> outcome: [ (name, certainty)... ]
// representing each person's prediction of each possible outcome
func outcomeArray(models map[string]map[string]interface{}, field string) []Pair {

	var oca []Pair
	for name, model := range models {
		pair := Pair{name, model[field].(float64)}
		oca = append(oca, pair)
	}

	sort.Slice(oca, func(i, j int) bool {
		return oca[i].Cert < oca[j].Cert
	})

	return oca
}

// receive the vec of tuples -> outcome: [ (name, certainty), ... ]
// return hashmap of payouts -> outcome: { name: payout, ... }
func payoutMap(oca []Pair, stake float64) map[string][]float64 {

	// create empty array for each name: [ payout0, payout1, ... ]
	var blankMap = func(oca []Pair) map[string][]float64 {

		pomap := make(map[string][]float64)
		for _, value := range oca {
			var container []float64
			pomap[value.Name] = container
		}
		return pomap
	}

	// calculate raw loss of bad prediction
	var rawLoss = func(certainty float64, stake float64) float64 {
		fraction := certainty / 100
		dim := (fraction * stake * 100) / 100 // previously used .round()
		return dim - stake
	}

	// moderate according to most correct person in field
	var modLoss = func(ocs []Pair) float64 {
		pop := ocs[len(ocs)-1]
		best := pop.Cert
		return best / 100
	}

	// sum the certainties of remaining field
	var sumCerts = func(rem []Pair) float64 {
		sum := 0.0
		for _, pair := range rem {
			sum += pair.Cert
		}
		return sum
	}

	var getConsecutive = func(data []Pair) int {

		alpha := data[0]
		var consecutive int

		for index, value := range data {
			if value.Cert == alpha.Cert {
				consecutive = index
				continue
			} else {
				break
			}
		}

		return consecutive + 1
	}

	pomap := blankMap(oca)

	for len(oca) > 0 {
		consecutive := getConsecutive(oca)
		trust := oca[0:consecutive] // copy consecutive elements
		oca = oca[consecutive:]     // remove consecutive
		if len(oca) == 0 {
			break
		}

		for _, current := range trust {
			rawLoss := rawLoss(current.Cert, stake)
			modLoss := modLoss(oca)
			paidLoss := rawLoss * modLoss
			pomap[current.Name] = append(pomap[current.Name], paidLoss)

			payout := math.Abs(paidLoss)
			mount := sumCerts(oca)

			for _, next := range oca {
				mass := next.Cert / mount
				portion := payout * mass
				pomap[next.Name] = append(pomap[next.Name], portion)
			}
		}
	}

	return pomap
}

// access vec name: [payouts, ...] and sum all values to get final payout
func aggregate(field string, name string, oculus map[string]map[string][]float64) float64 {

	ocmap := oculus[field]
	paySlice := ocmap[name]

	// round all values, sum, and round again for consistency
	// rounded := roundAll(paySlice)
	// let rounded: Vec<f64> = payvec.iter().map(|f| (f*100.0_f64).round() / 100.0_f64 ).collect();
	collapsed := sumPayouts(paySlice)
	//let final_round: f64 = (collapsed*100.0_f64).round() / 100.0_f64;

	return collapsed //finalRound
}

func roundAll(slice []float64) []float64 {
	sum := 0.0
	for _, value := range slice {
		sum += value
	}

	s := []float64{0.0, 1.2}
	return s
}

func sumPayouts(slice []float64) float64 {
	sum := 0.0
	for _, value := range slice {
		sum += value
	}
	return sum
}
