package scraper

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Veraticus/trappingway/internal/ffxiv"

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

	assert.Equal(t, 677, len(s.Listings.Listings))

	listing := s.Listings.Listings[100]
	assert.Equal(t, listing.Duty, "The Unending Coil of Bahamut (Ultimate)")
	assert.Equal(t, listing.Party[0].Job, ffxiv.AST)
	assert.True(t, listing.Party[0].Filled)
}
