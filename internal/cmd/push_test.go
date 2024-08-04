package cmd

import (
	"testing"

	"github.com/tmknom/cross/internal/testlib"
)

func TestPushRunner_push(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []struct {
		input *pushOptions
	}{
		{
			input: &pushOptions{
				branch: "test/push",
				dirs:   []string{"."},
				IO:     FakeTestIO(),
			},
		},
	}

	temp := t.TempDir()
	executeGit(t, ".", []string{"clone", "git@github.com:tmknom/cross.git", temp})
	executeGit(t, temp, []string{"switch", "-c", "test/push"})
	defer executeGit(t, temp, []string{"push", "origin", "-d", "test/push"})

	for _, tc := range cases {
		sut := newPushRunner(tc.input)

		err := sut.push(temp)
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}
	}
}
