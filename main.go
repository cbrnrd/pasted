package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cbrnrd/pasted/pkg/backends"
	"github.com/cbrnrd/pasted/pkg/config"
	"github.com/cbrnrd/pasted/pkg/transforms"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "pasted",
		Usage: "A simple pastebin application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Path to configuration file",
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
			// spew.Dump(config)
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

	go startPasteListener(backend, cfg, transformerChain)

	startWebServer(backend, cfg, transformerChain)
}

func startWebServer(backend backends.Backend, cfg *config.CLIConfig, chain *transforms.ChainTransformer) {
	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(httprate.LimitByIP(10, 1*time.Minute))

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
	http.ListenAndServe(cfg.HttpListenAddr, router)
}

func startPasteListener(backend backends.Backend, cfg *config.CLIConfig, chain *transforms.ChainTransformer) {
	l, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go handlePaste(conn, cfg, backend, chain)
	}
}

func handlePaste(conn net.Conn, cfg *config.CLIConfig, backend backends.Backend, chain *transforms.ChainTransformer) {
	defer conn.Close()

	transformed, err := chain.Transform(conn)
	if err != nil {
		fmt.Println("Error during transformation:", err)
		return
	}

	path, err := backend.Put(transformed)
	if err != nil {
		io.WriteString(conn, "Error storing paste: "+err.Error())
		return
	}

	path = cfg.Domain + "/" + path

	io.WriteString(conn, path)
}
