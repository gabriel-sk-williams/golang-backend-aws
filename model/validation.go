package model

import (
	"goyave.dev/goyave/v4/validation"
)

// SubmitModel()
var (
	SubmissionProps = validation.RuleSet{
		"puuid": validation.List{"required", "string"},
		"suuid": validation.List{"required", "string"},
		"model": validation.List{"required", "object"},
	}
)

var (
	PlayerCircleProps = validation.RuleSet{
		"puuid": validation.List{"required", "string"},
		"cuuid": validation.List{"required", "string"},
	}
)

var (
	PlayerSpaceProps = validation.RuleSet{
		"puuid": validation.List{"required", "string"},
		"suuid": validation.List{"required", "string"},
	}
)

var (
	CircleProps = validation.RuleSet{
		"cuuid": validation.List{"required", "string"},
	}
)

var (
	SpaceProps = validation.RuleSet{
		"uuid":    validation.List{"required", "string"},
		"fields":  validation.List{"required", "array:string"},
		"pattern": validation.List{"required", "string"},
		"stake":   validation.List{"required", "numeric"},
	}
)
