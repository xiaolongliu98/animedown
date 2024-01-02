package ratecalc

import (
	"fmt"
	"time"
)

type TimeDot struct {
	Time  time.Time
	Value int64 // total bytes
}

type Calculator struct {
	TimeWindow    int       // seconds
	ValueWindow   []TimeDot // total bytes slice
	DeltaValueSum int64     //
}

func NewCalculator(timeWindow ...int) *Calculator {
	if len(timeWindow) == 0 {
		timeWindow = []int{3}
	}
	return &Calculator{
		TimeWindow:  timeWindow[0],
		ValueWindow: make([]TimeDot, 0),
	}
}

func (c *Calculator) Add(value int64) {
	if len(c.ValueWindow) > 0 {
		c.DeltaValueSum += value - c.ValueWindow[len(c.ValueWindow)-1].Value
	}

	c.ValueWindow = append(c.ValueWindow, TimeDot{time.Now(), value})
	c.removeOld()
}

func (c *Calculator) removeOld() {
	// remove old values
	for len(c.ValueWindow) > 0 && c.ValueWindow[0].Time.Before(time.Now().Add(-time.Duration(c.TimeWindow)*time.Second)) {
		if len(c.ValueWindow) > 1 {
			c.DeltaValueSum -= c.ValueWindow[1].Value - c.ValueWindow[0].Value
		}
		c.ValueWindow = c.ValueWindow[1:]
	}
}

func (c *Calculator) GetAverage() int64 {
	c.removeOld()
	if len(c.ValueWindow) < 2 {
		return 0
	}
	return c.DeltaValueSum / int64(len(c.ValueWindow)-1)
}

// Reset resets the calculator.
func (c *Calculator) Reset() {
	c.ValueWindow = make([]TimeDot, 0)
	c.DeltaValueSum = 0
}

func (c *Calculator) GetAverageKiB() float64 {
	return float64(c.GetAverage()) / 1024
}

func (c *Calculator) GetAverageMiB() float64 {
	return float64(c.GetAverage()) / 1024 / 1024
}

func (c *Calculator) GetAverageAuto() string {
	avg := c.GetAverage()
	if avg < 1024 {
		return fmt.Sprintf("%.2fB/s", float64(avg))
	} else if avg < 1024*1024 {
		return fmt.Sprintf("%.2fKiB/s", float64(avg)/1024)
	} else {
		return fmt.Sprintf("%.2fMiB/s", float64(avg)/1024/1024)
	}
}
