package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tmknom/cross/internal/git"
)

func NewCommitCommand(io *IO) *cobra.Command {
	opts := &commitOptions{
		IO: io,
	}
	runner := newCommitRunner(opts)
	command := &cobra.Command{
		Use:   "commit",
		Short: "Commit all",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.run(cmd.Context()) },
	}
	command.PersistentFlags().StringVarP(&opts.branch, "branch", "b", "", "The branch name")
	command.PersistentFlags().StringVarP(&opts.message, "message", "m", "", "The commit message")
	command.PersistentFlags().StringSliceVarP(&opts.dirs, "dir", "d", []string{}, "The target directories")
	return command
}

type CommitRunner struct {
	opts *commitOptions
}

func newCommitRunner(opts *commitOptions) *CommitRunner {
	return &CommitRunner{
		opts: opts,
	}
}

type commitOptions struct {
	branch  string
	message string
	dirs    []string
	*IO
}

func (r *CommitRunner) run(ctx context.Context) error {
	log.Printf("Runner opts: %#v", r.opts)

	var result string
	for _, d := range r.opts.dirs {
		log.Printf("Commit: %s", d)
		commitHash, err := r.commitAll(d)
		if err != nil {
			return err
		}
		log.Printf("Latest commit hash: %s", commitHash)
	}
	_, err := fmt.Fprintf(r.opts.OutWriter, "%s", result)
	return err
}

func (r *CommitRunner) commitAll(path string) (string, error) {
	repo, err := git.OpenRepo(path)
	if err != nil {
		return "", err
	}

	err = repo.SwitchOrCreate(r.opts.branch)
	if err != nil {
		return "", err
	}

	err = repo.Add()
	if err != nil {
		return "", err
	}

	hash, err := repo.Commit(r.opts.message)
	if err != nil {
		return "", err
	}
	return hash, nil
}
