package cmd

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/tmknom/cross/internal/testlib"
)

func TestApp_Run(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []struct {
		input    []string
		expected string
	}{
		{
			input:    []string{"list", "--base", "../../", "--exclude", "tmp,.makefiles"},
			expected: ".",
		},
		{
			input:    []string{"exec", "--dir", ".", "--command", "pwd"},
			expected: "cross/internal/cmd",
		},
	}

	for _, tc := range cases {
		sut := NewApp(FakeTestIO())

		err := sut.Run(context.Background(), tc.input)
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}

		actual := sut.IO.OutWriter.(*bytes.Buffer).String()
		if !strings.Contains(actual, tc.expected) {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}

func FakeTestIO() *IO {
	return &IO{
		InReader:  &bytes.Buffer{},
		OutWriter: &bytes.Buffer{},
		ErrWriter: os.Stderr,
	}
}
