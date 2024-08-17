package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/connorvoisey/shgrid_api/pkg/load"
	"github.com/connorvoisey/shgrid_api/pkg/routes"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/humacli"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"3000"`
}

func Run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log, db, err := load.Init()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var api huma.API
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := http.NewServeMux()
		api = humago.New(router, huma.DefaultConfig("My API", "1.0.0"))

		routes.AddRoutes(api, log, db)
		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			log.Info().
				Int("Port", options.Port).
				Msg("Started server")
			err := http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
			if err != nil {
				log.Err(err).Msg("Failed to listen and serve")
				panic(err)
			}
		})
	})

	// Add a command to print the OpenAPI spec.
	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			b, _ := api.OpenAPI().MarshalJSON()
			fmt.Println(string(b))
		},
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()

	return nil
}
