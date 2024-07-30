package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tmknom/cross/internal/git"
	"github.com/tmknom/cross/internal/term"
)

func NewPullCommand(io *term.IO) *cobra.Command {
	opts := &pullOptions{
		IO: io,
	}
	runner := newPullRunner(opts)
	command := &cobra.Command{
		Use:   "pull",
		Short: "Pull all",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.run(cmd.Context()) },
	}
	command.PersistentFlags().StringVarP(&opts.branch, "branch", "b", "main", "The branch name")
	command.PersistentFlags().StringSliceVarP(&opts.dirs, "dir", "d", []string{}, "The target directories")
	return command
}

type PullRunner struct {
	opts *pullOptions
}

func newPullRunner(opts *pullOptions) *PullRunner {
	return &PullRunner{
		opts: opts,
	}
}

type pullOptions struct {
	branch string
	dirs   []string
	*term.IO
}

func (r *PullRunner) run(ctx context.Context) error {
	log.Printf("Runner opts: %#v", r.opts)

	var result string
	for _, d := range r.opts.dirs {
		log.Printf("Pull: %s", d)
		err := r.pull(d)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(r.opts.OutWriter, "%s", result)
	return err
}

func (r *PullRunner) pull(path string) error {
	repo, err := git.OpenRepo(path)
	if err != nil {
		return err
	}

	err = repo.Switch(r.opts.branch)
	if err != nil {
		return err
	}

	return repo.Pull()
}
