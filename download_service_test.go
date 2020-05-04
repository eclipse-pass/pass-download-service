package main_test

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	pass "github.com/oa-pass/pass-download-service"
)

type MockRequester func(*http.Request) (*http.Response, error)

func (f MockRequester) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type MockBinaryStore func(string, io.Reader, string) (string, error)

func (f MockBinaryStore) PostBinary(url string, body io.Reader, contentType string) (string, error) {
	return f(url, body, contentType)
}

func TestInputErrors(t *testing.T) {
	doiDoesNotExist := "test/nonExistentDOI"
	doiHasNoURLs := "test/noURLs"

	toTest := pass.DownloadService{
		DOIs: MockLookupService(func(doi string) (*pass.DoiInfo, error) {
			if doi == doiDoesNotExist {
				return nil, pass.ErrorBadInput("noDOI")
			}

			if doi == doiHasNoURLs {
				return &pass.DoiInfo{}, nil
			}

			return nil, nil
		}),
	}

	cases := map[string]string{
		"doi does not exist": doiDoesNotExist,
		"bad URL for doi":    doiHasNoURLs,
	}

	for name, doi := range cases {
		doi := doi
		t.Run(name, func(t *testing.T) {
			_, err := toTest.Download(doi, "foo:/bar")

			var badInput pass.ErrorBadInput
			if !errors.As(err, &badInput) {
				t.Fatalf("expected a bad input error %s", err)
			}
		})
	}
}

func TestErrorFetch(t *testing.T) {

	// We'll use these as both URLs and DOIs for the sake of verbosity
	badConnect := "http://example.org/badConnect"
	badErrorCode := "http://example.org/errorCode"

	toTest := pass.DownloadService{
		DOIs: MockLookupService(func(doi string) (*pass.DoiInfo, error) {
			return &pass.DoiInfo{
				Manuscripts: []pass.Manuscript{{Location: doi}},
			}, nil
		}),
		HTTP: MockRequester(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == badConnect {
				return nil, errors.New("cannot connect")
			}

			if req.URL.String() == badErrorCode {
				return &http.Response{StatusCode: 404, Body: ioutil.NopCloser(strings.NewReader(""))}, nil
			}

			return nil, nil
		}),
	}

	cases := map[string]string{
		"connection error": badConnect,
		"bad error code":   badErrorCode,
	}

	for name, doi := range cases {
		doi := doi
		t.Run(name, func(t *testing.T) {
			_, err := toTest.Download(doi, doi)
			if err == nil {
				t.Fatal("Should have gotten an error")
			}
		})
	}
}

func TestPostBody(t *testing.T) {

	doi := "abc/123"
	location := "http://example.org/unpaywall/file.pdf"
	dest := "http://fcrepo:8080/fcrepo/rest/bin/"
	fedoraURL := dest + "abc-123"
	expectedFileContent := "sdfkjsdfkjsdfj"
	expectedContentType := "application/pdf"

	toTest := pass.DownloadService{
		Dest: dest,
		DOIs: MockLookupService(func(d string) (*pass.DoiInfo, error) {
			if d == doi {
				return &pass.DoiInfo{
					Manuscripts: []pass.Manuscript{{Location: location}},
				}, nil
			}

			return nil, errors.New("oops")
		}),
		HTTP: MockRequester(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodGet {
				t.Fatalf("Wrong http method, should be GET, got %s", req.Method)
			}
			if req.URL.String() == location {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(expectedFileContent)),
					Header: http.Header(map[string][]string{
						"Content-Type": {expectedContentType},
					}),
				}, nil
			}
			return nil, errors.New("oops")
		}),
		Fedora: MockBinaryStore(func(url string, body io.Reader, mimetype string) (string, error) {
			if url != dest {
				return "", fmt.Errorf("deposit expected into %s, instead was %s", dest, url)
			}

			if mimetype != expectedContentType {
				return "", fmt.Errorf("expected content type %s, instead got %s", expectedContentType, mimetype)
			}

			content, _ := ioutil.ReadAll(body)
			if string(content) != expectedFileContent {
				return "", errors.New("Did not get expected content")
			}

			return fedoraURL, nil
		}),
	}

	url, err := toTest.Download(doi, location)

	if url != fedoraURL {
		t.Errorf("Dowmload service should have returned fedora url %s, instead it returned %s", fedoraURL, url)
	}

	if err != nil {
		t.Errorf("Download service errored, %v", err)
	}
}
