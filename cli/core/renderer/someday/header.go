package someday

import (
	"fmt"
	"os"
	"time"

	"github.com/Ladicle/tabwriter"
	. "github.com/logrusorgru/aurora/v3"

	"got/cli/lib"
	"got/cli/model"
)

type Header struct {
	Total   time.Duration
	Working time.Duration
	Break   time.Duration

	rules model.Rules
}

func NewHeader(rules model.Rules) *Header {
	return &Header{
		rules: rules,
	}
}

func (h *Header) Add(a *model.Activity) {
	if lib.IsOn(a) {
		return
	}

	h.Total += a.Duration
	if lib.IsWorking(a) {
		h.Working += a.Duration
	} else {
		h.Break += a.Duration
	}
}

func (h *Header) Render() {
	date := time.Now()
	if h.rules.Period == model.Someday {
		date = h.rules.Interval.Start
	}

	_, week := date.ISOWeek()
	fmt.Printf("\n---------------------- %s (week %d) -----------------------\n", date.Format("Monday, Jan 02, 2006"), week)
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)

	fmt.Fprintf(w, "working:\t%s\t\n", Bold(Green(lib.FormatDuration(h.Working))))
	fmt.Fprintf(w, "break:\t%s\t\n", Gray(8-1, lib.FormatDuration(h.Break)))
	fmt.Fprintf(w, " \t-----\t\n")
	fmt.Fprintf(w, "total:\t%s\t\n", lib.FormatDuration(h.Total))

	w.Flush()
}
