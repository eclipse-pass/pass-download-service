package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-test/deep"
	pass "github.com/oa-pass/pass-download-service"
)

type MockLookupService func(string) (*pass.DoiInfo, error)

func (f MockLookupService) Lookup(doi string) (*pass.DoiInfo, error) {
	return f(doi)
}

type NoLookupService struct{}

func (n NoLookupService) Lookup(doi string) (*pass.DoiInfo, error) {
	return nil, nil
}

func TestMethodNotAllowed(t *testing.T) {

	for _, method := range []string{http.MethodPost, http.MethodDelete, http.MethodPut} {
		resp := httptest.NewRecorder()
		pass.LookupServiceHandler(NoLookupService{}).ServeHTTP(resp, httptest.NewRequest(method, "/foo", nil))

		if resp.Code != http.StatusMethodNotAllowed {
			t.Errorf("Method should not be allowed: %s", method)
		}
	}
}

func TestNoDoi(t *testing.T) {

	resp := httptest.NewRecorder()
	pass.LookupServiceHandler(NoLookupService{}).ServeHTTP(
		resp, httptest.NewRequest(http.MethodGet, "/foo?param=notDoi", nil))

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected bad request error code")
	}
}

func TestResponse(t *testing.T) {
	testDoi := "abc/123"
	info := &pass.DoiInfo{
		Manuscripts: []pass.Manuscript{
			{
				Description: "One",
				Location:    "http://example.org/first",
			},
			{
				Description: "Two",
				Location:    "http://example.org/second",
			},
		},
	}

	resp := httptest.NewRecorder()
	pass.LookupServiceHandler(MockLookupService(func(doi string) (*pass.DoiInfo, error) {
		if doi == testDoi {
			return info, nil
		}
		t.Fatalf("DOI didn't match!")
		return nil, nil
	})).ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/foo?doi="+testDoi, nil))

	var returnedDoiInfo pass.DoiInfo
	err := json.Unmarshal(resp.Body.Bytes(), &returnedDoiInfo)
	if err != nil {
		t.Fatalf("Encountered error reading response: %v", err)
	}

	diffs := deep.Equal(info, &returnedDoiInfo)
	if len(diffs) > 0 {
		t.Fatalf("Got different response than expected:\n%s", strings.Join(diffs, "\n"))
	}

	if !strings.Contains(resp.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("Bad content type: %s", resp.Header().Get("Content-Type"))
	}
}
