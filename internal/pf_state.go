package murult

type PfState struct {
	listings map[string]*Listing
}

func NewPfState() *PfState {
	return &PfState{
		listings: make(map[string]*Listing),
	}
}

// GetListings returns an array of listings that is in the PF.
// Can be filtered based on the arguments.
func (pf *PfState) GetListings(duties []string) []*Listing {
	list := make([]*Listing, 0)

	for _, l := range pf.listings {
		if l.DataCentre == "Materia" {
			for _, d := range duties {
				if l.Duty == d {
					list = append(list, l)
					break
				}
			}
		}
	}

	return list
}

func (pf *PfState) Add(l *Listing) {
	// TODO: Check if we already have this creator
	// If we do, check which one is the latest one.
	pf.listings[l.Creator] = l
}
