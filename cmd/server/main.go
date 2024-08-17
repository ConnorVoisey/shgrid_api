package main

import (
	"context"
	"fmt"
	"os"

	"github.com/connorvoisey/shgrid_api/pkg/server"
)

func main() {
	ctx := context.Background()
	if err := server.Run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
