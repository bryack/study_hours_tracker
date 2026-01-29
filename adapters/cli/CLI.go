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
	in    io.Reader
}

func NewCLI(store store.SubjectStore, in io.Reader) *CLI {
	return &CLI{
		store: store,
		in:    in,
	}
}

func (cli *CLI) Run() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.store.RecordHour(extractSubjectAndHours(reader.Text()))
}

func extractSubjectAndHours(userInput string) (string, int) {
	str := strings.Split(userInput, " ")
	h, _ := strconv.Atoi(str[1])
	return str[0], h
}
