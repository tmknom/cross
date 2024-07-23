package cmd

import (
	"context"
	"strings"
	"testing"

	"github.com/tmknom/cross/internal/testlib"
)

func TestExecRunner_executeAll(t *testing.T) {
	cases := []struct {
		input    *execOptions
		expected string
	}{
		{
			input: &execOptions{
				command: "pwd",
				dirs:    []string{"."},
			},
			expected: "cross/internal/cmd",
		},
	}

	for _, tc := range cases {
		opts := tc.input
		opts.IO = FakeTestIO()
		sut := newExecRunner(opts)

		actual, err := sut.executeAll(context.Background())
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}

		if !strings.Contains(actual, tc.expected) {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
