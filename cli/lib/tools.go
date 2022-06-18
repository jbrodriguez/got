package lib

import (
	"fmt"
	"log"
	"os"
	"time"

	"got/cli/model"
)

// IsEmpty checks if the string is empty.
func IsEmpty(s string) bool {
	return len(s) == 0
}

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

func AppendToFile(file, message string) error {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString(message); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func IsCurrent(a *model.Activity) bool {
	return a.Tag == "current"
}

func IsOn(a *model.Activity) bool {
	return a.Tag == "on"
}

func IsWorking(a *model.Activity) bool {
	return a.Tag != "break"
}

func GetWeekday(date time.Time) int {
	day := date.Weekday()
	if int(day) == 0 {
		return int(model.Sunday)
	}

	return int(day) - 1
}
