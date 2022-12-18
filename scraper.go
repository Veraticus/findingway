package murult

import (
	"strings"

	"github.com/gocolly/colly"
)

type Scraper struct {
	Url      string
	Listings *Listings
}

func NewScraper(url string) *Scraper {
	return &Scraper{
		Url:      url,
		Listings: &Listings{Listings: []*Listing{}},
	}
}

func (s *Scraper) Scrape() error {
	listings := &Listings{}

	c := colly.NewCollector()

	c.OnHTML("#listings.list .listing", func(e *colly.HTMLElement) {
		listing := &Listing{Party: []*Slot{}}

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
			slot := NewSlot()
			class := p.Attr("class")

			if strings.Contains(class, "tank") {
				slot.Roles.Tank = true
			}

			if strings.Contains(class, "healer") {
				slot.Roles.Healer = true
			}

			if strings.Contains(class, "dps") {
				slot.Roles.Dps = true
			}

			if strings.Contains(class, "empty") {
				slot.Roles.Empty = true
			}

			if strings.Contains(class, "filled") {
				slot.Filled = true
				slot.Job = p.Attr("title")
			}

			listing.Party = append(listing.Party, slot)
		})

		listings.Add(listing)
	})

	c.Visit(s.Url)

	s.Listings = listings

	return nil
}
