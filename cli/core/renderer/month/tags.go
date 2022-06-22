package month

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

type Fx func(time.Duration) any

type Tags struct {
	rules    model.Rules
	sections []model.Section

	weeksPerMonth int
	firstWeek     int
	weeks         []string

	projects map[string][]time.Duration
	totals   []time.Duration
	working  []time.Duration
	breaks   []time.Duration
}

func NewTags(rules model.Rules) *Tags {
	firstWeek, weeksPerMonth, weeks := getWeeks(rules)
	return &Tags{
		rules:         rules,
		weeksPerMonth: weeksPerMonth,
		firstWeek:     firstWeek,
		weeks:         weeks,
		projects:      make(map[string][]time.Duration),
		totals:        make([]time.Duration, weeksPerMonth),
		working:       make([]time.Duration, weeksPerMonth),
		breaks:        make([]time.Duration, weeksPerMonth),
	}
}

func (t *Tags) Add(a *model.Activity) {
	if lib.IsOn(a) {
		return
	}

	week := t.getWeekIndex(a.End.Local())

	if a.Tag == "break" {
		t.breaks[week] += a.Duration
		t.breaks[t.weeksPerMonth-1] += a.Duration

		t.totals[week] += a.Duration
		t.totals[t.weeksPerMonth-1] += a.Duration

		return
	}

	if _, ok := t.projects[a.Tag]; !ok {
		t.projects[a.Tag] = make([]time.Duration, t.weeksPerMonth)
	}

	t.projects[a.Tag][week] += a.Duration
	t.projects[a.Tag][t.weeksPerMonth-1] += a.Duration

	t.working[week] += a.Duration
	t.working[t.weeksPerMonth-1] += a.Duration

	t.totals[week] += a.Duration
	t.totals[t.weeksPerMonth-1] += a.Duration
}

func (t *Tags) Render() {
	keys := make([]string, 0, len(t.projects))
	for k := range t.projects {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i int, j int) bool {
		return t.projects[keys[i]][t.weeksPerMonth-1] > t.projects[keys[j]][t.weeksPerMonth-1]
	})

	date := t.rules.Interval.Start.Local()

	fmt.Printf("\n----------------------------- %s %s -----------------------------\t\n",
		date.Format("January"),
		date.Format("2006"),
	)
	fmt.Println()

	separator := t.createSeparator()
	splitter := t.createSplitter()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)

	t.renderWeekNames(w)

	fmt.Fprintf(w, separator)

	t.renderProjects(w, keys)

	fmt.Fprintf(w, separator)

	t.renderWorking(w)
	t.renderBreak(w)

	fmt.Fprintf(w, splitter)

	t.renderTotal(w)

	w.Flush()
}

func (t *Tags) getWeekIndex(date time.Time) int {
	_, week := date.ISOWeek()
	return week - t.firstWeek
}

func (t *Tags) renderWeekNames(w *tabwriter.Writer) {
	line := " "
	values := make([]any, t.weeksPerMonth)
	for i := 0; i < t.weeksPerMonth-1; i++ {
		line += "\t%s"
		values[i] = t.weeks[i]
	}
	line += "\t\t%s\t\n"
	values[t.weeksPerMonth-1] = t.weeks[t.weeksPerMonth-1]

	fmt.Fprintf(w, line, values...)
}

func (t *Tags) renderProjects(w *tabwriter.Writer, keys []string) {
	for _, tag := range keys {
		if tag == "break" {
			continue
		}

		project := t.projects[tag]

		t.renderProject(w, tag, project)
	}
}

func (t *Tags) renderProject(w *tabwriter.Writer, tag string, project []time.Duration) {
	line, values := renderLine(tag, project, t.weeksPerMonth, formatProjects, nil, true)
	fmt.Fprintf(w, line, values...)
}

func (t *Tags) renderWorking(w *tabwriter.Writer) {
	line, values := renderLine("working", t.working, t.weeksPerMonth, formatWorking, formatTotalWorking, false)
	fmt.Fprintf(w, line, values...)
}

func (t *Tags) renderBreak(w *tabwriter.Writer) {
	line, values := renderLine("break", t.breaks, t.weeksPerMonth, formatBreaks, nil, false)
	fmt.Fprintf(w, line, values...)
}

func (t *Tags) renderTotal(w *tabwriter.Writer) {
	line, values := renderLine("total", t.totals, t.weeksPerMonth, formatTotals, nil, false)
	fmt.Fprintf(w, line, values...)
}

func (t *Tags) createSeparator() string {
	line := ""

	for i := 0; i < t.weeksPerMonth+2; i++ {
		line += " \t"
	}
	line += "\n"

	return line
}

func (t *Tags) createSplitter() string {
	// fmt.Fprintf(w, " \t-----\t-----\t-----\t-----\t-----\t \t-----\t\n")
	line := " "

	for i := 0; i < t.weeksPerMonth-1; i++ {
		line += "\t-----"
	}
	line += "\t \t------\t\n"

	return line
}

func getWeeks(rules model.Rules) (int, int, []string) {
	_, firstWeek := rules.Interval.Start.Local().ISOWeek()
	_, lastWeek := rules.Interval.End.Local().ISOWeek()

	// +1 for the difference
	// +1 to accommodate totals
	weeksPerMonth := (lastWeek - firstWeek) + 2
	weeks := make([]string, weeksPerMonth)

	current := rules.Interval.Start

	for i := 0; i < weeksPerMonth; i++ {
		weeks[i] = current.Format("Jan 02")
		current = current.AddDate(0, 0, 7)
	}
	weeks[weeksPerMonth-1] = "Total"

	return firstWeek, weeksPerMonth, weeks
}

func renderLine(tag string, samples []time.Duration, weeksPerMonth int, fn, fnLast Fx, emphasis bool) (string, []any) {
	line := "%s\t"
	values := make([]any, weeksPerMonth+1)
	values[0] = tag
	if emphasis {
		values[0] = Blue(tag)
	}
	for i := 1; i < weeksPerMonth; i++ {
		line += "%s\t"
		values[i] = fn(samples[i-1])
	}
	line += " \t%s\t\n"
	values[weeksPerMonth] = fn(samples[weeksPerMonth-1])
	if fnLast != nil {
		values[weeksPerMonth] = fnLast(samples[weeksPerMonth-1])
	}

	return line, values
}

func formatProjects(duration time.Duration) any {
	if duration <= 0 {
		return Gray(8-1, lib.FormatDuration(duration))
	}

	return Green(lib.FormatDuration(duration))
}

func formatWorking(duration time.Duration) any {
	return formatProjects(duration)
}

func formatBreaks(duration time.Duration) any {
	return Gray(8-1, lib.FormatDuration(duration))
}

func formatTotalWorking(duration time.Duration) any {
	return Bold(Green(lib.FormatDuration(duration)))
}

func formatTotals(duration time.Duration) any {
	return lib.FormatDuration(duration)
}
