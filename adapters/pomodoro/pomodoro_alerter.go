// Package pomodoro provides infrastructure implementations for timer functionality.
package pomodoro

import (
	"fmt"
	"io"
	"time"
)

type Alerter struct {
	ScheduleFunc func(duration time.Duration, message string, out io.Writer)
	WaitFunc     func(duration time.Duration)
}

func (a Alerter) ScheduleAlert(duration time.Duration, message string, out io.Writer) {
	a.ScheduleFunc(duration, message, out)
}

func (a Alerter) Wait(duration time.Duration) {
	a.WaitFunc(duration)
}

func RealScheduleAlert(duration time.Duration, message string, out io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintln(out, message)
	})
}

func RealWait(duration time.Duration) {
	time.Sleep(duration)
}
