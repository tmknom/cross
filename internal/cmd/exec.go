package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tmknom/cross/internal/errlib"
	"github.com/tmknom/cross/internal/term"
)

func NewExecCommand(io *term.IO) *cobra.Command {
	opts := &execOptions{
		IO: io,
	}
	runner := newExecRunner(opts)
	command := &cobra.Command{
		Use:   "exec",
		Short: "Execute a command",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.run(cmd.Context()) },
	}
	command.PersistentFlags().StringVarP(&opts.command, "command", "c", "", "The shell command to execute")
	command.PersistentFlags().StringSliceVarP(&opts.dirs, "dir", "d", []string{}, "The target directories")
	return command
}

type ExecRunner struct {
	opts *execOptions
}

func newExecRunner(opts *execOptions) *ExecRunner {
	return &ExecRunner{
		opts: opts,
	}
}

type execOptions struct {
	command string
	dirs    []string
	*term.IO
}

func (r *ExecRunner) run(ctx context.Context) error {
	log.Printf("Runner opts: %#v", r.opts)
	stdout, err := r.executeAll(ctx)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(r.opts.OutWriter, "%s", stdout)
	return err
}

func (r *ExecRunner) executeAll(ctx context.Context) (string, error) {
	name, args, err := r.parseCommand()
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for _, workdir := range r.opts.dirs {
		log.Printf("execute: %s: %s", r.opts.command, workdir)
		stdout, err := r.execute(ctx, name, args, workdir)
		if err != nil {
			return "", err
		}
		b.WriteString(fmt.Sprintf("%s\n", stdout))
		log.Printf("stdout:\n%s", stdout)
	}
	return b.String(), nil
}

func (r *ExecRunner) parseCommand() (string, []string, error) {
	elements := strings.Split(r.opts.command, " ")
	if len(elements) == 0 {
		return "", nil, fmt.Errorf("no command to execute")
	}

	name := elements[0]
	if len(elements) == 1 {
		return name, nil, nil
	}
	return name, elements[1:], nil
}

func (r *ExecRunner) execute(ctx context.Context, name string, args []string, workdir string) (string, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = workdir
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%s\n", cmd.String()))
		b.WriteString(fmt.Sprintf("Stderr\n%s\n", stdout.String()))
		b.WriteString(fmt.Sprintf("Stdout\n%s\n", stderr.String()))
		b.WriteString(fmt.Sprintf("Workdir: %v\n", cmd.Dir))
		return "", errlib.Wrapf(err, "%s", b.String())
	}
	return stdout.String(), nil
}
