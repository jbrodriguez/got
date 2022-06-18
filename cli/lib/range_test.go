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
