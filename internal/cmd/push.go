package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tmknom/cross/internal/git"
	"github.com/tmknom/cross/internal/term"
)

func NewPushCommand(io *term.IO) *cobra.Command {
	opts := &pushOptions{
		IO: io,
	}
	runner := newPushRunner(opts)
	command := &cobra.Command{
		Use:   "push",
		Short: "Push all",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.run(cmd.Context()) },
	}
	command.PersistentFlags().StringVarP(&opts.branch, "branch", "b", "main", "The branch name")
	command.PersistentFlags().StringSliceVarP(&opts.dirs, "dir", "d", []string{}, "The target directories")
	return command
}

type PushRunner struct {
	opts *pushOptions
}

func newPushRunner(opts *pushOptions) *PushRunner {
	return &PushRunner{
		opts: opts,
	}
}

type pushOptions struct {
	branch string
	dirs   []string
	*term.IO
}

func (r *PushRunner) run(ctx context.Context) error {
	log.Printf("Runner opts: %#v", r.opts)

	var result string
	for _, d := range r.opts.dirs {
		log.Printf("Push: %s", d)
		err := r.push(d)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(r.opts.OutWriter, "%s", result)
	return err
}

func (r *PushRunner) push(path string) error {
	repo, err := git.OpenRepo(path)
	if err != nil {
		return err
	}

	return repo.Push(r.opts.branch)
}
