package ffxiv

import (
	"strings"
)

type Listings struct {
	Listings []*Listing
}

type Listing struct {
	DataCentre  string
	PfCategory  string
	Duty        string `selector:".left .duty"`
	Tags        string `selector:".left .description span"`
	TagsColor   string `selector:".left .description span" attr:"class"`
	Description string
	MinIL       string `selector:".middle .stat .value"`
	Creator     string `selector:".right .creator .text"`
	World       string `selector:".right .world .text"`
	Expires     string `selector:".right .expires .text"`
	Updated     string `selector:".right .updated .text"`
	Party       []*Slot
}

type Slot struct {
	Roles  Roles
	Job    Job
	Filled bool
}

func NewSlot() *Slot {
	return &Slot{
		Roles: Roles{Roles: []Role{}},
	}
}

func (ls *Listings) ForDutyAndDataCentres(duty string, dataCentres []string) []*Listing {
	listings := []*Listing{}

	for _, l := range ls.Listings {
		if l.Duty == duty {
			for _, dataCentre := range dataCentres {
				if l.DataCentre == dataCentre {
					listings = append(listings, l)
				}
			}
		}
	}

	return listings
}

func (ls *Listings) Add(l *Listing) {
	for _, existingListing := range ls.Listings {
		if existingListing.Creator == l.Creator {
			return
		}
	}

	ls.Listings = append(ls.Listings, l)
}

func (l *Listing) PartyDisplay() string {
	var party strings.Builder

	for _, slot := range l.Party {
		if slot.Filled {
			party.WriteString(slot.Job.Emoji() + " ")
		} else {
			party.WriteString(slot.Roles.Emoji() + " ")
		}
	}

	return party.String()

}

func (l *Listing) GetExpires() string {
	return "<:hourglass:991379574187372655> " + l.Expires
}

func (l *Listing) GetUpdated() string {
	return "<:stopwatch:991379573000388758> " + l.Updated
}

func (l *Listing) GetTags() string {
	if len(l.Tags) == 0 {
		return "_ _"
	}
	return l.Tags
}

func (l *Listing) GetDescription() string {
	return "```" + l.Description + "```"
}
