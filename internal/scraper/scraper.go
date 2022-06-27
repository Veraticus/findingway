package scraper

import (
	"fmt"
	"strings"

	"github.com/Veraticus/trappingway/internal/ffxiv"

	"github.com/gocolly/colly"
)

type Scraper struct {
	Url      string
	Listings *ffxiv.Listings
}

func New(url string) *Scraper {
	return &Scraper{
		Url:      url,
		Listings: &ffxiv.Listings{Listings: []*ffxiv.Listing{}},
	}
}

func (s *Scraper) Scrape() error {
	listings := []*ffxiv.Listing{}

	c := colly.NewCollector()

	fmt.Printf("Starting on URL: %v\n", s.Url)

	c.OnHTML("#listings.list .listing", func(e *colly.HTMLElement) {
		fmt.Printf("Working %v!\n", e)
		listing := &ffxiv.Listing{}

		// We can unmarshal a fair amount of information
		e.Unmarshal(listing)

		// Get attributes which are unmarshall-able
		listing.DataCentre = e.Attr("data-centre")
		listing.PfCategory = e.Attr("data-pf-category")

		// Get everything else that isn't easily inferred; first description
		description := e.ChildText(".left .description")
		description = strings.TrimSpace(strings.Replace(description, listing.Tags, "", -1))
		listing.Description = description

		// Then the party list
		e.ForEach(".party .slot", func(s int, p *colly.HTMLElement) {

		})

		fmt.Printf("Wound up with %+v\n", listing)
		listings = append(listings, listing)
	})

	c.Visit(s.Url + "/listings")

	s.Listings.Listings = listings

	return nil
}
