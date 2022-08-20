package lib

import (
	"regexp"
	"strings"
	"time"

	"got/cli/model"
)

var (
	emptyStr   string
	activityRx = regexp.MustCompile(`(?P<End>\d{4}-(?:0[1-9]|1[0-2])-(?:0[1-9]|[1-2]\d|3[0-1])T(?:[0-1]\d|2[0-3]):[0-5]\d:[0-5]\dZ)\s+(?P<Tag>.*)\:\s+(?P<Task>.*)`)
	onRx       = regexp.MustCompile(`(?P<End>\d{4}-(?:0[1-9]|1[0-2])-(?:0[1-9]|[1-2]\d|3[0-1])T(?:[0-1]\d|2[0-3]):[0-5]\d:[0-5]\dZ)\s+on$`)
)

func parseTime(s string) (time.Time, error) {
	return time.ParseInLocation(time.RFC3339, s, time.Local)
}

func ParseActivity(line string) (*model.Activity, error) {
	original := strings.Trim(line, " \t\n\r")

	activity := &model.Activity{}

	var err error
	var end time.Time
	var tag, task string

	if onRx.MatchString(original) {
		matches := onRx.FindStringSubmatch(original)
		index := onRx.SubexpIndex("End")

		end, err = parseTime(matches[index])
		if err != nil {
			return nil, err
		}

		tag = "on"
	} else if activityRx.MatchString(original) {
		matches := activityRx.FindStringSubmatch(original)

		end, err = parseTime(matches[activityRx.SubexpIndex("End")])
		if err != nil {
			return nil, err
		}

		tag = matches[activityRx.SubexpIndex("Tag")]
		task = matches[activityRx.SubexpIndex("Task")]
	}

	activity.End = end
	activity.Tag = tag
	activity.Task = task

	return activity, nil
}
