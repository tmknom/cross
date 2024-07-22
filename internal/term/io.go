package term

import (
	"bufio"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type IO struct {
	InReader  io.Reader
	OutWriter io.Writer
	ErrWriter io.Writer
}

func (i *IO) IsPipe() bool {
	f, ok := i.InReader.(*os.File)
	if !ok {
		return false
	}
	return !terminal.IsTerminal(int(f.Fd()))
}

func (i *IO) Read() []string {
	lines := make([]string, 0, 64)
	scanner := bufio.NewScanner(i.InReader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
