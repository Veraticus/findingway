package murult

type PfState struct {
	Listings map[string]*Listing
}

func NewPfState() *PfState {
	return &PfState{
		Listings: make(map[string]*Listing),
	}
}

func (pf *PfState) FilterForUltimatesInMateria(duties []string) {
	listings := make(map[string]*Listing)

	for _, l := range pf.Listings {
		if l.DataCentre == "Materia" {
			for _, d := range duties {
				if l.Duty == d {
					listings[l.Creator] = l
					break
				}
			}
		}
	}

	pf.Listings = listings
}

func (pf *PfState) Add(l *Listing) {
	pf.Listings[l.Creator] = l
}
