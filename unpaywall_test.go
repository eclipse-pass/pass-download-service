package main_test

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/go-test/deep"
	pass "github.com/oa-pass/pass-download-service"
)

func TestUnpaywall(t *testing.T) {
	doi := "test/foo"
	email := "foo@example.org"
	baseuri := "http://example.org/unpaywall/v2"

	// Based on things in real_response.json
	expected := &pass.DoiInfo{
		Manuscripts: []pass.Manuscript{
			{
				Location:              "https://dash.harvard.edu/bitstream/1/12285462/1/Nanometer-Scale%20Thermometry.pdf",
				RepositoryInstitution: "Harvard University - Digital Access to Scholarship at Harvard (DASH)",
				Type:                  "application/pdf",
				Source:                "Unpaywall",
				Name:                  "Nanometer-Scale Thermometry.pdf",
			},
		},
	}

	file, err := os.Open("testdata/real_response.json")
	if err != nil {
		t.Fatalf("Could not open test response: %v", err)
	}

	toTest := pass.UnpaywallService{
		HTTP: MockRequester(func(req *http.Request) (*http.Response, error) {
			expectedURI := baseuri + "/" + doi + "?email=" + email
			if req.URL.String() != expectedURI {
				t.Fatalf("Did not get expected unpaywall request URL.  Expected: %s, got: %s", expectedURI, req.URL.String())
			}

			if req.Method != http.MethodGet {
				t.Fatalf("Expected GET method, got %s", req.Method)
			}

			return &http.Response{
				StatusCode: 200,
				Body:       file,
			}, nil
		}),
		Email:   email,
		Baseuri: baseuri,
		Cache:   pass.NewDoiCache(pass.DoiCacheConfig{}),
	}

	result, err := toTest.Lookup(doi)
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	diffs := deep.Equal(result, expected)
	if len(diffs) > 0 {
		t.Fatalf("result differed from expected %s", strings.Join(diffs, "\n"))
	}

}
