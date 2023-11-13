package scraper

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Veraticus/findingway/internal/ffxiv"

	"github.com/stretchr/testify/assert"
)

func newTestServer() (*httptest.Server, error) {
	mux := http.NewServeMux()

	listings, err := os.ReadFile("testdata/listings.html")
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
	assert.NoError(t, err)
	defer ts.Close()

	s := &Scraper{Url: ts.URL}
	listings, err := s.Scrape()
	assert.NoError(t, err)

	assert.Equal(t, 656, len(listings.Listings))

	listing := listings.Listings[100]
	assert.Equal(t, listing.Duty, "The Unending Coil of Bahamut (Ultimate)")
	assert.Equal(t, listing.Party[0].Job, ffxiv.AST)
	assert.True(t, listing.Party[0].Filled)
}
