package week

import (
	"fmt"
	"os"

	// "text/tabwriter"
	"time"

	"github.com/Ladicle/tabwriter"
	. "github.com/logrusorgru/aurora/v3"

	"got/cli/lib"
	"got/cli/model"
)

type Header struct {
	rules model.Rules

	Total   [weekDays]time.Duration
	Working [weekDays]time.Duration
	Break   [weekDays]time.Duration

	days []string
}

func NewHeader(rules model.Rules) *Header {
	return &Header{
		rules: rules,
		days:  getDays(rules.Interval.Start.Local()),
	}
}

func (h *Header) Add(a *model.Activity) {
	day := lib.GetWeekday(a.End.Local())
	h.Total[day] += a.Duration
	h.Total[weekDays-1] += a.Duration
	if lib.IsWorking(a) {
		h.Working[day] += a.Duration
		h.Working[weekDays-1] += a.Duration
	} else {
		h.Break[day] += a.Duration
		h.Break[weekDays-1] += a.Duration
	}
}

func (h *Header) Render() {
	date := h.rules.Interval.Start

	_, week := date.ISOWeek()
	fmt.Printf("\n---------------------- (week %d) -----------------------\n", week)
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)

	fmt.Fprintf(w, " \t%s\t%s\t%s\t%s\t%s\t%s\t%s\t \t%s\t\n",
		h.days[0],
		h.days[1],
		h.days[2],
		h.days[3],
		h.days[4],
		h.days[5],
		h.days[6],
		h.days[7],
	)
	fmt.Fprintf(w, " \t \t \t \t \t \t \t \t \t \t\n")
	fmt.Fprintf(w, "Working:\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t \t%s\t\n",
		Green(lib.FormatDuration(h.Working[0])),
		Green(lib.FormatDuration(h.Working[1])),
		Green(lib.FormatDuration(h.Working[2])),
		Green(lib.FormatDuration(h.Working[3])),
		Green(lib.FormatDuration(h.Working[4])),
		Green(lib.FormatDuration(h.Working[5])),
		Green(lib.FormatDuration(h.Working[6])),
		Bold(Green(lib.FormatDuration(h.Working[7]))),
	)
	fmt.Fprintf(w, "Break:\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t \t%s\t\n",
		Gray(8-1, lib.FormatDuration(h.Break[0])),
		Gray(8-1, lib.FormatDuration(h.Break[1])),
		Gray(8-1, lib.FormatDuration(h.Break[2])),
		Gray(8-1, lib.FormatDuration(h.Break[3])),
		Gray(8-1, lib.FormatDuration(h.Break[4])),
		Gray(8-1, lib.FormatDuration(h.Break[5])),
		Gray(8-1, lib.FormatDuration(h.Break[6])),
		Gray(8-1, lib.FormatDuration(h.Break[7])),
	)
	fmt.Fprintf(w, " \t-----\t-----\t-----\t-----\t-----\t-----\t-----\t \t-----\t\n")
	fmt.Fprintf(w, "Total:\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t \t%s\t\n",
		lib.FormatDuration(h.Total[0]),
		lib.FormatDuration(h.Total[1]),
		lib.FormatDuration(h.Total[2]),
		lib.FormatDuration(h.Total[3]),
		lib.FormatDuration(h.Total[4]),
		lib.FormatDuration(h.Total[5]),
		lib.FormatDuration(h.Total[6]),
		lib.FormatDuration(h.Total[7]),
	)

	w.Flush()
}
