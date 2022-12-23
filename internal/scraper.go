package murult

import (
	"strings"

	"github.com/gocolly/colly"
)

type Scraper struct {
	Url string
}

func NewScraper(url string) *Scraper {
	return &Scraper{
		Url: url,
	}
}

func (s *Scraper) Scrape() (*PfState, error) {
	pf := NewPfState()
	collector := colly.NewCollector()

	collector.OnHTML("#listings.list .listing", func(e *colly.HTMLElement) {
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

		pf.Add(listing)
	})

	err := collector.Visit(s.Url)

	if err != nil {
		return nil, err
	} else {
		return pf, nil
	}
}