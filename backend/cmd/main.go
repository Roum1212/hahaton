package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Roum1212/hahaton/cmd/hahaton"
)

func main() {
	if err := run(); err != nil {
		if _, writeError := fmt.Fprintln(os.Stderr, err); writeError != nil {
			os.Exit(1)
		}

		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return hahaton.Run(ctx)
}
