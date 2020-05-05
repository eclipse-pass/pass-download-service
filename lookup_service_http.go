package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// LookupService looks up a DOI and provides information associated with it
type LookupService interface {
	Lookup(doi string) (*DoiInfo, error)
}

func LookupServiceHandler(svc LookupService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		doi := r.URL.Query().Get("doi")

		if doi == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("No DOI parameter provided"))
			return
		}

		info, err := svc.Lookup(doi)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Content-Type", "application/json;charset=utf-8")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error encoding JSON response: %s", err)
		}
	})
}
