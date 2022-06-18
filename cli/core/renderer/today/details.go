package today

import (
	"fmt"
	"os"

	"github.com/Ladicle/tabwriter"
	. "github.com/logrusorgru/aurora/v3"

	"got/cli/lib"
	"got/cli/model"
)

type Details struct {
	w *tabwriter.Writer

	rules   model.Rules
	current *model.Activity
}

func NewDetails(rules model.Rules) *Details {
	return &Details{
		w:     tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0),
		rules: rules,
	}
}

func (d *Details) Add(a *model.Activity) {
	if lib.IsOn(a) {
		return
	}

	if lib.IsCurrent(a) {
		d.current = a
		return
	}

	fmt.Fprintf(
		d.w,
		"%s\t%s - %s\t%s\t%s\t\n",
		Green(lib.FormatDuration(a.Duration)),
		a.Start.Local().Format("15:04"),
		a.End.Local().Format("15:04"),
		Blue(a.Tag),
		a.Task,
	)
}

func (d *Details) Render() {
	fmt.Println("\n---------------------- Details -----------------------")
	fmt.Println()

	fmt.Fprintf(
		d.w,
		"%s\t%s - %s\t%s\t%s\t\n",
		Green(lib.FormatDuration(d.current.Duration)),
		d.current.Start.Local().Format("15:04"),
		d.current.End.Local().Format("15:04"),
		Blue(d.current.Tag),
		d.current.Task,
	)

	d.w.Flush()
}
