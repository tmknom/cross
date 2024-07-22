package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmknom/cross/internal/cmd"
	"github.com/tmknom/cross/internal/term"
)

// Specify explicitly in ldflags
// For full details, see Makefile and .goreleaser.yml
var (
	name    = ""
	version = ""
	commit  = ""
	date    = ""
	url     = ""
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	ctx := context.Background()
	io := &term.IO{
		InReader:  os.Stdin,
		OutWriter: os.Stdout,
		ErrWriter: os.Stderr,
	}
	cmd.AppName = name
	cmd.AppVersion = fmt.Sprintf("%s version %s (%s:%s)\n%s\n", name, version, commit, date, url)
	return cmd.NewApp(io).Run(ctx, os.Args[1:])
}
