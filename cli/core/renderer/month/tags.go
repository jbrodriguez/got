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

type Tags struct {
	rules    model.Rules
	sections []model.Section

	weeksPerMonth int
	zeroIndex     int
	weeks         []string

	projects map[string][]time.Duration
	totals   []time.Duration
	working  []time.Duration
	breaks   []time.Duration
}

func NewTags(rules model.Rules) *Tags {
	zeroIndex, weeksPerMonth, weeks := getWeeks(rules)

	fmt.Printf("zero(%d), wpm(%d), weeks(%+v)\n", zeroIndex, weeksPerMonth, weeks)

	return &Tags{
		rules:         rules,
		weeksPerMonth: weeksPerMonth,
		zeroIndex:     zeroIndex,
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

	fmt.Printf("\n------------------- %s to %s (month %d) -------------------------------\t\n",
		date.Format("Jan 02"),
		t.rules.Interval.End.Local().Format("Jan 02"),
		date.Month(),
	)
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)

	line1 := " "
	line2 := ""
	count := 0
	values := make([]any, t.weeksPerMonth)
	for i := 0; i < t.weeksPerMonth-1; i++ {
		line1 += "\t%s"
		line2 += " \t"
		count += 1
		values[i] = t.weeks[i]
	}
	line1 += " \t%s\t\n"
	line2 += "\n"
	values[t.weeksPerMonth-1] = t.weeks[t.weeksPerMonth-1]

	fmt.Fprintf(w, line1, values...)
	fmt.Fprintf(w, line2)

	for _, tag := range keys {
		if tag == "break" {
			continue
		}

		project := t.projects[tag]

		line1 := "%s\t"
		line2 := ""
		count := 0
		values := make([]any, t.weeksPerMonth+1)
		values[0] = Blue(tag)
		for i := 1; i < t.weeksPerMonth; i++ {
			line1 += "%s\t"
			line2 += " \t"
			count += 1
			values[i] = getView(project[i-1])
		}
		line1 += " \t%s\t\n"
		line2 += "\n"
		values[t.weeksPerMonth] = getView(project[t.weeksPerMonth-1])

		fmt.Fprintf(w, line1, values...)
	}

	fmt.Fprint(w, line2)

	line1 = "%s\t"
	line2 = ""
	count = 0
	values = make([]any, t.weeksPerMonth+1)
	values[0] = Blue("working")
	for i := 1; i < t.weeksPerMonth; i++ {
		line1 += "%s\t"
		line2 += " \t"
		count += 1
		values[i] = getView(t.working[i-1])
	}
	line1 += " \t%s\t\n"
	line2 += "\n"
	values[t.weeksPerMonth] = getView(t.working[t.weeksPerMonth-1])

	fmt.Fprintf(w, line1, values...)

	line1 = "%s\t"
	line2 = ""
	count = 0
	values = make([]any, t.weeksPerMonth+1)
	values[0] = Blue("break")
	for i := 1; i < t.weeksPerMonth; i++ {
		line1 += "%s\t"
		line2 += " \t"
		count += 1
		values[i] = Gray(8-1, lib.FormatDuration(t.breaks[i-1]))
	}
	line1 += " \t%s\t\n"
	line2 += "\n"
	values[t.weeksPerMonth] = Gray(8-1, lib.FormatDuration(t.breaks[t.weeksPerMonth-1]))

	fmt.Fprintf(w, line1, values...)

	w.Flush()
}

func getWeeks(rules model.Rules) (int, int, []string) {
	_, firstWeek := rules.Interval.Start.Local().ISOWeek()
	_, lastWeek := rules.Interval.End.Local().ISOWeek()

	first := lib.GetRange(model.Week, rules.Interval.Start.Local())

	// one is for the difference
	// the other is to accommodate totals
	weeksPerMonth := (lastWeek - firstWeek) + 2
	weeks := make([]string, weeksPerMonth)
	current := first.Start
	for i := 0; i < weeksPerMonth; i++ {
		weeks[i] = current.Format("Jan 02")
		current = current.AddDate(0, 0, 7)
	}
	weeks[weeksPerMonth-1] = "Total"

	return firstWeek, weeksPerMonth, weeks
}

func (t *Tags) getWeekIndex(date time.Time) int {
	_, week := date.ISOWeek()
	return week - t.zeroIndex
}

func getView(duration time.Duration) interface{} {
	if duration <= 0 {
		return Gray(8-1, lib.FormatDuration(duration))
	}

	return Green(lib.FormatDuration(duration))
}
