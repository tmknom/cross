package cmd

import (
	"bytes"
	"context"
	"github.com/tmknom/cross/internal/term"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			input:    []string{""},
			expected: "Cross directory management tool\n\n",
		},
	}

	for _, tc := range cases {
		sut := NewApp(FakeTestIO())

		err := sut.Run(context.Background(), tc.input)
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}

		actual := sut.IO.OutWriter.(*bytes.Buffer).String()
		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}

func FakeTestIO() *term.IO {
	return &term.IO{
		InReader:  &bytes.Buffer{},
		OutWriter: &bytes.Buffer{},
		ErrWriter: os.Stderr,
	}
}
