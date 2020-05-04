package main

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type DownloadService struct {
	HTTP   Requester     // Http client for downloading content
	Fedora BinaryStore   // PASS/Fedora client
	Dest   string        // URL of Fedora container where binaries will be deposited into
	DOIs   LookupService // DOI lookup service (for verifying validity of download URI for a given DOI)
}

// Binarystore is a place where binary content can be POSTed.  If successful, the URL of the
// newly-stored content will be returned.
type BinaryStore interface {
	PostBinary(url string, body io.Reader, contentType string) (string, error)
}

// Download verifies that the given url is valid for a given DOI, downloads it into Fedora,
// Then returns the resulting URL of the binary.  Note:  It does *not* create a File entity.
func (d DownloadService) Download(doi, url string) (string, error) {
	info, err := d.DOIs.Lookup(doi)
	if err != nil {
		return "", errors.Wrapf(err, "could not lookup doi %s", doi)
	}

	if err = d.verifyURL(doi, info, url); err != nil {
		return "", errors.Wrapf(err, "could not validate url %s for doi %s", url, doi)
	}

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := d.HTTP.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "could not fetch content URL")
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", errors.Errorf("download of '%s' faied with %d %s", url, resp.StatusCode, string(body))
	}

	return d.Fedora.PostBinary(d.Dest, resp.Body, resp.Header.Get(headerContentType))

}

func (d DownloadService) verifyURL(doi string, info *DoiInfo, url string) error {
	for _, m := range info.Manuscripts {
		if m.Location == url {
			return nil // We found the matching URL.  Done!
		}
	}

	return ErrorBadInput("no matching URL found for DOI")
}
