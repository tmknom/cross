package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/tmknom/cross/internal/testlib"
)

func TestCommitRunner_commitAll(t *testing.T) {
	cases := []struct {
		input    *commitOptions
		expected string
	}{
		{
			input: &commitOptions{
				branch:  "test-branch",
				message: "test commit message",
				dirs:    []string{"."},
				IO:      FakeTestIO(),
			},
			expected: "test commit message",
		},
	}

	for _, tc := range cases {
		repo := initRepo(t)
		sut := newCommitRunner(tc.input)

		hash, err := sut.commitAll(repo)
		if err != nil {
			t.Fatalf(testlib.FormatError(err, sut, tc.input))
		}

		actual := executeGit(t, repo, []string{"log", hash, "-1", "--oneline"})
		if !strings.Contains(actual, tc.expected) {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}

func initRepo(t *testing.T) string {
	temp := t.TempDir()
	_ = executeGit(t, temp, []string{"init"})
	_ = executeGit(t, temp, []string{"commit", "-m", "initial commit", "--allow-empty"})
	_ = os.WriteFile(temp+"/foo.txt", []byte("foo"), fs.ModePerm)
	return temp
}

func executeGit(t *testing.T, workDir string, args []string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = workDir
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	err := cmd.Run()

	if err != nil {
		logCommand(cmd)
		t.Fatalf(testlib.FormatError(err, nil, nil))
	}
	return cmd.Stdout.(*bytes.Buffer).String()
}

func logCommand(cmd *exec.Cmd) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s\n", cmd.String()))
	b.WriteString(fmt.Sprintf("Workdir: %v\n", cmd.Dir))
	b.WriteString(fmt.Sprintf("Stdout:\n%s", cmd.Stdout.(*bytes.Buffer).String()))
	b.WriteString(fmt.Sprintf("Stderr:\n%s", cmd.Stderr.(*bytes.Buffer).String()))
	log.Printf("%s\n", b.String())
}
