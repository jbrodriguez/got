package model

import "time"

type Rules struct {
	DataDir  string
	Period   Period
	Interval Range
	Base     time.Time
}
