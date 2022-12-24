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
		posts := &Post{Party: []*Slot{}}

		// We can unmarshal a fair amount of information
		e.Unmarshal(posts)

		// Get attributes which are unmarshall-able
		posts.DataCentre = e.Attr("data-centre")
		posts.PfCategory = e.Attr("data-pf-category")

		// Get everything else that isn't easily inferred; first description
		description := e.ChildText(".left .description")
		description = strings.TrimSpace(strings.Replace(description, posts.Tags, "", -1))
		posts.Description = description

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

			posts.Party = append(posts.Party, slot)
		})

		pf.Add(posts)
	})

	err := collector.Visit(s.Url)

	if err != nil {
		return nil, err
	} else {
		return pf, nil
	}
}
