package cmd

import (
	"time"

	"github.com/pkg/errors"

	"got/cli/core/renderer/daily"
	"got/cli/core/renderer/month"
	"got/cli/core/renderer/someday"
	"got/cli/core/renderer/today"
	"got/cli/core/renderer/week"
	"got/cli/core/reporter"
	"got/cli/lib"
	"got/cli/model"
)

type Report struct {
	Today    bool   `short:"t" help:"report for the current day" `
	Someday  bool   `short:"s" help:"report for a specific day"`
	Week     bool   `short:"w" help:"report for a week"`
	Month    bool   `short:"m" help:"report for a month"`
	Calendar bool   `short:"c" help:"report for a calendar month"`
	Daily    bool   `short:"d" help:"report for each day of a week"`
	Param    string `arg:"" optional:"" help:"date to report on"`
}

func (r *Report) Run(ctx *Context) error {
	period := model.Unknown
	switch {
	case r.Today || (!r.Today && !r.Someday && !r.Week && !r.Month && !r.Calendar && !r.Daily && r.Param == ""):
		period = model.Today
	case (r.Someday && r.Param != "") || (r.Param != "" && !r.Someday && !r.Week && !r.Month && !r.Calendar && !r.Daily):
		period = model.Someday
	case r.Week:
		period = model.Week
	case r.Month:
		period = model.Month
	case r.Calendar:
		period = model.Calendar
	case r.Daily:
		period = model.Daily
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
		Base:     date,
	}

	var renderer model.Renderer
	switch period {
	case model.Calendar:
		renderer = month.NewMonth(rules)
	case model.Month:
		renderer = month.NewMonth(rules)
	case model.Week:
		renderer = week.NewWeek(rules)
	case model.Someday:
		renderer = someday.NewSomeday(rules)
	case model.Daily:
		renderer = daily.NewDaily(rules)
	case model.Today:
		fallthrough
	default:
		renderer = today.NewToday(rules)

	}

	rpt := reporter.CreateReporter(renderer, rules)

	return rpt.Run()
}
