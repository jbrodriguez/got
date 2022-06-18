package today

import (
	"got/cli/model"
	"time"
)

type Today struct {
	rules    model.Rules
	sections []model.Section
	last     *model.Activity
}

func NewToday(rules model.Rules) *Today {
	return &Today{
		rules: rules,
		sections: []model.Section{
			NewHeader(rules),
			NewTags(),
			NewDetails(rules),
		},
	}
}

func (t *Today) Add(a *model.Activity) {
	t.last = a
	for _, s := range t.sections {
		s.Add(a)
	}
}

func (t *Today) Render() {
	now := time.Now()

	if t.last == nil {
		end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).UTC()
		t.last = &model.Activity{
			Start: end,
			End:   end,
		}
	}

	current := &model.Activity{
		Tag:      "current",
		Start:    t.last.End,
		End:      now.UTC(),
		Duration: now.Sub(t.last.End),
	}

	t.Add(current)

	for _, s := range t.sections {
		s.Render()
	}
}
