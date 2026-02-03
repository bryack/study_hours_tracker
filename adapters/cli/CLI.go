// Package cli provides a command-line interface for the study hours tracker.
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

const (
	GreetingString  = "Let's study\nType {subject} {hours} to track hours\nOr type 'pomodoro' {subject} to use pomodoro tracker\nType 'quit' to exit"
	PomodoroCommand = "pomodoro"
	QuitCommand     = "quit"
)

var (
	ErrNotEnoughArgs = errors.New("should be 2 arguments")
	ErrInvalidHours  = errors.New("failed to parse hours")
)

// CLI provides an interactive command-line interface for tracking study hours.
type CLI struct {
	in      *bufio.Scanner
	out     io.Writer
	session *StudySession
}

// NewCLI creates a new CLI with the given dependencies.
func NewCLI(store store.SubjectStore, in io.Reader, out io.Writer, pomodoroRunner PomodoroRunner) *CLI {
	return &CLI{
		in:      bufio.NewScanner(in),
		out:     out,
		session: NewStudySession(store, pomodoroRunner),
	}
}

// Run starts the interactive CLI loop.
func (cli *CLI) Run() error {
	fmt.Fprintln(cli.out, GreetingString)

	for cli.in.Scan() {
		input := cli.in.Text()
		if input == QuitCommand {
			fmt.Fprintln(cli.out, "Goodbye!")
			break
		}
		s, h, isPomodoro, err := extractSubjectAndHours(cli.in.Text())
		if err != nil {
			fmt.Fprintf(cli.out, "failed to extract subject and hours: %v\n", err)
			continue
		}

		if isPomodoro {
			fmt.Fprintln(cli.out, "Pomodoro started...")
			if err := cli.session.RecordPomodoro(s); err != nil {
				fmt.Fprintf(cli.out, "failed to record pomodoro: %v\n", err)
			}
		} else {
			if err := cli.session.RecordManual(s, h); err != nil {
				fmt.Fprintf(cli.out, "failed to record hours: %v\n", err)
			}
		}
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

	if args[0] == PomodoroCommand {
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
