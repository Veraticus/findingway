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

func (ls *Listings) ForDataCentreAndDuty(datacentre, duty string) []*Listing {
	listings := []*Listing{}

	for _, l := range ls.Listings {
		if l.Duty == duty && l.DataCentre == datacentre {
			listings = append(listings, l)
		}
	}

	return listings
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
	return "<:ffxivhourglass:987141579879878676> " + l.Expires
}

func (l *Listing) GetUpdated() string {
	return "<:ffxivstopwatch:987141580869730324> " + l.Updated
}

func (l *Listing) GetTags() string {
	return strings.Replace(strings.Replace(l.Tags, "[", "", -1), "]", "", -1)
}
