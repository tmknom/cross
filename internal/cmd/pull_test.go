package cmd

import (
	"testing"

	"github.com/tmknom/cross/internal/testlib"
)

func TestPullRunner_pull(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []struct {
		input *pullOptions
	}{
		{
			input: &pullOptions{
				branch: "main",
				dirs:   []string{"."},
				IO:     FakeTestIO(),
			},
		},
	}

	for _, tc := range cases {
		temp := t.TempDir()
		executeGit(t, ".", []string{"clone", "https://github.com/tmknom/cross.git", temp})
		sut := newPullRunner(tc.input)

		err := sut.pull(temp)
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}
	}
}
