package someday

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

	rules model.Rules
}

func NewDetails(rules model.Rules) *Details {
	return &Details{
		w:     tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0),
		rules: rules,
	}
}

func (h *Details) Add(a *model.Activity) {
	if lib.IsOn(a) {
		return
	}

	fmt.Fprintf(
		h.w,
		"%s\t%s - %s\t%s\t%s\t\n",
		Green(lib.FormatDuration(a.Duration)),
		a.Start.Local().Format("15:04"),
		a.End.Local().Format("15:04"),
		Blue(a.Tag),
		a.Task,
	)
}

func (h *Details) Render() {
	fmt.Println("\n---------------------- Details -----------------------")
	fmt.Println()

	h.w.Flush()
}
