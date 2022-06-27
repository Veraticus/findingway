package scraper

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestServer() (*httptest.Server, error) {
	mux := http.NewServeMux()

	listings, err := ioutil.ReadFile("testdata/listings.html")
	if err != nil {
		return nil, err
	}

	mux.HandleFunc("/listings", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(listings)
	})

	return httptest.NewServer(mux), nil
}

func TestGetListings(t *testing.T) {
	ts, err := newTestServer()
	assert.Nil(t, err)
	defer ts.Close()

	s := New(ts.URL)
	s.Scrape()

	assert.True(t, false)
}
