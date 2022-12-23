package murult

type PfState struct {
	listings map[string]*Post
}

func NewPfState() *PfState {
	return &PfState{
		listings: make(map[string]*Post),
	}
}

// GetListings returns an array of listings that is in the PF.
// Can be filtered based on the arguments.
func (pf *PfState) GetListings(duties []string) map[string]*Post {
	list := make(map[string]*Post, 0)

	for _, l := range pf.listings {
		if l.DataCentre == "Materia" {
			for _, d := range duties {
				if l.Duty == d {
					list[l.Creator] = l
				}
			}
		}
	}

	return list
}

func (pf *PfState) Add(l *Post) {
	// TODO: Check if we already have this creator
	// If we do, check which one is the latest one.
	pf.listings[l.Creator] = l
}
