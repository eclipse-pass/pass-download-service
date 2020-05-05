package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// UnpaywallService looks up DOI info from unpaywall
type UnpaywallService struct {
	HTTP    Requester // Http client for interacting with unpaywall API
	Email   string    // Email for unpaywall requests
	Baseuri string    // Unpaywall baseURI
	Cache   *DoiCache // Can be nil if no caching is desired
}

// DOI lookup response from unpaywall
type unpaywallDOIResponse struct {
	BestOaLocation unpaywallLocation `json:"best_oa_location"`
}

type unpaywallLocation struct {
	URLForPdf             string `json:"url_for_pdf"`
	Version               string `json:"version"`
	RepositoryInstitution string `json:"repository_institution"`
}

// Lookup looks up DOI info for a given DOI
func (u UnpaywallService) Lookup(doi string) (*DoiInfo, error) {

	generator := func() (*DoiInfo, error) {

		results, err := u.get(u.apiRequestURI(doi))

		if err != nil {
			return nil, fmt.Errorf("unpaywall API request failed: %w", err)
		}

		var doiResponse DoiInfo

		// For now we'll only return the best location for the manuscript
		location := results.BestOaLocation

		// Get the file name from the decoded url for pdf
		// but log any problems do not cause response to fail
		var fileName string
		decodedURLForPdf, err := url.QueryUnescape(location.URLForPdf)
		if err != nil {
			log.Printf("file name decoding failed: %s", err)
		} else {
			splitURLForPdf := strings.Split(decodedURLForPdf, "/")
			fileName = splitURLForPdf[len(splitURLForPdf)-1]
		}

		if location.URLForPdf != "" {
			doiResponse.Manuscripts = append(doiResponse.Manuscripts, Manuscript{
				Location:              location.URLForPdf,
				RepositoryInstitution: location.RepositoryInstitution,
				Type:                  "application/pdf",
				Source:                "Unpaywall",
				Name:                  fileName,
			})
		}

		return &doiResponse, nil
	}

	if u.Cache != nil {
		return u.Cache.GetOrAdd(doi, generator)
	}

	return generator()
}

func (u UnpaywallService) apiRequestURI(doi string) string {
	return fmt.Sprintf("%s/%s?email=%s", u.Baseuri, doi, u.Email)
}

func (u UnpaywallService) get(uri string) (*unpaywallDOIResponse, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("could not form unpaywall API request: %w", err)
	}

	resp, err := mustSucceed(u.HTTP.Do(req))
	if err != nil {
		return nil, fmt.Errorf("unpaywall request failed: %w", err)
	}

	defer resp.Body.Close()

	var raw unpaywallDOIResponse
	return &raw, json.NewDecoder(resp.Body).Decode(&raw)
}
