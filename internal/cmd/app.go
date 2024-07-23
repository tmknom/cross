package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tmknom/cross/internal/term"
)

// AppName is the cli name (set by main.go)
var AppName string

// AppVersion is the current version (set by main.go)
var AppVersion string

type App struct {
	*term.IO
	rootCmd *cobra.Command
}

func NewApp(io *term.IO) *App {
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
	case "true", "1", "yes":
		return true
	default:
		return false
	}
}
