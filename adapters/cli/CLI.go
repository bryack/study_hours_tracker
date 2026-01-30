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

var (
	ErrNotEnoughArgs = errors.New("should be 2 arguments")
	ErrInvalidHours  = errors.New("failed to parse hours")
)

type CLI struct {
	store store.SubjectStore
	in    *bufio.Scanner
}

func NewCLI(store store.SubjectStore, in io.Reader) *CLI {
	return &CLI{
		store: store,
		in:    bufio.NewScanner(in),
	}
}

func (cli *CLI) Run() error {
	cli.in.Scan()
	s, h, err := extractSubjectAndHours(cli.in.Text())
	if err != nil {
		return fmt.Errorf("failed to extract subject and hours: %w", err)
	}
	cli.store.RecordHour(s, h)
	return nil
}

func extractSubjectAndHours(userInput string) (string, int, error) {
	str := strings.Split(userInput, " ")
	if len(str) < 2 {
		return "", 0, fmt.Errorf("failed to parse: %w, got: %d", ErrNotEnoughArgs, len(str))
	}
	h, err := strconv.Atoi(str[1])
	if err != nil {
		return "", 0, fmt.Errorf("%w %v: %v", ErrInvalidHours, str[1], err)
	}
	if h <= 0 {
		return "", 0, fmt.Errorf("%w %d, should be 1 or more: %v", ErrInvalidHours, h, err)
	}
	return str[0], h, nil
}
