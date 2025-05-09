package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tmknom/cross/internal/dir"
	"github.com/tmknom/cross/internal/errlib"
)

func NewListCommand(io *IO) *cobra.Command {
	opts := &listOptions{
		IO: io,
	}
	runner := newListRunner(opts)
	command := &cobra.Command{
		Use:   "list",
		Short: "List directories under Git version control",
		RunE:  func(cmd *cobra.Command, args []string) error { return runner.run(cmd.Context()) },
	}
	command.PersistentFlags().StringVarP(&opts.base, "base", "b", ".", "The base directory that contains repositories")
	command.PersistentFlags().StringVarP(&opts.format, "format", "f", "default", "The output format (default, csv)")
	command.PersistentFlags().StringSliceVarP(&opts.excludes, "exclude", "e", []string{}, "The exclude directories")
	return command
}

type ListRunner struct {
	opts *listOptions
}

func newListRunner(opts *listOptions) *ListRunner {
	return &ListRunner{
		opts: opts,
	}
}

type listOptions struct {
	base     string
	format   string
	excludes []string
	*IO
}

func (r *ListRunner) run(ctx context.Context) error {
	log.Printf("Runner opts: %#v", r.opts)
	dirs, err := r.listGitDirs()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(r.opts.OutWriter, r.output(dirs))

	return err
}

func (r *ListRunner) listGitDirs() ([]string, error) {
	result := make([]string, 0, 64)
	base := dir.NewBaseDir(r.opts.base).Abs()
	log.Printf("Walk base dir: %s", base)

	err := filepath.WalkDir(base, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return errlib.Wrapf(err, "invalid base: %s", base)
		}
		if !entry.IsDir() {
			return nil
		}
		for _, exclude := range r.opts.excludes {
			if strings.Contains(path, exclude) {
				return nil
			}
		}
		if entry.Name() == ".git" {
			result = append(result, filepath.Clean(filepath.Dir(path)))
		}
		return nil
	})

	sort.Strings(result)
	log.Printf("Find dirs: %v", result)
	return result, err
}

func (r *ListRunner) output(dirs []string) string {
	switch r.opts.format {
	case "csv":
		return fmt.Sprintf(strings.Join(dirs, ","))
	default:
		return fmt.Sprintf(strings.Join(dirs, "\n"))
	}
}
