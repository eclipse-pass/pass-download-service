package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var httpClient *http.Client = &http.Client{
	Timeout: 10 * time.Second,
}

func TestIntegration(t *testing.T) {
	doi := "10.1038/nature12373"
	lookupEndpoint := "http://localhost:6502/lookup"
	downloadEndpoint := "http://localhost:6502/download"

	info := lookupDOI(t, fmt.Sprintf("%s?doi=%s", lookupEndpoint, url.QueryEscape(doi)))

	if len(info.Manuscripts) == 0 {
		t.Fatal("expected to see at least one manuscript")
	}

	downloadURL := info.Manuscripts[0].Location

	binaryURI := postBinary(t, doi, fmt.Sprintf("%s?doi=%s&url=%s", downloadEndpoint, url.QueryEscape(doi), downloadURL))

	println(binaryURI)

	// Now, make sure we can HEAD the created binary
	headRequest, _ := http.NewRequest(http.MethodHead, binaryURI, nil)
	headRequest.SetBasicAuth("fedoraAdmin", "moo")
	resp, err := mustSucceed(httpClient.Do(headRequest))
	if err != nil {
		t.Fatalf("HEAD of resulting binary failed: %v", err)
	}

	if resp.Header.Get("Content-Type") != "application/pdf;charset=ISO-8859-1" {
		t.Fatalf("Got wrong content type for PDF!")
	}

	if resp.Header.Get("Content-Length") != "1035596" {
		t.Fatalf("Didn't get expected content length for pdf file")
	}
}

func lookupDOI(t *testing.T, url string) *DoiInfo {
	resp, err := mustSucceed(httpClient.Get(url))
	if err != nil {
		t.Fatalf("could not look up DOI from unpaywall %v", err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var info *DoiInfo
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&info)
	if err != nil {
		t.Fatalf("Could not deserialize response %s\n%v", body, err)
	}

	return info

}

func postBinary(t *testing.T, doi string, url string) string {
	resp, err := mustSucceed(httpClient.Post(url, "text/plain", nil))
	if err != nil {
		t.Fatalf("could not download manuscript from %s:  %v", url, err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Could not read body: %v", err)
	}

	urlFromBody := string(body)
	urlFromHeader := resp.Header.Get("Location")

	if urlFromBody != urlFromHeader {
		t.Fatalf("URLs from body and Location header don't agree: %s, %s", urlFromBody, urlFromHeader)
	}

	return urlFromHeader
}
