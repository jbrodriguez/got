package model

type Period int

const (
	Unknown Period = iota
	Today
	Someday
	Week
	Month
	Calendar
)
