package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tmknom/cross/internal/testlib"
)

func TestListRunner_listGitDirs(t *testing.T) {
	cases := []struct {
		input    *listOptions
		expected []string
	}{
		{
			input: &listOptions{
				base:     "../../",
				excludes: []string{"tmp", ".makefiles"},
			},
			expected: []string{"."},
		},
	}

	for _, tc := range cases {
		opts := tc.input
		opts.IO = FakeTestIO()
		sut := newListRunner(opts)

		actual, err := sut.listGitDirs()
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}

		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}
