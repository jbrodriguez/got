package someday

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Ladicle/tabwriter"
	. "github.com/logrusorgru/aurora/v3"

	"got/cli/lib"
	"got/cli/model"
)

type Tag struct {
	Duration time.Duration
	Task     []string
}

type Tags struct {
	list map[string]*Tag
}

func NewTags() *Tags {
	return &Tags{
		list: make(map[string]*Tag),
	}
}

func (t *Tags) Add(a *model.Activity) {
	if lib.IsOn(a) {
		return
	}

	if _, ok := t.list[a.Tag]; !ok {
		t.list[a.Tag] = &Tag{
			Duration: a.Duration,
			Task:     []string{a.Task},
		}
		return
	}

	t.list[a.Tag].Duration += a.Duration
	t.list[a.Tag].Task = append(t.list[a.Tag].Task, a.Task)
}

func (t *Tags) Render() {
	keys := make([]string, 0, len(t.list))
	for k := range t.list {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i int, j int) bool {
		return t.list[keys[i]].Duration > t.list[keys[j]].Duration
	})

	fmt.Println("\n---------------------- Tags -----------------------")
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for i, tag := range keys {
		tk := t.list[keys[i]]
		fmt.Fprintf(w, "%s\t%s\t%s\t\n",
			Blue(tag),
			Green(lib.FormatDuration(tk.Duration)),
			strings.Join(tk.Task, ", "),
		)
	}
	w.Flush()
}
