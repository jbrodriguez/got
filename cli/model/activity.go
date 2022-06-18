package model

import (
	"time"
)

type Activity struct {
	ID       int
	Tag      string
	Task     string
	Start    time.Time
	End      time.Time
	Duration time.Duration
}
