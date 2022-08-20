package lib

import (
	"time"

	"got/cli/model"
)

func GetRange(period model.Period, date time.Time) model.Range {
	var start, end time.Time

	switch period {
	case model.Today:
		start = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).UTC()
		end = time.Now().UTC()
	case model.Someday:
		start = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local).UTC()
		end = time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, -1, time.Local).UTC()
	case model.Daily:
		fallthrough
	case model.Week:
		start, end = weekInterval(date.ISOWeek())
	case model.Month:
		start, end = monthInterval(date.Year(), date.Month())
	case model.Calendar:
		ms, me := monthInterval(date.Year(), date.Month())
		start, _ = weekInterval(ms.Local().ISOWeek())
		_, end = weekInterval(me.Local().ISOWeek())
	}

	return model.Range{
		Start: start,
		End:   end,
	}
}

// https://stackoverflow.com/questions/53076290/how-to-get-the-start-date-and-end-date-of-current-month-using-golang
func monthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.Local).UTC()
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.Local).UTC()
	return firstDay, lastDay
}

// https://stackoverflow.com/questions/52300644/date-range-by-week-number-golang
func weekInterval(year, week int) (start, end time.Time) {
	start = weekStart(year, week)
	future := start.AddDate(0, 0, 7)
	end = time.Date(future.Year(), future.Month(), future.Day(), 0, 0, 0, -1, time.Local)
	return start.UTC(), end.UTC()
}

func weekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.Local)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

func GetYears(start, end time.Time) []int {
	if start.Year() == end.Year() {
		return []int{start.Year()}
	}

	var years []int
	for year := start.Year(); year <= end.Year(); year++ {
		years = append(years, year)
	}
	return years
}

func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
