package reporter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"got/cli/lib"
	"got/cli/model"
)

var (
	// emptyStr    string
	whitespaces = "\t\n\r "
	// oneDay      = 24 * time.Hour
)

type Reporter struct {
	renderer model.Renderer
	rules    model.Rules
}

func CreateReporter(renderer model.Renderer, rules model.Rules) *Reporter {
	return &Reporter{
		renderer: renderer,
		rules:    rules,
	}
}

func (r *Reporter) Run() error {
	activities, err := r.GetActivities(r.rules.DataDir, r.rules.Interval)
	if err != nil {
		return err
	}

	r.PrintReport(activities)

	return nil
}

func (r *Reporter) GetActivities(dataDir string, interval model.Range) ([]*model.Activity, error) {
	years := lib.GetYears(interval.Start, interval.End)

	list := []*model.Activity{}
	activityID := 1

	for _, year := range years {
		f := filepath.Join(dataDir, fmt.Sprintf("%d.log", year))
		file, err := os.Open(f)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		prev := &model.Activity{}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := strings.Trim(scanner.Text(), whitespaces) // Read line

			// Ignore blank or comment lines
			if lib.IsEmpty(text) || strings.HasPrefix(text, "#") {
				continue
			}

			activity, err := lib.ParseActivity(text)
			if err != nil {
				return nil, err
			}

			if activity.End.Before(interval.Start) || activity.End.After(interval.End) {
				continue
			}

			activity.ID = activityID
			if activity.Tag == "on" {
				activity.Start = activity.End
			} else {
				activity.Start = prev.End
				activity.Duration = activity.End.Sub(activity.Start)
			}

			list = append(list, activity)

			activityID++
			prev = activity
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (r *Reporter) PrintReport(activities []*model.Activity) {
	for _, activity := range activities {
		r.renderer.Add(activity)
	}

	r.renderer.Render()
}
