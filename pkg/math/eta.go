package math

import "time"

type ETA struct {
	lastTime time.Time
	lastVal  float64
}

func NewETA() *ETA {
	return &ETA{
		lastTime: time.Now(),
		lastVal:  0,
	}
}

func (e *ETA) UpdatePercents(progress100 float64) time.Duration {
	return e.Update(progress100 / 100)
}

func (e *ETA) Update(progress01 float64) time.Duration {
	now := time.Now()
	since := now.Sub(e.lastTime)
	e.lastTime = now

	diff := progress01 - e.lastVal
	e.lastVal = progress01

	remain := 1 - progress01

	left := float64(since) * remain / diff

	dur := time.Duration(left)
	dur -= dur % time.Millisecond

	if dur < 0 {
		return 0
	}
	return dur
}
