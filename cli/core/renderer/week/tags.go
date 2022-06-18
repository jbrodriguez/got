package week

import (
	"fmt"
	"os"
	"sort"

	"time"

	"github.com/Ladicle/tabwriter"
	. "github.com/logrusorgru/aurora/v3"

	"got/cli/lib"
	"got/cli/model"
)

type ReportWeek [weekDays]time.Duration

type Tags struct {
	rules model.Rules

	projects map[string]*ReportWeek
	days     []string
	totals   *ReportWeek
	working  *ReportWeek
	breaks   *ReportWeek
}

func NewTags(rules model.Rules) *Tags {
	return &Tags{
		rules:    rules,
		days:     getDays(rules.Interval.Start.Local()),
		projects: make(map[string]*ReportWeek),
		totals:   &ReportWeek{},
		working:  &ReportWeek{},
		breaks:   &ReportWeek{},
	}
}

func (t *Tags) Add(a *model.Activity) {
	if lib.IsOn(a) {
		return
	}

	day := lib.GetWeekday(a.End.Local())

	if a.Tag == "break" {
		t.breaks[day] += a.Duration
		t.breaks[weekDays-1] += a.Duration

		t.totals[day] += a.Duration
		t.totals[weekDays-1] += a.Duration

		return
	}

	if _, ok := t.projects[a.Tag]; !ok {
		t.projects[a.Tag] = &ReportWeek{}
	}

	t.projects[a.Tag][day] += a.Duration
	t.projects[a.Tag][weekDays-1] += a.Duration

	t.working[day] += a.Duration
	t.working[weekDays-1] += a.Duration

	t.totals[day] += a.Duration
	t.totals[weekDays-1] += a.Duration
}

func (t *Tags) Render() {
	keys := make([]string, 0, len(t.projects))
	for k := range t.projects {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i int, j int) bool {
		return t.projects[keys[i]][weekDays-1] > t.projects[keys[j]][weekDays-1]
	})

	date := t.rules.Interval.Start.Local()
	_, week := date.ISOWeek()

	fmt.Printf("\n------------------- %s to %s (week %d) -------------------------------\t\n",
		date.Format("Jan 02"),
		t.rules.Interval.End.Local().Format("Jan 02"),
		week,
	)
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)

	fmt.Fprintf(w, " \t%s\t%s\t%s\t%s\t%s\t%s\t%s\t \t%s\t\n",
		t.days[0],
		t.days[1],
		t.days[2],
		t.days[3],
		t.days[4],
		t.days[5],
		t.days[6],
		t.days[7],
	)
	fmt.Fprintf(w, " \t \t \t \t \t \t \t \t \t \t\n")

	for _, tag := range keys {
		if tag == "break" {
			continue
		}

		project := t.projects[tag]
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t \t%s\t\n",
			Blue(tag),
			getView(project[0]),
			getView(project[1]),
			getView(project[2]),
			getView(project[3]),
			getView(project[4]),
			getView(project[5]),
			getView(project[6]),
			getView(project[7]),
		)
	}

	fmt.Fprintf(w, " \t \t \t \t \t \t \t \t \t \t\n")

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t \t%s\t\n",
		"working",
		getView(t.working[0]),
		getView(t.working[1]),
		getView(t.working[2]),
		getView(t.working[3]),
		getView(t.working[4]),
		getView(t.working[5]),
		getView(t.working[6]),
		Bold(Green(lib.FormatDuration(t.working[7]))),
	)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t \t%s\t\n",
		"break",
		Gray(8-1, lib.FormatDuration(t.breaks[0])),
		Gray(8-1, lib.FormatDuration(t.breaks[1])),
		Gray(8-1, lib.FormatDuration(t.breaks[2])),
		Gray(8-1, lib.FormatDuration(t.breaks[3])),
		Gray(8-1, lib.FormatDuration(t.breaks[4])),
		Gray(8-1, lib.FormatDuration(t.breaks[5])),
		Gray(8-1, lib.FormatDuration(t.breaks[6])),
		Gray(8-1, lib.FormatDuration(t.breaks[7])),
	)

	fmt.Fprintf(w, " \t-----\t-----\t-----\t-----\t-----\t-----\t-----\t \t-----\t\n")

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t \t%s\t\n",
		"total",
		lib.FormatDuration(t.totals[0]),
		lib.FormatDuration(t.totals[1]),
		lib.FormatDuration(t.totals[2]),
		lib.FormatDuration(t.totals[3]),
		lib.FormatDuration(t.totals[4]),
		lib.FormatDuration(t.totals[5]),
		lib.FormatDuration(t.totals[6]),
		lib.FormatDuration(t.totals[7]),
	)

	w.Flush()
}

func getView(duration time.Duration) interface{} {
	if duration <= 0 {
		return Gray(8-1, lib.FormatDuration(duration))
	}

	return Green(lib.FormatDuration(duration))
}
