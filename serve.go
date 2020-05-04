package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/urfave/cli/v2"
)

type serveOpts struct {
	port                int
	downloadDest        string
	unpaywallEmail      string
	unpaywallBaseURI    string
	publicFedoraBaseURI string
	fedoraBaseURI       string
	fedoraUsername      string
	fedoraPassword      string
}

func serve() *cli.Command {

	var opts serveOpts

	return &cli.Command{
		Name:  "serve",
		Usage: "Start the user service web service",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Port for serving http user service",
				Required:    false,
				Destination: &opts.port,
				EnvVars:     []string{"DOWNLOAD_SERVICE_PORT"},
				Value:       8091,
			},
			&cli.StringFlag{
				Name:        "download.dest",
				Usage:       "URI of Fedora container to deposit binaries into",
				Required:    false,
				Destination: &opts.downloadDest,
				EnvVars:     []string{"DOWNLOAD_SERVICE_DEST"},
			},

			&cli.StringFlag{
				Name:        "unpaywall.email",
				Usage:       "Email used for making unpaywall API requests",
				Required:    false,
				Destination: &opts.unpaywallEmail,
				EnvVars:     []string{"UNPAYWALL_REQUEST_EMAIL"},
			},
			&cli.StringFlag{
				Name:        "unpaywall.baseuri",
				Usage:       "Unpaywall API BaseURI",
				Required:    false,
				Destination: &opts.unpaywallBaseURI,
				EnvVars:     []string{"UNPAYWALL_BASEURI"},
			},
			&cli.StringFlag{
				Name:        "fedora.public.baseurl",
				Usage:       "External (public) PASS baseurl",
				Destination: &opts.publicFedoraBaseURI,
				EnvVars:     []string{"PASS_EXTERNAL_FEDORA_BASEURL"},
			},
			&cli.StringFlag{
				Name:        "fedora.internal.baseurl",
				Usage:       "Internal (private) PASS baseuri",
				Destination: &opts.fedoraBaseURI,
				EnvVars:     []string{"PASS_FEDORA_BASEURL"},
			},
			&cli.StringFlag{
				Name:        "fedora.username",
				Usage:       "Username for basic auth to Fedora",
				Destination: &opts.fedoraUsername,
				EnvVars:     []string{"PASS_FEDORA_USER"},
			},
			&cli.StringFlag{
				Name:        "password, p",
				Usage:       "Password for basic auth to Fedora",
				Destination: &opts.fedoraPassword,
				EnvVars:     []string{"PASS_FEDORA_PASSWORD"},
			},
		},
		Action: func(c *cli.Context) error {
			return serveAction(opts)
		},
	}
}

func serveAction(opts serveOpts) error {

	httpClient := &http.Client{
		Timeout: 20 * time.Second,
	}

	var fedoraCredentials *Credentials
	if opts.fedoraUsername != "" {
		fedoraCredentials = &Credentials{
			Username: opts.fedoraUsername,
			Password: opts.fedoraPassword,
		}
	}

	unpaywall := UnpaywallService{
		HTTP:    httpClient,
		Baseuri: opts.unpaywallBaseURI,
		Email:   opts.unpaywallEmail,
		Cache: NewDoiCache(DoiCacheConfig{
			MaxAge:  1 * time.Minute,
			MaxSize: 100,
		}),
	}

	downloadService := DownloadService{
		HTTP: httpClient,
		DOIs: unpaywall,
		Dest: opts.downloadDest,
		Fedora: &InternalPassClient{
			Requester:       httpClient,
			Credentials:     fedoraCredentials,
			ExternalBaseURI: opts.publicFedoraBaseURI,
			InternalBaseURI: opts.fedoraBaseURI,
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/lookup", LookupServiceHandler(unpaywall))
	mux.Handle("/download", DownloadServiceHandler(downloadService))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.port),
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	done := make(chan error, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Listening on port %d", opts.port)
		done <- server.ListenAndServe()
	}()

	select {
	case <-stop:
		_ = server.Shutdown(context.Background())
		log.Printf("Goodbye!")
		return nil
	case err := <-done:
		return err
	}
}
