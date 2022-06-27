package ffxiv

import ()

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

type Role int

const (
	DPS Role = iota
	Healer
	Tank
	Any
)

type Job int

const (
	GNB Job = iota
	PLD
	GLD
	DRK
	WAR
	MRD
	ACN // Arcanist
	SGE
	AST
	WHM
	CNJ
	SAM
	DRG
	NIN
	MNK
	RPR
	BRD
	MCH
	DNC
	BLM
	BLU
	SMN
	RDM
	LNC
	PUG
	ROG
	THM
	ARC // Archer
)

type Slot struct {
	Role   Role
	Job    Job
	Filled bool
}

func (ls *Listings) ForDuty(duty string) []*Listing {
	listings := []*Listing{}

	for _, l := range ls.Listings {
		if l.Duty == duty {
			listings = append(listings, l)
		}
	}

	return listings
}
