package main_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	pass "github.com/oa-pass/pass-download-service"
)

func TestSimpleCase(t *testing.T) {
	cache := pass.NewDoiCache(pass.DoiCacheConfig{})

	expected := &pass.DoiInfo{
		Manuscripts: []pass.Manuscript{
			{
				Description: "Foo",
				Location:    "Bar",
			},
		},
	}

	manuscripts, _ := cache.GetOrAdd("foo", func() (*pass.DoiInfo, error) {
		return &pass.DoiInfo{
			Manuscripts: []pass.Manuscript{
				{
					Description: "Foo",
					Location:    "Bar",
				},
			},
		}, nil
	})

	diffs := deep.Equal(manuscripts, expected)
	if len(diffs) > 0 {
		t.Fatalf("Found differences between cached results and expected:\n%s", strings.Join(diffs, "\n"))
	}
}

func TestEvictSize(t *testing.T) {
	cache := pass.NewDoiCache(pass.DoiCacheConfig{
		MaxSize: 1,
		MaxAge:  1 * time.Second,
	})

	assertComputed(t, cache, "foo")
	assertNotComputed(t, cache, "foo") // foo should be cached here
	assertComputed(t, cache, "bar")
	assertComputed(t, cache, "foo") // no longer cached if it'd been evicted
}

func TestEvictTimeout(t *testing.T) {
	cache := pass.NewDoiCache(pass.DoiCacheConfig{
		MaxSize: 1,
		MaxAge:  1 * time.Microsecond,
	})

	assertComputed(t, cache, "foo")

	for i := 1; i < 10; i++ {
		if didCompute(cache, "foo") {
			return // good, it was evicted and had to be re-computed
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("Cache entry should have been evicted by now")
}

// Make sure only one simultaneous/contested add wins
func TestContested(t *testing.T) {
	cache := pass.NewDoiCache(pass.DoiCacheConfig{})

	errChan := make(chan error)
	exec1 := make(chan bool)
	ready1 := make(chan bool)
	ready2 := make(chan bool)
	result1 := make(chan *pass.DoiInfo)
	result2 := make(chan *pass.DoiInfo)

	expected := &pass.DoiInfo{
		Manuscripts: []pass.Manuscript{
			{
				Description: "Foo",
				Location:    "Bar",
			},
		},
	}

	// 1: This will execute and calculate the result once we signal it to do so
	// on the exec channel
	go func() {
		result, _ := cache.GetOrAdd("foo", func() (*pass.DoiInfo, error) {
			ready1 <- true
			<-exec1
			return &pass.DoiInfo{
				Manuscripts: []pass.Manuscript{
					{
						Description: "Foo",
						Location:    "Bar",
					},
				},
			}, nil
		})

		result1 <- result
	}()

	<-ready1 // Wait until our generator function is running, but paused until we signal

	// 2: This will block, and return the result from 1
	go func() {
		ready2 <- true
		result, _ := cache.GetOrAdd("foo", func() (*pass.DoiInfo, error) {
			// This shouldn't execute
			errChan <- errors.New("cache function executed when not expected to")
			return &pass.DoiInfo{}, nil
		})
		result2 <- result
		errChan <- nil
	}()

	<-ready2 // Wait until 2 blocks on 1 finishing

	exec1 <- true // Let 1 execute

	if diffs := deep.Equal(expected, <-result2); len(diffs) > 0 {
		t.Fatalf("Did not get expected cached result:\n%s", strings.Join(diffs, "\n"))
	}

	if diffs := deep.Equal(expected, <-result1); len(diffs) > 0 {
		t.Fatalf("Did not get expected cached result:\n%s", strings.Join(diffs, "\n"))
	}

	if err := <-errChan; err != nil {
		t.Fatal(err)
	}

}

func TestError(t *testing.T) {
	cache := pass.NewDoiCache(pass.DoiCacheConfig{})

	_, err := cache.GetOrAdd("foo", func() (*pass.DoiInfo, error) {
		return nil, fmt.Errorf("error")
	})

	if err == nil {
		t.Fatalf("Should have gotten an error!")
	}

	// After our error, we should be able to add just fine
	assertComputed(t, cache, "foo")
}

// expects the cache to execute the generator function for a given key
func assertComputed(t *testing.T, cache *pass.DoiCache, doi string) {
	t.Helper()
	if !didCompute(cache, doi) {
		t.Fatalf("Cache did compute and cache when expected")
	}
}

func assertNotComputed(t *testing.T, cache *pass.DoiCache, doi string) {
	t.Helper()
	if didCompute(cache, doi) {
		t.Fatalf("Cache computed and cached result when not expected to")
	}
}

func didCompute(cache *pass.DoiCache, doi string) bool {
	var computed bool
	_, _ = cache.GetOrAdd(doi, func() (*pass.DoiInfo, error) {
		computed = true
		return &pass.DoiInfo{
			Manuscripts: []pass.Manuscript{
				{
					Description: "Foo",
					Location:    "Bar",
				},
			},
		}, nil
	})

	return computed
}
