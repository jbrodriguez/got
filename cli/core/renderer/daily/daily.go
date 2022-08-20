package daily

import (
	"got/cli/model"
)

const weekDays = 8

type Daily struct {
	rules    model.Rules
	sections []model.Section
}

func NewDaily(rules model.Rules) *Daily {
	return &Daily{
		rules: rules,
		sections: []model.Section{
			NewDetails(rules),
		},
	}
}

func (d *Daily) Add(a *model.Activity) {
	for _, s := range d.sections {
		s.Add(a)
	}
}

func (d *Daily) Render() {
	for _, s := range d.sections {
		s.Render()
	}
}
