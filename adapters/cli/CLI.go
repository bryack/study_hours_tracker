package cli

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/bryack/study_hours_tracker/store"
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

func (cli *CLI) Run() {
	cli.in.Scan()
	cli.store.RecordHour(extractSubjectAndHours(cli.in.Text()))
}

func extractSubjectAndHours(userInput string) (string, int) {
	str := strings.Split(userInput, " ")
	h, _ := strconv.Atoi(str[1])
	return str[0], h
}
