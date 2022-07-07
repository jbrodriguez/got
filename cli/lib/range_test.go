package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetYears(t *testing.T) {
	assert.Equal(t, []int{2022}, GetYears(time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)))
	assert.Equal(t, []int{2021, 2022}, GetYears(time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)))
	assert.Equal(t, []int{2020, 2021, 2022}, GetYears(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)))
}

func TestMonthInterval(t *testing.T) {
	bom, eom := monthInterval(2022, 6)
	assert.Equal(t, time.Date(2022, 6, 1, 0, 0, 0, 0, time.Local), bom.Local())
	assert.Equal(t, time.Date(2022, 6, 30, 23, 59, 59, 999999999, time.Local), eom.Local())

	_, bw := bom.Local().ISOWeek()
	_, ew := eom.Local().ISOWeek()
	assert.Equal(t, 22, bw)
	assert.Equal(t, 26, ew)
}

func TestCalendarInterval(t *testing.T) {
	bom, eom := monthInterval(2022, 7)
	assert.Equal(t, time.Date(2022, 7, 1, 0, 0, 0, 0, time.Local), bom.Local())
	assert.Equal(t, time.Date(2022, 7, 31, 23, 59, 59, 999999999, time.Local), eom.Local())

	start, _ := weekInterval(bom.Local().ISOWeek())
	_, end := weekInterval(eom.Local().ISOWeek())
	assert.Equal(t, time.Date(2022, 6, 27, 0, 0, 0, 0, time.Local), start.Local())
	assert.Equal(t, time.Date(2022, 7, 31, 23, 59, 59, 999999999, time.Local), end.Local())
}
