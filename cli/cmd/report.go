package cmd

import (
	"time"

	"github.com/pkg/errors"

	"got/cli/core/renderer/month"
	"got/cli/core/renderer/someday"
	"got/cli/core/renderer/today"
	"got/cli/core/renderer/week"
	"got/cli/core/reporter"
	"got/cli/lib"
	"got/cli/model"
)

type Report struct {
	Today   bool   `short:"t" help:"report for the current day" `
	Someday bool   `short:"s" help:"report for a specific day"`
	Week    bool   `short:"w" help:"report for the current week"`
	Month   bool   `short:"m" help:"report for the current month"`
	Param   string `arg:"" optional:"" help:"date to report on"`
}

func (r *Report) Run(ctx *Context) error {
	period := model.Unknown
	switch {
	case r.Today || (!r.Today && !r.Someday && !r.Week && !r.Month && r.Param == ""):
		period = model.Today
	case (r.Someday && r.Param != "") || (r.Param != "" && !r.Someday && !r.Week && !r.Month):
		period = model.Someday
	case r.Week:
		period = model.Week
	case r.Month:
		period = model.Month
	}

	if period == model.Unknown {
		return errors.Errorf("incorrect arguments")
	}

	date := time.Now()
	if r.Param != "" {
		var err error
		date, err = time.Parse("2006-01-02", r.Param)
		if err != nil {
			return err
		}
	}

	interval := lib.GetRange(period, date)

	rules := model.Rules{
		Period:   period,
		Interval: interval,
		DataDir:  ctx.DataDir,
	}

	var renderer model.Renderer
	switch period {
	case model.Month:
		renderer = month.NewMonth(rules)
	case model.Week:
		renderer = week.NewWeek(rules)
	case model.Someday:
		renderer = someday.NewSomeday(rules)
	case model.Today:
		fallthrough
	default:
		renderer = today.NewToday(rules)

	}

	rpt := reporter.CreateReporter(renderer, rules)

	return rpt.Run()
}
