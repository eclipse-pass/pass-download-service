package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// HTTP headers
const (
	headerUserAgent   = "User-Agent"
	headerContentType = "Content-Type"
	headerLocation    = "Location"
)

// InternalPassClient uses "private" backend URIs for interacting with the PASS repository
// It is intended for use on private networks.  Public URIs will be
// converted to private URIs when accessing the repository.
type InternalPassClient struct {
	Requester
	ExternalBaseURI string
	InternalBaseURI string
	Credentials     *Credentials
}

type Credentials struct {
	Username string
	Password string
}

// Requester performs http requests
type Requester interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c *InternalPassClient) PostBinary(url string, body io.Reader, contentType string) (string, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return "", errors.Wrapf(err, "could not build http request to %s", url)
	}

	if c.Credentials != nil {
		request.SetBasicAuth(c.Credentials.Username, c.Credentials.Password)
	}
	request.Header.Set(headerUserAgent, "pass-download-service")
	request.Header.Set(headerContentType, contentType)

	resp, err := c.Do(request)
	if err != nil {
		return "", errors.Wrapf(err, "error connecting to %s", url)
	}

	if resp.StatusCode > 299 {
		msg, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New("got error from Fedora: " + string(msg))
	}

	// Consume and discard the body
	defer resp.Body.Close()
	_, _ = ioutil.ReadAll(resp.Body)

	return c.translateToPublic(resp.Header.Get(headerLocation))
}

func (c *InternalPassClient) translateToPublic(uri string) (string, error) {
	if !strings.HasPrefix(uri, c.ExternalBaseURI) &&
		!strings.HasPrefix(uri, c.InternalBaseURI) {
		return uri, fmt.Errorf(`uri "%s" must start with internal or external baseuri"`, uri)
	}
	return strings.Replace(uri, c.InternalBaseURI, c.ExternalBaseURI, 1), nil
}

func mustSucceed(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return resp, err
	}

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		if resp.Body == nil {
			return nil, fmt.Errorf("request failed with code %d", resp.StatusCode)
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("request failed with code %d and message '%s'", resp.StatusCode, body)
	}

	return resp, err
}
