package main

import (
	"errors"
	"fmt"
	"net/http"
)

type Downloader interface {
	Download(doi, url string) (string, error)
}

func DownloadServiceHandler(svc Downloader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		doi := r.URL.Query().Get("doi")
		uri := r.URL.Query().Get("url")

		if doi == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("No DOI parameter provided"))
			return
		}

		if uri == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("No URL parameter provided"))
			return
		}

		downloadLocation, err := svc.Download(doi, uri)
		if err != nil {
			var badRequest ErrorBadInput
			if errors.As(err, &badRequest) {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			fmt.Fprintf(w, "%v", err)
			return
		}

		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Location", downloadLocation)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(downloadLocation))
	})
}
