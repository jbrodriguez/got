package month

import "got/cli/model"

type Month struct {
	rules    model.Rules
	sections []model.Section
}

func NewMonth(rules model.Rules) *Month {
	return &Month{
		rules: rules,
		sections: []model.Section{
			NewTags(rules),
		},
	}
}

func (m *Month) Add(a *model.Activity) {
	for _, s := range m.sections {
		s.Add(a)
	}
}

func (m *Month) Render() {
	for _, s := range m.sections {
		s.Render()
	}
}
