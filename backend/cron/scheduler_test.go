package cron

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var monday0am = time.Date(2022, 11, 21, 0, 0, 0, 0, time.UTC)

func TestShouldRun_GoodDay_WrongHour_FirstTime_NotRun(t *testing.T) {
	assert.False(t, NewSimpleScheduler().ShouldRun(1, 1, monday0am, time.Time{}))
}

func TestShouldRun_GoodDay_GoodHour_FirstTime_Run(t *testing.T) {
	assert.True(t, NewSimpleScheduler().ShouldRun(1, 0, monday0am, time.Time{}))
}

func TestShouldRun_GoodDay_WrongHour_SecondTime_NotRun(t *testing.T) {
	assert.False(t, NewSimpleScheduler().ShouldRun(1, 1, monday0am, monday0am))
}

func TestShouldRun_GoodDay_GoodHour_SecondTime_NotRun(t *testing.T) {
	assert.True(t, NewSimpleScheduler().ShouldRun(1, 1, monday0am.Add(1*time.Hour), monday0am))
}

func TestShouldRun_WrongDay_GoodHour_FirstTime_NotRun(t *testing.T) {
	assert.False(t, NewSimpleScheduler().ShouldRun(1, 1, monday0am.Add(25*time.Hour), time.Time{}))
}

func TestShouldRun_WrongDay_GoodHour_SecondTime_NotRun(t *testing.T) {
	assert.False(t, NewSimpleScheduler().ShouldRun(1, 1, monday0am.Add(25*time.Hour), monday0am))
}

func TestShouldRun_EveryDay_WrongHour_FirstTime_NotRun(t *testing.T) {
	assert.False(t, NewSimpleScheduler().ShouldRun(0, 1, monday0am, time.Time{}))
}

func TestShouldRun_EveryDay_GoodHour_SecondTime_NotRun(t *testing.T) {
	assert.False(t, NewSimpleScheduler().ShouldRun(0, 0, monday0am, monday0am))
}

func TestShouldRun_EveryDay_GoodHour_FirstTime_Run(t *testing.T) {
	assert.True(t, NewSimpleScheduler().ShouldRun(0, 0, monday0am, time.Time{}))
}

func TestShouldRun_Sunday_GoodHour_FirstTime_Run(t *testing.T) {
	assert.True(t, NewSimpleScheduler().ShouldRun(7, 0, monday0am.AddDate(0, 0, 6), time.Time{}))
}
