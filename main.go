package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/cbrnrd/pasted/pkg/backends"
	"github.com/cbrnrd/pasted/pkg/config"
	"github.com/cbrnrd/pasted/pkg/transforms"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "pasted",
		Usage: "A simple pastebin application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Usage: "Path to configuration file",
				Sources: cli.EnvVars("PASTED_CONFIG"),
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			configPath := c.String("config")
			if configPath == "" {
				return fmt.Errorf("config file is required")
			}

			configFile, err := os.Open(configPath)
			if err != nil {
				return fmt.Errorf("could not open config file: %v", err)
			}
			defer configFile.Close()

			var config config.CLIConfig
			if err := yaml.NewDecoder(configFile).Decode(&config); err != nil {
				return fmt.Errorf("could not parse config file: %v", err)
			}
			spew.Dump(config)
			startListeners(&config)
			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func startListeners(cfg *config.CLIConfig) {

	backend, err := cfg.GetBackend()
	if err != nil {
		panic(err)
	}

	tfs, err := cfg.GetTransforms()
	if err != nil {
		panic(err)
	}

	transformerChain := transforms.NewChainTransformer(
		tfs...,
	)

	go startPasteListener(backend, transformerChain)

	startWebServer(backend, transformerChain)
}

func startWebServer(backend backends.Backend, chain *transforms.ChainTransformer) {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/{key}", func(w http.ResponseWriter, r *http.Request) {
		pr, pw := io.Pipe()
		errCh := make(chan error, 1)

		go func() {
			defer pw.Close()
			key := chi.URLParam(r, "key")
			err := backend.Get(key, pw)
			if err != nil {
				errCh <- err
				return
			}
			errCh <- nil
		}()

		reversed, transformErr := chain.ReverseTransform(pr)

		// Wait for the goroutine's error, if any
		goroutineErr := <-errCh

		if goroutineErr != nil || transformErr != nil {
			// Close the pipe to ensure no further writes
			pr.Close()
			http.Error(w, "Error retrieving paste", http.StatusInternalServerError)
			return
		}

		_, err := io.Copy(w, reversed)
		if err != nil {
			http.Error(w, "Error writing paste", http.StatusInternalServerError)
			return
		}

	})
	http.ListenAndServe(":8080", router)
}

func startPasteListener(backend backends.Backend, chain *transforms.ChainTransformer) {
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go handlePaste(conn, backend, chain)
	}
}

func handlePaste(conn net.Conn, backend backends.Backend, chain *transforms.ChainTransformer) {
	defer conn.Close()

	transformed, err := chain.Transform(conn)
	if err != nil {
		fmt.Println("Error during transformation:", err)
		return
	}

	path, err := backend.Put(transformed)
	if err != nil {
		io.WriteString(conn, "Error storing paste: "+err.Error())
		panic(err)
	}

	io.WriteString(conn, path)
}
