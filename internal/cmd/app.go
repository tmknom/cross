package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// AppName is the cli name (set by main.go)
var AppName string

// AppVersion is the current version (set by main.go)
var AppVersion string

type App struct {
	*IO
	rootCmd *cobra.Command
}

func NewApp(io *IO) *App {
	return &App{
		IO: io,
		rootCmd: &cobra.Command{
			Short: "Cross directory management tool",
		},
	}
}

func (a *App) Run(ctx context.Context, args []string) error {
	a.rootCmd.SetContext(ctx)

	// setup help message
	a.rootCmd.Use = AppName
	a.rootCmd.Version = AppVersion

	// override default settings
	a.rootCmd.SetArgs(args)
	a.rootCmd.SetIn(a.IO.InReader)
	a.rootCmd.SetOut(a.IO.OutWriter)
	a.rootCmd.SetErr(a.IO.ErrWriter)

	// setup log
	cobra.OnInitialize(func() { a.setupLog(args) })

	// setup version option
	a.rootCmd.SetVersionTemplate(AppVersion)

	// setup sub commands
	a.rootCmd.AddCommand(NewListCommand(a.IO))
	a.rootCmd.AddCommand(NewExecCommand(a.IO))
	a.rootCmd.AddCommand(NewCommitCommand(a.IO))
	a.rootCmd.AddCommand(NewPullCommand(a.IO))
	a.rootCmd.AddCommand(NewPushCommand(a.IO))

	return a.rootCmd.Execute()
}

func (a *App) setupLog(args []string) {
	log.SetOutput(io.Discard)
	if a.isDebug() {
		log.SetOutput(a.IO.ErrWriter)
	}
	log.SetPrefix(fmt.Sprintf("[%s] ", AppName))
	log.Printf("Start: args: %v", args)
}

func (a *App) isDebug() bool {
	switch os.Getenv("CROSS_DEBUG") {
	case "false", "0", "no":
		return false
	default:
		return true
	}
}

type IO struct {
	InReader  io.Reader
	OutWriter io.Writer
	ErrWriter io.Writer
}

func (i *IO) IsPipe() bool {
	f, ok := i.InReader.(*os.File)
	if !ok {
		return false
	}
	return !terminal.IsTerminal(int(f.Fd()))
}

func (i *IO) Read() []string {
	lines := make([]string, 0, 64)
	scanner := bufio.NewScanner(i.InReader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
