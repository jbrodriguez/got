package someday

import (
	"got/cli/model"
)

type Someday struct {
	rules    model.Rules
	sections []model.Section
	last     *model.Activity
}

func NewSomeday(rules model.Rules) *Someday {
	return &Someday{
		rules: rules,
		sections: []model.Section{
			NewHeader(rules),
			NewTags(),
			NewDetails(rules),
		},
	}
}

func (t *Someday) Add(a *model.Activity) {
	t.last = a
	for _, s := range t.sections {
		s.Add(a)
	}
}

func (t *Someday) Render() {
	for _, s := range t.sections {
		s.Render()
	}
}
