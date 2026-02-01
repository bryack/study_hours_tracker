package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/bryack/study_hours_tracker/store"
)

var (
	ErrNotEnoughArgs = errors.New("should be 2 arguments")
	ErrInvalidHours  = errors.New("failed to parse hours")
)

type Sleeper interface {
	Sleep(duration time.Duration)
}

type CLI struct {
	store   store.SubjectStore
	in      *bufio.Scanner
	out     io.Writer
	sleeper Sleeper
}

type PomodoroSleeper struct{}

func (ps *PomodoroSleeper) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

func NewCLI(store store.SubjectStore, in io.Reader, out io.Writer, sleeper Sleeper) *CLI {
	return &CLI{
		store:   store,
		in:      bufio.NewScanner(in),
		out:     out,
		sleeper: sleeper,
	}
}

func (cli *CLI) Run() error {
	cli.in.Scan()
	s, h, d, err := extractSubjectAndHours(cli.in.Text())
	if err != nil {
		return fmt.Errorf("failed to extract subject and hours: %w", err)
	}

	if d > 0 {
		fmt.Fprintln(cli.out, "Pomodoro started...")
		cli.sleeper.Sleep(d)
	}

	cli.store.RecordHour(s, h)
	return nil
}

func extractSubjectAndHours(userInput string) (string, int, time.Duration, error) {
	args := strings.Split(userInput, " ")
	if len(args) < 2 {
		return "", 0, 0, fmt.Errorf("failed to parse: %w, got: %d", ErrNotEnoughArgs, len(args))
	}

	if args[0] == "pomodoro" {
		return args[1], 1, 25 * time.Minute, nil
	}

	h, err := strconv.Atoi(args[1])
	if err != nil {
		return "", 0, 0, fmt.Errorf("%w %v: %v", ErrInvalidHours, args[1], err)
	}
	if h <= 0 {
		return "", 0, 0, fmt.Errorf("%w %d, should be 1 or more: %v", ErrInvalidHours, h, err)
	}
	return args[0], h, 0, nil
}
