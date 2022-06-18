package week

import (
	"got/cli/model"
	"time"
)

const weekDays = 8

type Week struct {
	rules    model.Rules
	sections []model.Section
}

func NewWeek(rules model.Rules) *Week {
	return &Week{
		rules: rules,
		sections: []model.Section{
			NewTags(rules),
		},
	}
}

func (w *Week) Add(a *model.Activity) {
	for _, s := range w.sections {
		s.Add(a)
	}
}

func (w *Week) Render() {
	for _, s := range w.sections {
		s.Render()
	}
}

func getDays(start time.Time) []string {
	days := make([]string, 8)
	for i := 0; i < 7; i++ {
		days[i] = start.AddDate(0, 0, i).Format("Jan 02")
	}
	days[7] = "Total"
	return days
}
