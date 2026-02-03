package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bryack/study_hours_tracker/store"
)

const GretingString = "Let's study\nType {subject} {hours} to track hours\nOr type 'pomodoro' {subject} to use pomodoro tracker"

var (
	ErrNotEnoughArgs = errors.New("should be 2 arguments")
	ErrInvalidHours  = errors.New("failed to parse hours")
)

type PomodoroRunner interface {
	Start()
}

type CLI struct {
	store          store.SubjectStore
	in             *bufio.Scanner
	out            io.Writer
	pomodoroRunner PomodoroRunner
}

func NewCLI(store store.SubjectStore, in io.Reader, out io.Writer, pomodoroRunner PomodoroRunner) *CLI {
	return &CLI{
		store:          store,
		in:             bufio.NewScanner(in),
		out:            out,
		pomodoroRunner: pomodoroRunner,
	}
}

func (cli *CLI) Run() error {
	fmt.Fprintln(cli.out, GretingString)

	for cli.in.Scan() {
		s, h, isPomodoro, err := extractSubjectAndHours(cli.in.Text())
		if err != nil {
			fmt.Fprintf(cli.out, "failed to extract subject and hours: %v\n", err)
			continue
		}

		if isPomodoro {
			fmt.Fprintln(cli.out, "Pomodoro started...")
			cli.pomodoroRunner.Start()
		}
		cli.store.RecordHour(s, h)
	}

	if err := cli.in.Err(); err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	return nil
}

func extractSubjectAndHours(userInput string) (subject string, hours int, isPomodoro bool, err error) {
	args := strings.Split(userInput, " ")
	if len(args) < 2 {
		return "", 0, false, fmt.Errorf("failed to parse: %w, got: %d", ErrNotEnoughArgs, len(args))
	}

	if args[0] == "pomodoro" {
		return args[1], 1, true, nil
	}

	h, err := strconv.Atoi(args[1])
	if err != nil {
		return "", 0, false, fmt.Errorf("%w %v: %v", ErrInvalidHours, args[1], err)
	}
	if h <= 0 {
		return "", 0, false, fmt.Errorf("%w %d, should be 1 or more: %v", ErrInvalidHours, h, err)
	}
	return args[0], h, false, nil
}
