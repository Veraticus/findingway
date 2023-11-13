package scraper

import (
	"fmt"
	"strings"

	"github.com/Veraticus/findingway/internal/ffxiv"

	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	Url string
}

func (s *Scraper) Scrape() (*ffxiv.Listings, error) {
	listings := &ffxiv.Listings{}

	c := colly.NewCollector()
	errors := []error{}

	c.OnHTML("#listings.list .listing", func(e *colly.HTMLElement) {
		listing := &ffxiv.Listing{Party: []*ffxiv.Slot{}}

		// We can unmarshal a fair amount of information
		e.Unmarshal(listing)

		// Get attributes which are unmarshall-able
		listing.DataCentre = e.Attr("data-centre")
		listing.PfCategory = e.Attr("data-pf-category")
		listing.Id = e.Attr("data-id")

		// Get everything else that isn't easily inferred; first description
		description := e.ChildText(".left .description")
		description = strings.TrimSpace(strings.Replace(description, listing.Tags, "", -1))
		listing.Description = description

		// Then the party list
		e.ForEach(".party .slot", func(s int, p *colly.HTMLElement) {
			slot := ffxiv.NewSlot()
			class := p.Attr("class")

			if strings.Contains(class, "dps") {
				slot.Roles.Roles = append(slot.Roles.Roles, ffxiv.DPS)
			}

			if strings.Contains(class, "healer") {
				slot.Roles.Roles = append(slot.Roles.Roles, ffxiv.Healer)
			}

			if strings.Contains(class, "tank") {
				slot.Roles.Roles = append(slot.Roles.Roles, ffxiv.Tank)
			}

			if strings.Contains(class, "empty") {
				slot.Roles.Roles = append(slot.Roles.Roles, ffxiv.Empty)
			}

			if strings.Contains(class, "filled") {
				slot.Filled = true
				slot.Job = ffxiv.JobFromAbbreviation(p.Attr("title"))
			}

			listing.Party = append(listing.Party, slot)
		})

		listings.Add(listing)
	})
	c.Visit(s.Url + "/listings")

	if len(errors) > 0 {
		return nil, fmt.Errorf("Could not scrape listings: %w", errors[0])
	}

	return listings, nil
}
